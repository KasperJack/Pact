    
namespace: {

    download: #Option & {
        description:    "how you want to download the binary"
        reserved_flags: ["http", "sftp","ssh"]
    }


    compiler: #Option & {
        description:    "what compiler backend you want to use"
        reserved_flags: ["gcc", "llvm", "msvc"]
    }



    languages: #Selection & {
        description:    "what languages you want"
        reserved_flags: ["en", "fr", "ar", "jp"]
    }



    extra: #Selection & {
        description:    "what extra content you want"
        reserved_flags: ["patch-01", "patch-02", "patch-03", "patch-04"]

    }


}


