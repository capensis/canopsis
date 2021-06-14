// security contains implementation of authentication and authorization methods.
package security

import (
	"context"
	"path/filepath"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/security/mongoadapter"
	"github.com/casbin/casbin/v2"
)

// Enforcer is the API interface of casbin enforcer.
// Interface casbin.IEnforcer is not used because if cannot be mocked by mockgen.
type Enforcer interface {
	Enforce(rvals ...interface{}) (bool, error)
	StartAutoLoadPolicy(context.Context)
	LoadPolicy() error
}

type enforcer struct {
	*casbin.SyncedEnforcer
}

func (e *enforcer) StartAutoLoadPolicy(ctx context.Context) {
	e.SyncedEnforcer.StartAutoLoadPolicy(autoLoadInterval)
	defer e.SyncedEnforcer.StopAutoLoadPolicy()
	<-ctx.Done()
}

const modelFilePath = "/api/security/rbac_model.conf"
const autoLoadInterval = 10 * time.Second

// NewEnforcer creates new synced enforcer with mongo adapter.
func NewEnforcer(configDir string, client mongo.DbClient) (Enforcer, error) {
	a := mongoadapter.NewAdapter(client)
	casbinEnforcer, err := casbin.NewSyncedEnforcer(filepath.Join(configDir, modelFilePath), a)
	if err != nil {
		return nil, err
	}

	return &enforcer{casbinEnforcer}, nil
}
