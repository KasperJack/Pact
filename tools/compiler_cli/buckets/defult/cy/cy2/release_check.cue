
//import "list"

// Types must match those used in namespace.cue
#Option: {
	type:           "option"
	description:    string
	reserved_flags: [...string]
}

#Selection: {
	type:           "selection"
	description:    string
	reserved_flags: [...string]
}

// The namespace is imported from the user's namespace.cue
namespace: [string]: #Option | #Selection

// The release must satisfy this schema
release: #Release

#Release: {
	tags?:    #Tags
	options?: #Options
}

// ----- Tags (passive) -----
#Tags: {
    for k, v in namespace {
        if v.type == "option" {
            // string that equals one of the reserved_flags
            "\(k)": or(v.reserved_flags)
        }
        if v.type == "selection" {
            // list where each element is one of the reserved_flags
            "\(k)": [...or(v.reserved_flags)]
        }
    }
}

// ----- Options (active) -----
#Options: {
    for k, v in namespace {
        "\(k)"?: {
            [string]: {
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

