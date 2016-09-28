package plugin

type (
	Rewrite struct {
		Base `json:",squash"`
		From string `json:"from"`
		To   string `json:"to"`
		Code string `json:"code"`
		When string `json:"when"`
	}
)
