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




type Manifest struct {
	Scope     *Scope
	Shortcuts *ActionSet[ShortcutBody]
	Commands  *ActionSet[CommandBody]
	
}

func NewManifest() *Manifest {
	return &Manifest{
		Scope:     nil,
		Shortcuts: NewActionSet[ShortcutBody](),
		Commands:  NewActionSet[CommandBody](),
	}
}







type Scope struct {
	User   *ScopeTarget `hcl:"user,block"`
	System *ScopeTarget `hcl:"system,block"`
}

type ScopeTarget struct {
	InstallPath string `hcl:"install_path,attr"`
}