package stores

import (
	"sync"

	"github.com/Everlane/evan/common"
)

// Stores deployments in the process's local memory.
type ProcessLocalStore struct {
	// Two-level map: first is application repository canonical name, second
	// is environment.
	applications map[string]map[string]common.Deployment
	mutex sync.Mutex
}

func NewProcessLocalStore() *ProcessLocalStore {
	return &ProcessLocalStore{
		applications: make(map[string]map[string]common.Deployment),
	}
}

func (store *ProcessLocalStore) SaveDeployment(deployment common.Deployment) error {
	application := store.keyForApplication(deployment.Application())
	environment := deployment.Environment()

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.applications[application] == nil {
		store.applications[application] = make(map[string]common.Deployment)
	}
	store.applications[application][environment] = deployment
	return nil
}

func (store *ProcessLocalStore) FindDeployment(app common.Application, environment string) (common.Deployment, error) {
	application := store.keyForApplication(app)

	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.applications[application][environment], nil
}

func (store *ProcessLocalStore) HasActiveDeployment(app common.Application, environment string) (bool, error) {
	deployment, err := store.FindDeployment(app, environment)
	if err != nil {
		return false, err
	}
	if deployment == nil {
		return false, nil
	}

	switch deployment.Status().State {
	case common.DEPLOYMENT_PENDING:
	case common.RUNNING_PRECONDITIONS:
	case common.RUNNING_PHASE:
		return true, nil
	}
	return false, nil
}

func (store *ProcessLocalStore) keyForApplication(application common.Application) string {
	return common.CanonicalNameForRepository(application.Repository())
}
