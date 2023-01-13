package gotest

type Assertion struct {
	Method string   `json:"method"`
	Args   []string `json:"args"`
}

func MakeEqualAssertion() Assertion {
	return Assertion{
		Method: "equal",
		Args:   []string{"res.status", "200"},
	}
}

func MakePropertyAssertion() Assertion {
	return Assertion{
		Method: "property",
		Args:   []string{"res.body", `"id"`},
	}
}
