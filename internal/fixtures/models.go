package fixtures

type Fixtures []Fixture

type Fixture struct {
	Namespace string `json:"namespace"`
	Host      string `json:"host"`
	Tests     []Test `json:"tests"`
}

type Test struct {
	Name       string      `json:"name"`
	URLPath    string      `json:"url_path"`
	JSONPath   string      `json:"json_path"`
	Method     string      `json:"method"`
	UpdateTime string      `json:"update_time"`
	Expression string      `json:"expression"`
	Body       interface{} `json:"body"`
}
