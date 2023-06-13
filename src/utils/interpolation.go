package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func InterpolateMapWithItself(toInterpolate map[string]string) (map[string]string, error) {
	re := regexp.MustCompile(`\{[A-Z\.\_]+\}`)
	var err error
	needsInterpolation := true

	for needsInterpolation {
		needsInterpolation = false

		for key, val := range toInterpolate {
			if re.MatchString(val) {
				needsInterpolation = true
			}

			toInterpolate[key] = re.ReplaceAllStringFunc(val, func(m string) string {
				if val, ok := toInterpolate[m[1:len(m)-1]]; !ok {
					err = fmt.Errorf("%s is not a valid interpolation var", m)
				} else {
					return val
				}

				return ""
			})
		}

	}

	return toInterpolate, err
}

func Interpolate(text string, vars map[string]string) (string, error) {
	re := regexp.MustCompile(`\{[A-Z\.\_]+\}`)
	var err error
	var needsInterpolation bool

	retString := text
	for needsInterpolation = true; needsInterpolation; needsInterpolation = re.MatchString(retString) {
		retString = re.ReplaceAllStringFunc(retString, func(m string) string {
			if val, ok := vars[m[1:len(m)-1]]; !ok {
				err = fmt.Errorf("%s is not a valid interpolation var", m)
			} else {
				return val
			}

			return ""
		})
	}

	return retString, err
}

func InterpolateSetVars(text string, vars map[string]string) string {
	keys := []string{}
	for key := range vars {
		keys = append(keys, strings.ToUpper(key))
	}
	if len(keys) == 0 {
		return text
	}

	re := regexp.MustCompile(`\{(?:` + strings.Join(keys, "|") + `)\}`)
	var needsInterpolation bool

	retString := text
	for needsInterpolation = true; needsInterpolation; needsInterpolation = re.MatchString(retString) {
		retString = re.ReplaceAllStringFunc(retString, func(m string) string {
			return vars[strings.ToLower(m)[1:len(m)-1]]
		})
	}

	return retString
}

func InterpolateSlice(texts []string, vars map[string]string) ([]string, error) {
	var err error

	retString := []string{}
	for _, t := range texts {
		interpolatedText, err := Interpolate(t, vars)
		if err != nil {
			return nil, fmt.Errorf("interpolating %s: %w", t, err)
		}

		retString = append(retString, interpolatedText)
	}
	return retString, err
}

func InterpolateMap(m map[string]string, vars map[string]string) (map[string]string, error) {
	ret := map[string]string{}
	for k, v := range m {
		interpolatedKey, err := Interpolate(k, vars)
		if err != nil {
			return nil, fmt.Errorf("interpolating key %s: %w", k, err)
		}

		interpolatedVal, err := Interpolate(v, vars)
		if err != nil {
			return nil, fmt.Errorf("interpolating val %s: %w", v, err)
		}

		ret[interpolatedKey] = interpolatedVal
	}

	return ret, nil
}

func MergeMaps(maps ...map[string]string) map[string]string {
	merged := map[string]string{}

	for _, m := range maps {
		if m == nil {
			continue
		}

		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}
