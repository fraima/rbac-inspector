package inspector

type k8s interface {
	ClusterRoleBuildingChan() (<-chan []ClusterRoleBuilding, error)
	Stop()
}

type inspector struct {
	k8s k8s
}
