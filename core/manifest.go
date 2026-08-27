package core


import("fmt")



type ActionSet[T any] struct {
	Unconditional []T
	Tagged        map[string]T
}

func NewActionSet[T any]() *ActionSet[T] {
	return &ActionSet[T]{
		Tagged: make(map[string]T),
	}
}

func (s *ActionSet[T]) AddUnconditional(body T) {
	s.Unconditional = append(s.Unconditional, body)
}

func (s *ActionSet[T]) AddTagged(id string, body T) error {
	if _, exists := s.Tagged[id]; exists {
		return fmt.Errorf("duplicate label %q", id)
	}
	s.Tagged[id] = body
	return nil
}



type AddPathBody struct {
	Dir string `hcl:"dir"`

}


type CommandBody struct {
	Exe  string `hcl:"exe"`
	Args string `hcl:"args,optional"`
}


type ShortcutBody struct {
	DisplayName string `hcl:"display_name,optional"`
	Exe         string `hcl:"exe"`
	Icon        string `hcl:"icon,optional"`
	Args        string `hcl:"args,optional"`

}








type ManifestBody struct {
	InstallPath string `hcl:"install_path"`
	Shortcuts   *ActionSet[ShortcutBody]
	Commands    *ActionSet[CommandBody]
	AddPaths    *ActionSet[AddPathBody]
}

type Manifest struct {
	User   *ManifestBody
	System *ManifestBody
}

func (m *Manifest) HasUser() bool {
	return m.User != nil
}

func (m *Manifest) HasSystem() bool {
	return m.System != nil
}

func NewManifest() *Manifest {
	return &Manifest{}
}

// NewManifestBody allocates an empty body with its ActionSets initialized.
// Call this when the parser encounters a user{} or system{} block.
func NewManifestBody() *ManifestBody {
	return &ManifestBody{
		Commands:  NewActionSet[CommandBody](),
		Shortcuts: NewActionSet[ShortcutBody](),
		AddPaths:  NewActionSet[AddPathBody](),
	}
}