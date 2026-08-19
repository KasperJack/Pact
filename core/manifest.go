package core



type Shortcut struct {
    Name string `hcl:"name,optional"`
    Exe  string `hcl:"exe"`
    Icon string `hcl:"icon,optional"`
    Args string `hcl:"args,optional"`
}



type Command struct {
    Exe  string `hcl:"exe"`
    Args string `hcl:"args,optional"`   // default args baked into shim
}
