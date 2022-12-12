package inspector

import (
	"context"
)

type k8s interface {
	GetClusterRoleBindingsList(ctx context.Context) (ClusterRoleBuilding, error)
}
