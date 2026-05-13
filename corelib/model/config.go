package model


type Config struct {
  Release  *Release  `hcl:"release,block"`
  Package  *Package  `hcl:"package,block"`

}

