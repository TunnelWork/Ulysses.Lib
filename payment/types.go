package payment

type (
	InputType string
)

// InputType compatible for gatewayConfigTemplate
var (
	Text     InputType = "text"
	Password InputType = "password"
	Textarea InputType = "textarea"

	Number InputType = "number"

	Radiogroup InputType = "radiogroup"
	Dropdown   InputType = "dropdown"
)

// P stands for Parameters and is a shortcut for map[string]interface{}
type P map[string]interface{}
