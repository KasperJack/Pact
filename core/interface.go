package core



type OptionBody struct {
	Default     bool     `hcl:"default"`
	Label       string   `hcl:"label"`
	Description string   `hcl:"description,optional"`
	Binding     []string `hcl:"binding"`
}


type InterfaceBody struct {
	Options map[string]OptionBody
}

type Interface struct {
	User   *InterfaceBody
	System *InterfaceBody
}

func (i *Interface) HasUser() bool   { return i.User != nil }
func (i *Interface) HasSystem() bool { return i.System != nil }

func NewInterface() *Interface { return &Interface{} }

func NewInterfaceBody() *InterfaceBody {
	return &InterfaceBody{Options: make(map[string]OptionBody)}
}