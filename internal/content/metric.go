package content

type Metric struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

func (m Metric) Type() string {
	return "metric"
}
