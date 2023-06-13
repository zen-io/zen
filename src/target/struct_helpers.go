package target

import environs "github.com/baulos-io/baulos/src/environments"

type BuildFields struct {
	Name       string   `mapstructure:"name" desc:"Name for the target"`
	Srcs       []string `mapstructure:"srcs" desc:"Sources for the build"`
	Outs       []string `mapstructure:"outs" desc:"Outs for the build"`
	BaseFields `mapstructure:",squash"`
}

type BaseFields struct {
	Name        string            `mapstructure:"name" desc:"Name for the target"`
	Description string            `mapstructure:"desc"`
	Labels      []string          `mapstructure:"labels" desc:"Labels to apply to the targets"` //
	Deps        []string          `mapstructure:"deps" desc:"Build dependencies"`
	PassEnv     []string          `mapstructure:"pass_env"`
	SecretEnv   []string          `mapstructure:"secret_env"`
	Env         map[string]string `mapstructure:"env"`
	Tools       map[string]string `mapstructure:"tools"`
	Visibility  []string          `mapstructure:"visibility"`
}

func (bf *BuildFields) GetBuildMods() []TargetOption {
	return bf.GetBaseMods()
}

func (bf *BaseFields) GetBaseMods() []TargetOption {
	return []TargetOption{
		WithVisibility(bf.Visibility),
		WithDescription(bf.Description),
		WithLabels(bf.Labels),
		WithPassEnv(bf.PassEnv),
		WithSecretEnvVars(bf.SecretEnv),
		WithEnvVars(bf.Env),
		WithTools(bf.Tools),
	}
}

type DeployFields struct {
	Environments map[string]*environs.Environment `mapstructure:"environments" desc:"Deployment Environments"`
}

func (df *DeployFields) GetDeployMods() []TargetOption {
	return []TargetOption{
		WithEnvironments(df.Environments),
	}
}
