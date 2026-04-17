_all_flags: [ for ns in namespace for f in ns.reserved_flags { f } ]

_validate_uniqueness: {
    for i, f1 in _all_flags {
        for j, f2 in _all_flags if i < j {
            if f1 == f2 {
                //error 
                "DUPLICATE_FLAG_ERROR": "Flag '\(f1)' is used more than once!"
                "\(f1)": "conflict"
            }
        }
    }
}









release: {
    compiler: {
        "_": {
            flags: "gcc"
        }
       // this is a passive cause it prefixed with '_' the handller does not need to know about so "_" and it does not need a defult , passive act conly as discover for search it a way to tell this pacakge has only 'gcc' from option compiler if you have more you need to set a defult and tell the hanndler you var to expect
    }

    download: {
        "protocl": {
            flags: ["http","ssh"]
            // this has no defult it will promt the user during install 
        }
    }

    languages: {
        "_": {
            flags: ["en", "fr"]
            // this is passive but is is an array [] cause langaues is of type "selection" means this realse will install all of those 
        }


        // you can only use a name spaces once meaing you can't have a languages as passive and active only one of each 

        ui: {
            flags: ["en", "ar"]
            default: []
        }
    }

    extra: {
        "_": {
            flags: "patch-01"
        }

        optional: {
            flags: ["patch-01", "patch-03"]
        }
    }
}