package pipeline


type VerifierPackage interface {
    Validate() error
}

type PublisherPackage interface {
    VerifierPackage
    Build() error
    RebuildIndex() error
}