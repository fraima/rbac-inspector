package k8s

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	v1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	v1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/rest"

	"github.com/fraima/rbac-inspector/internal/inspector"
)

const (
	versionV1       = "v1"
	versionV1alpha1 = "v1alpha1"
	versionV1beta1  = "v1beta1"
)

type k8s struct {
	cliV1       *v1.RbacV1Client
	cliV1alpha1 *v1alpha1.RbacV1alpha1Client
	cliV1beta1  *v1beta1.RbacV1beta1Client

	watchers sync.Map
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

	k := new(k8s)
	k.cliV1, err = v1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.cliV1alpha1, err = v1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.cliV1beta1, err = v1beta1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (s *k8s) ClusterRoleBuildingChan() (<-chan inspector.ClusterRoleBuilding, error) {
	watcherV1, err := s.cliV1.ClusterRoleBindings().Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", versionV1, err)
	}

	watcherV1alpha1, err := s.cliV1alpha1.ClusterRoleBindings().Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", versionV1alpha1, err)
	}
	watcherV1beta1, err := s.cliV1alpha1.ClusterRoleBindings().Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", versionV1beta1, err)
	}

	rChan := make(chan inspector.ClusterRoleBuilding)
	s.watchers.Store(uuid.New().String(), stopWatcher(rChan, watcherV1.Stop, watcherV1alpha1.Stop, watcherV1beta1.Stop))

	go func() {
		var (
			clb *inspector.ClusterRoleBuilding
			err error
		)
		for {
			select {
			case event := <-watcherV1.ResultChan():
				clb, err = convertClusterRoleBuildingV1(event)
			case event := <-watcherV1alpha1.ResultChan():
				clb, err = convertClusterRoleBuildingV1alpha1(event)
			case event := <-watcherV1beta1.ResultChan():
				clb, err = convertClusterRoleBuildingV1beta1(event)
			}
			if err != nil {
				zap.L().Warn("converting", zap.Any("event", err))
				continue
			}

			rChan <- *clb
		}

	}()

	return rChan, nil
}

func (s *k8s) ClusterRole(ctx context.Context, role inspector.ClusterRole) (inspector.ClusterRole, error) {
	switch role.Version {
	case versionV1:
		return s.getClusterRoleV1(ctx, role)
	case versionV1alpha1:
		return s.getClusterRoleV1alpha1(ctx, role)
	case versionV1beta1:
		return s.getClusterRoleV1beta1(ctx, role)
	}
	return role, fmt.Errorf("version %s is not exist", role.Version)
}

func (s *k8s) Stop() {
	s.watchers.Range(func(_, value any) bool {
		watcherStop := value.(func())
		watcherStop()
		return true
	})
}

func stopWatcher(rChan chan inspector.ClusterRoleBuilding, stopWatchers ...func()) func() {
	return func() {
		for _, sw := range stopWatchers {
			sw()
		}
		close(rChan)
	}
}
