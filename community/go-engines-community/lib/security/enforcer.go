// security contains implementation of authentication and authorization methods.
package security

//go:generate mockgen -destination=../../mocks/lib/security/security.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security Enforcer,Provider,HttpProvider,UserProvider

import (
	"context"
	"path/filepath"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/mongoadapter"
	"github.com/casbin/casbin/v2"
)

// Enforcer is the API interface of casbin enforcer.
// Interface casbin.IEnforcer is not used because if cannot be mocked by mockgen.
type Enforcer interface {
	Enforce(rvals ...interface{}) (bool, error)
	StartAutoLoadPolicy(context.Context, time.Duration)
	LoadPolicy() error
	GetRolesForUser(name string, domain ...string) ([]string, error)
	GetPermissionsForUser(user string, domain ...string) ([][]string, error)
	HasPermissionForUser(user string, permission ...string) bool
}

type enforcer struct {
	*casbin.SyncedEnforcer
}

func (e *enforcer) StartAutoLoadPolicy(ctx context.Context, interval time.Duration) {
	e.SyncedEnforcer.StartAutoLoadPolicy(interval)
	defer e.SyncedEnforcer.StopAutoLoadPolicy()
	<-ctx.Done()
}

const modelFilePath = "/api/security/rbac_model.conf"

// NewEnforcer creates new synced enforcer with mongo adapter.
func NewEnforcer(configDir string, client mongo.DbClient) (Enforcer, error) {
	a := mongoadapter.NewAdapter(client)
	casbinEnforcer, err := casbin.NewSyncedEnforcer(filepath.Join(configDir, modelFilePath), a)
	if err != nil {
		return nil, err
	}

	return &enforcer{casbinEnforcer}, nil
}
