package target

func (t *Target) EnvVars() map[string]string {
	return t.Env
}

func (t *Target) ShouldFlattenOuts() bool {
	return t.flattenOuts
}

func (t *Target) ShouldInterpolate() bool {
	return !t.noInterpolation
}

func (t *Target) Path() string {
	return t._original_path
}

func (t *Target) ShouldClean() bool {
	return t._clean
}
