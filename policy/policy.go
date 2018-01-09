package policy

type Condition struct {
}

type Policy struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Subjects    []string  `json:"subjects"`
	Actions     []string  `json:"actions"`
	Effect      string    `json:"effect"`
	Conditions  Condition `json:"conditions"`
	Resources   []string  `json:"resources"`

	Row int `json:"-"`
}
