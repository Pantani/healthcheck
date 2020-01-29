package fixtures

type Fixtures []Fixture

type Fixture struct {
	Name  string `json:"name"`
	Host  string `json:"host"`
	Tests []Test `json:"tests"`
}

type Test struct {
	Name       string      `json:"name"`
	UrlPath    string      `json:"url_path"`
	JsonPath   string      `json:"json_path"`
	Method     string      `json:"method"`
	UpdateTime string      `json:"update_time"`
	Expression string      `json:"expression"`
	Body       interface{} `json:"body"`
}
