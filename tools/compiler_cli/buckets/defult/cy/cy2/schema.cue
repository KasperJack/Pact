package main

import "list"

// ---- Type definitions ----
// check ainterfases are not in the keywords "version,source"
#Option: {
    type:           "option"
    description:    string
    reserved_flags: list.MinItems(1) & [...string]
}

#Selection: {
    type:           "selection"
    description:    string
    reserved_flags: list.MinItems(1) & [...string]
}

// Every key in namespace must be an Option or Selection
namespace: [string]: #Option | #Selection

// ---- Duplicate flag check across all namespace entries ----
_allFlags: [ for _, v in namespace for f in v.reserved_flags { f } ]
_noDuplicates: list.UniqueItems & _allFlags




// ---- Release schema, auto-generated from namespace ----
#Release: {
    tags?: {
        for k, v in namespace if v.type == "option" {
            "\(k)"?: or(v.reserved_flags)
        }
        for k, v in namespace if v.type == "selection" {
            "\(k)"?: [...or(v.reserved_flags)]
        }
    }
    options?: {
        for k, v in namespace {
            "\(k)"?: [string]: {
                flags: [...or(v.reserved_flags)]
                if v.type == "option" {
                    default?: or(v.reserved_flags)
                }
                if v.type == "selection" {
                    default?: [...or(v.reserved_flags)]
                }
            }
        }
    }
}

// Every release in the package must satisfy #Release
release: #Release
