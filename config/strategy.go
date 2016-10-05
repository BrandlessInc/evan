package config

type Strategy struct {
	Preconditions []Precondition
	Phases        []Phase
}

type Phase interface{}
