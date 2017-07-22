package stores

import (
	"testing"

	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/context"

	"github.com/stretchr/testify/require"
)

type stubApplication struct {
	name string
}

func (app *stubApplication) Name() string {
	return app.name
}

func (app *stubApplication) Repository() common.Repository {
	return nil
}

func (app *stubApplication) StrategyForEnvironment(environment string) common.Strategy {
	return nil
}

func newDeployment(application common.Application) *context.Deployment {
	var strategy common.Strategy = nil
	environment := "environment"
	ref := ""
	flags := map[string]interface{}{
		"force": false,
	}
	return context.NewDeployment(application, environment, strategy, ref, flags)
}

func TestEnqueuesDeployment(t *testing.T) {
	store := NewProcessLocalStore()
	application := &stubApplication { name: "application", }

	deployment1 := newDeployment(application)
	store.EnqueueDeployment(deployment1)

	deployment2 := newDeployment(application)
	store.EnqueueDeployment(deployment2)

	enqueuedDeployments := store.enqueuedDeployments["application"]["environment"]
	require.Equal(t, 2, len(enqueuedDeployments))
}

func TestFindsDeployment(t *testing.T) {
	store := NewProcessLocalStore()
	application := &stubApplication { name: "application", }

	deployment, err := store.FindDeployment(application, "environment")
	require.Nil(t, deployment)
	require.Nil(t, err)

	deployment = newDeployment(application)
	store.SaveDeployment(deployment)

	deployment, err = store.FindDeployment(application, "environment")
	require.Equal(t, deployment, deployment)
	require.Nil(t, err)
}
