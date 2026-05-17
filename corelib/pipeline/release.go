package pipeline


type VerifierRelease interface {
    Validate() error
}

type PublisherRelease interface {
    VerifierRelease
    Build() error
}