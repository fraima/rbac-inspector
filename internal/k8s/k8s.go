package k8s

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	v1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	v1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	v1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/rest"

	"github.com/fraima/rbac-inspector/internal/inspector"
	"github.com/google/uuid"
)

type k8s struct {
	cliV1       v1.ClusterRoleBindingInterface
	cliV1alpha1 v1alpha1.ClusterRoleBindingInterface
	cliV1beta1  v1beta1.ClusterRoleBindingInterface

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

func (s *k8s) RbacChan() (<-chan []inspector.ClusterRoleBuilding, error) {
	watcherV1, err := s.cliV1.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("v1: %w", err)
	}
	watcherV1alpha1, err := s.cliV1alpha1.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("v1alpha1: %w", err)
	}
	watcherV1beta1, err := s.cliV1alpha1.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("v1beta1: %w", err)
	}

	rChan := make(chan []inspector.ClusterRoleBuilding)
	s.watchers.Store(uuid.New().String(), stopWatcher(rChan, watcherV1.Stop, watcherV1alpha1.Stop, watcherV1beta1.Stop))

	go func() {
		var event watch.Event
		for {
			select {
			case event = <-watcherV1.ResultChan():
			case event = <-watcherV1alpha1.ResultChan():
			case event = <-watcherV1beta1.ResultChan():
			}

			crbList, ok := event.Object.(*rbac.ClusterRoleBindingList)
			if !ok {
				zap.L().Warn("converting", zap.Any("event", event))
				continue
			}

			rChan <- convert(*crbList)
		}

	}()

	return rChan, nil
}

func (s *k8s) Stop() {
	s.watchers.Range(func(_, value any) bool {
		watcherStop := value.(func())
		watcherStop()
		return true
	})
}

func convert(src rbac.ClusterRoleBindingList) []inspector.ClusterRoleBuilding {
	resultList := make([]inspector.ClusterRoleBuilding, 0, len(src.Items))
	for _, i := range src.Items {
		ri := inspector.ClusterRoleBuilding{}
		for _, subject := range i.Subjects {
			switch subject.Kind {
			case rbac.GroupKind:
				ri.Group = subject.Name
			case rbac.ServiceAccountKind:
				ri.ServiceAccount = subject.Name
				ri.Namespace = subject.Namespace
			case rbac.UserKind:
				ri.User = subject.Name
			}
		}
		resultList = append(resultList, ri)
	}

	return resultList
}

func stopWatcher(rChan chan []inspector.ClusterRoleBuilding, stopWatchers ...func()) func() {
	return func() {
		for _, sw := range stopWatchers {
			sw()
		}
		close(rChan)
	}
}
