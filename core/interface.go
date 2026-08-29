package core

type Option struct {
	Common
	Default     bool
	Label       string // optional, trim
	Description string // optional, trim
	Binding     []string
}

type Interface struct {
	User   []Option
	System []Option
}