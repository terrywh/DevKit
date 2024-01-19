package k8s

// Request ...
type Request struct {
	ClusterID string `json:"cluster_id"`
	Namespace string `json:"namespace"`
	Pod string `json:"pod"`
	Command string `json:"command"`

	Rows int `json:"rows"`
	Cols int `json:"cols"`
}
type Description struct {
	ClusterID string `json:"cluster_id"`
	Namespace string `json:"namespace"`
	Svc []Svc `json:"svc"`
	Pod []Pod `json:"pod"`
	Node []Node `json:"node"`
}
type Svc struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ClusterIP string `json:"cluster_ip"`
	Port string `json:"port"`
}
type Pod struct {
	Name string `json:"name"`
	Node string `json:"node"`
	IP string `json:"ip"`
	Status string `json:"status"`
}
type Node struct {
	IP string `json:"ip"`
	Status string `json:"status"`
}

