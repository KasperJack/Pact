release: {


tags: {
        compiler: #OptionTag    & ["gcc"]         
        languages: #SelectionTag  & ["en", "fr"]
    }



public: {


    compiler: {
        backend: #Pick & {
            flags:   ["gcc", "llvm"]
            default: false
        }
    }


    languages: {
        region: #Set & {
            flags:   ["en", "fr", "ar"]
            default: ["en"]
        }
    }
}


}