latest_version = "2.7.0-1"

# upstream version -> latest full version for that upstream
version_mappings = {
  "2.7.0" = "2.7.0-1"
  "2.2.2" = "2.2.2-1"

}

# upstream version -> all known full versions/revisions for it
revision_history = {
  "2.7.0" = ["2.7.0-1"]
  "2.2.2" = ["2.2.2-1"]
}

# full version -> its upstream version (reverse lookup)
upstream_of = {

  "2.7.0-1" = "2.7.0"
  "2.2.2-1" = "2.2.2"
}


yanked = {}