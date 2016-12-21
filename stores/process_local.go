package stores

import (
	"sync"

	"github.com/Everlane/evan/common"
)

// Stores deployments in the process's local memory.
type ProcessLocalStore struct {
	// Two-level map: first is application repository canonical name, second
	// is environment.
	deployments map[string]map[string]common.Deployment
	// Same map structure as `deployments`.
	enqueuedDeployments map[string]map[string][]common.Deployment

	mutex sync.Mutex
}

func NewProcessLocalStore() *ProcessLocalStore {
	return &ProcessLocalStore{
		deployments:         make(map[string]map[string]common.Deployment),
		enqueuedDeployments: make(map[string]map[string][]common.Deployment),
	}
}

func (store *ProcessLocalStore) SaveDeployment(deployment common.Deployment) error {
	application := store.keyForApplication(deployment.Application())
	environment := deployment.Environment()

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.deployments[application] == nil {
		store.deployments[application] = make(map[string]common.Deployment)
	}
	store.deployments[application][environment] = deployment
	return nil
}

func (store *ProcessLocalStore) FindDeployment(app common.Application, environment string) (common.Deployment, error) {
	application := store.keyForApplication(app)

	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.deployments[application][environment], nil
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

func (store *ProcessLocalStore) EnqueueDeployment(deployment common.Deployment) error {
	application := store.keyForApplication(deployment.Application())
	environment := deployment.Environment()

	store.mutex.Lock()
	defer store.mutex.Unlock()

	deployments := store.findOrCreateEnqueuedDeployments(application, environment)
	store.enqueuedDeployments[application][environment] = append(deployments, deployment)

	return nil
}

func (store *ProcessLocalStore) FindEnqueuedDeployments(app common.Application, environment string) ([]common.Deployment, error) {
	application := store.keyForApplication(app)

	return store.findOrCreateEnqueuedDeployments(application, environment), nil
}

// `make` doesn't support double-nesting so we need to create the inner nest.
// Returns the slice of enqueued deployments once everything is set up.
func (store *ProcessLocalStore) findOrCreateEnqueuedDeployments(application string, environment string) []common.Deployment {
	if store.enqueuedDeployments[application] == nil {
		store.enqueuedDeployments[application] = make(map[string][]common.Deployment)
	}

	return store.enqueuedDeployments[application][environment]
}

func (store *ProcessLocalStore) keyForApplication(application common.Application) string {
	return application.Name()
}
