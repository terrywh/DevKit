package entity

type Cluster struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
}

func (c *Cluster) ApplyDefaults() {
	if c.Namespace == "" {
		c.Namespace = "wemeet"
	}
}
