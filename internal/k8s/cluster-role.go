package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fraima/rbac-inspector/internal/inspector"
)

func (s *k8s) getClusterRoleV1(ctx context.Context, role inspector.ClusterRole) (inspector.ClusterRole, error) {
	cl, err := s.cliV1.ClusterRoles().Get(ctx, role.Name, metav1.GetOptions{})
	if err != nil {
		return role, fmt.Errorf("%s: %w", versionV1, err)
	}

	role.Rules = make([]inspector.Rule, 0, len(cl.Rules))
	for _, r := range cl.Rules {
		role.Rules = append(role.Rules,
			inspector.Rule{
				Verbs:     r.Verbs,
				Resources: r.Resources,
			},
		)
	}
	return role, nil
}

func (s *k8s) getClusterRoleV1alpha1(ctx context.Context, role inspector.ClusterRole) (inspector.ClusterRole, error) {
	cl, err := s.cliV1alpha1.ClusterRoles().Get(ctx, role.Name, metav1.GetOptions{})
	if err != nil {
		return role, fmt.Errorf("%s: %w", versionV1, err)
	}

	role.Rules = make([]inspector.Rule, 0, len(cl.Rules))
	for _, r := range cl.Rules {
		role.Rules = append(role.Rules,
			inspector.Rule{
				Verbs:     r.Verbs,
				Resources: r.Resources,
			},
		)
	}
	return role, nil
}

func (s *k8s) getClusterRoleV1beta1(ctx context.Context, role inspector.ClusterRole) (inspector.ClusterRole, error) {
	cl, err := s.cliV1beta1.ClusterRoles().Get(ctx, role.Name, metav1.GetOptions{})
	if err != nil {
		return role, fmt.Errorf("%s: %w", versionV1, err)
	}

	role.Rules = make([]inspector.Rule, 0, len(cl.Rules))
	for _, r := range cl.Rules {
		role.Rules = append(role.Rules,
			inspector.Rule{
				Verbs:     r.Verbs,
				Resources: r.Resources,
			},
		)
	}
	return role, nil
}
