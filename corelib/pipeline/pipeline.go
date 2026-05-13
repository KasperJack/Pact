 package pipeline


type Checker interface {
    Validate() error
}

type Pipeline interface {
    Checker
    Build()  error
    //Hashes() error
    //Stage()  error
    //Push()   error
}
