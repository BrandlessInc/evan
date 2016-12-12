package common

type Application interface {
	Name() string
	Repository() Repository
	// Returns the strategy to use for deploying to a given environment or nil
	// if no strategy could be determined.
	StrategyForEnvironment(string) Strategy
}

type Repository interface {
	Owner() string
	Name() string
}

func CanonicalNameForRepository(repository Repository) string {
	owner := repository.Owner()
	name := repository.Name()
	return owner + "/" + name
}
