package strategy

// Describes how an application will be deployed to a environment & target.
type Strategy struct {
	Preconditions []Precondition
	Phases        []Phase
	Reporter      Reporter
}

type Phase interface{}
type Reporter interface{}
