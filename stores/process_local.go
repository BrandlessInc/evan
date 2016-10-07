package stores

import (
	"github.com/Everlane/evan/common"
)

// Stores deployments in the process's local memory.
type ProcessLocalStore struct {
	// Two-level map: first is application repository canonical name, second
	// is environment.
	applications map[string]map[string]common.Deployment
}

func (store *ProcessLocalStore) SaveDeployment(deployment common.Deployment) error {
	applicationKey := store.keyForApplication(deployment.Application())
	environment := deployment.Environment()
	store.applications[applicationKey][environment] = deployment
	return nil
}

func (store *ProcessLocalStore) keyForApplication(application common.Application) string {
	return common.CanonicalNameForRepository(application.Repository())
}

func (store *ProcessLocalStore) FindDeployment(application common.Application, environment string) (common.Deployment, error) {
	applicationKey := store.keyForApplication(application)
	return store.applications[applicationKey][environment], nil
}
