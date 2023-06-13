package target

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/baulos-io/baulos/src/utils"
)

// Checks if text is a reference to a target. This is true if the text starts with // or :
func IsTargetReference(text string) bool {
	return regexp.MustCompile(`^(\/\/|:)`).MatchString(text)
}

func CopyWithInterpolate(from, to string, target *Target, runCtx *RuntimeContext, interpolateVars ...map[string]string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return fmt.Errorf("reading from %s: %w", from, err)
	}

	if utils.CheckFileCanInterpolate(data) {
		interpolatedData, err := target.Interpolate(string(data), interpolateVars...)
		if err != nil {
			return fmt.Errorf("interpolating data: %w", err)
		}
		data = []byte(interpolatedData)
	}

	if err := os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return fmt.Errorf("opening folder %s: %w", filepath.Dir(to), err)
	}

	if err := os.WriteFile(to, data, 0644); err != nil {
		return fmt.Errorf("writing to %s, %w", to, err)
	}

	return nil
}

func (target *Target) Interpolate(text string, custom ...map[string]string) (string, error) {
	interpolateVars := utils.MergeMaps(
		append([]map[string]string{target.EnvVars()}, custom...)...,
	)

	return utils.Interpolate(text, interpolateVars)
}

func (target *Target) InterpolateMyself(runCtx *RuntimeContext) error {
	srcs := map[string][]string{}
	for sName, sSrcs := range target.Srcs {
		srcs[sName] = make([]string, 0)

		for _, src := range sSrcs {
			if interpolatedSrc, err := target.Interpolate(src); err != nil {
				return err
			} else {
				srcs[sName] = append(srcs[sName], interpolatedSrc)
			}
		}
	}

	target.Srcs = srcs

	outs := make([]string, 0)
	for _, o := range target.Outs {
		if interpolatedOut, err := target.Interpolate(o); err != nil {
			return err
		} else {
			outs = append(outs, interpolatedOut)
		}
	}
	target.Outs = outs

	for _, script := range target.Scripts {
		deps := []string{}

		for _, dep := range target.Scripts["build"].Deps {
			if interpolatedDep, err := target.Interpolate(dep); err != nil {
				return err
			} else {
				deps = append(deps, interpolatedDep)
			}
		}
		script.Deps = deps
	}

	tools := map[string]string{}
	for toolName, toolValue := range target.Tools {
		if interpolatedTool, err := target.Interpolate(toolValue); err != nil {
			return err
		} else {
			tools[toolName] = interpolatedTool
		}
	}
	target.Tools = tools

	return nil
}
