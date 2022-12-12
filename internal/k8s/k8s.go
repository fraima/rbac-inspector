package k8s

import (
	"context"
	"os"

	"github.com/fraima/rbac-inspector/internal/inspector"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	v1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	v1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/rest"
)

type k8s struct {
	cliV1       v1.ClusterRoleBindingInterface
	cliV1alpha1 v1alpha1.ClusterRoleBindingInterface
	cliV1beta1  v1beta1.ClusterRoleBindingInterface
}

func Connect(kubeHost, kubeTokenFile string) (*k8s, error) {
	token, err := os.ReadFile(kubeTokenFile)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{
		Host:            kubeHost,
		APIPath:         "/",
		BearerToken:     string(token),
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}

	clientV1, err := v1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clientV1alpha1, err := v1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clientV1beta1, err := v1beta1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &k8s{
		cliV1:       clientV1.ClusterRoleBindings(),
		cliV1alpha1: clientV1alpha1.ClusterRoleBindings(),
		cliV1beta1:  clientV1beta1.ClusterRoleBindings(),
	}, nil
}

func (s *k8s) GetClusterRoleBuildingList(ctx context.Context) ([]inspector.ClusterRoleBuilding, error) {
	list, err := s.cliV1.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	resultList := make([]inspector.ClusterRoleBuilding, 0, len(list.Items))
	for _, i := range list.Items {
		resultList = append(resultList, convert(i))
	}
	return resultList, nil
}

func convert(src rbac.ClusterRoleBinding) (dst inspector.ClusterRoleBuilding) {
	for _, subject := range src.Subjects {
		switch subject.Kind {
		case rbac.GroupKind:
			dst.Group = subject.Name
		case rbac.ServiceAccountKind:
			dst.ServiceAccount = subject.Name
			dst.Namespace = subject.Namespace
		case rbac.UserKind:
			dst.User = subject.Name
		}
	}
	return
}
