package k8s

import (
	"fmt"

	rbac "k8s.io/api/rbac/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/fraima/rbac-inspector/internal/inspector"
)

func convertClusterRoleBuildingV1(e watch.Event) (*inspector.ClusterRoleBuilding, error) {
	src, ok := e.Object.(*v1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("converting : %v", e)
	}

	clb := &inspector.ClusterRoleBuilding{
		Role: inspector.ClusterRole{
			Name: src.RoleRef.Name,
		},
	}
	for _, subject := range src.Subjects {
		switch subject.Kind {
		case rbac.GroupKind:
			clb.Group = subject.Name
		case rbac.ServiceAccountKind:
			clb.ServiceAccount = subject.Name
			clb.Namespace = subject.Namespace
		case rbac.UserKind:
			clb.User = subject.Name
		}
	}

	return clb, nil
}

func convertClusterRoleBuildingV1alpha1(e watch.Event) (*inspector.ClusterRoleBuilding, error) {
	src, ok := e.Object.(*v1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("converting : %v", e)
	}

	clb := &inspector.ClusterRoleBuilding{
		Role: inspector.ClusterRole{
			Name: src.RoleRef.Name,
		},
	}
	for _, subject := range src.Subjects {
		switch subject.Kind {
		case rbac.GroupKind:
			clb.Group = subject.Name
		case rbac.ServiceAccountKind:
			clb.ServiceAccount = subject.Name
			clb.Namespace = subject.Namespace
		case rbac.UserKind:
			clb.User = subject.Name
		}
	}

	return clb, nil
}

func convertClusterRoleBuildingV1beta1(e watch.Event) (*inspector.ClusterRoleBuilding, error) {
	src, ok := e.Object.(*v1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("%v", e)
	}

	clb := &inspector.ClusterRoleBuilding{
		Role: inspector.ClusterRole{
			Name: src.RoleRef.Name,
		},
	}
	for _, subject := range src.Subjects {
		switch subject.Kind {
		case rbac.GroupKind:
			clb.Group = subject.Name
		case rbac.ServiceAccountKind:
			clb.ServiceAccount = subject.Name
			clb.Namespace = subject.Namespace
		case rbac.UserKind:
			clb.User = subject.Name
		}
	}

	return clb, nil
}
