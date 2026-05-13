 package verifier



type Pipeline interface {
    Validate() error
    Build()  error
    //Hashes() error
    //Stage()  error
    //Push()   error
}
