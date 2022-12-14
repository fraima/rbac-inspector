package inspector

type ClusterRoleBuilding struct {
	Role ClusterRole

	Group          string
	ServiceAccount string
	Namespace      string
	User           string
}

type ClusterRole struct {
	Name    string
	Version string

	Rules []Rule
}

type Rule struct {
	Verbs     []string
	Resources []string
}
