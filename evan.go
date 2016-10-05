package evan

// A single code-base deployed to 1+ targets for 1+ environments.
type Application struct {
	Targets              map[string]Target
	Environments         []string
	TargetForEnvironment func(string) *Target
}

type Target interface{}
