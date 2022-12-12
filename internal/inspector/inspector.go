package inspector

type k8s interface {
	ClusterRoleChan() (<-chan []ClusterRole, error)
	Stop()
}

type inspector struct {
}
