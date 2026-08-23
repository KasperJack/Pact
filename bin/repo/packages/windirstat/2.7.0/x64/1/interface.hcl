













option "desktop_shortcut" {
  type        = "boolean"
  default     = true
  label       = "Create Desktop Shortcut"
  description = "Adds a shortcut to your desktop for quick access."
  binding     = ["engine.shortcut.desktop"]
}



option "progamfiles_install" {
  type        = "boolean"
  label       = "Install Scope"
  description = "Install pkg in program files"

  binding     = "engine.scope.1"
}






option "install_scope" {
  type        = "enum"
  values      = ["per_user", "per_machine"]
  default     = "per_user"
  label       = "Install Scope"
  description = "Install for current user only or all system users."
  binding     = "engine.scope"
}








/*
option "install_scope" {
  type        = "Path" // will be checked if the path is valid and if the user has permission to write to it
  default     = "some/path"
  label       = "Install Scope"
  description = "Install for current user only or all system users."
  binding     = "engine.scope"
}
*/