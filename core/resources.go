package core

 
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)




var allScopes = []Scope{ScopeUser, ScopeSystem}



type Resource interface {
	Kind() string
	Scopes() []Scope
	BoundOption() string
	Key() string
}



type Bound struct {
	Option string `json:"-"` // excluded from Key
}



 
func (b Bound) BoundOption() string { return b.Option }




func keyOf(kind string, v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic("core: cannot encode " + kind + " resource: " + err.Error()) // programmer error
	}
	sum := sha256.Sum256(append([]byte(kind+"\x00"), b...))
	return kind + ":" + hex.EncodeToString(sum[:])[:12]
}









// ---------------------------------------------------------------- shortcut
 
type Shortcut struct {
	Bound
	Exe         string `json:"exe"`
	DisplayName string `json:"display_name,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Args        string `json:"args,omitempty"`
}
 
func (Shortcut) Kind() string    { return "shortcut" }
func (Shortcut) Scopes() []Scope { return allScopes }
func (s Shortcut) Key() string   { return keyOf(s.Kind(), s) }
 
// ----------------------------------------------------------------- command
 
type Command struct {
	Bound
	Exe  string `json:"exe"`
	Args string `json:"args,omitempty"`
}
 


func (Command) Kind() string    { return "command" }
func (Command) Scopes() []Scope { return allScopes }
func (c Command) Key() string   { return keyOf(c.Kind(), c) }
 

// ---------------------------------------------------------------- add_path
 
type AddPath struct {
	Bound
	Dir string `json:"dir"`
}
 
func (AddPath) Kind() string    { return "add_path" }
func (AddPath) Scopes() []Scope { return allScopes }
func (a AddPath) Key() string   { return keyOf(a.Kind(), a) }

