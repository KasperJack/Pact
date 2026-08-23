




















shortcut {  
    display_name = "gsd" // optional fall back to name from exe  
    exe  = "gg.exe"
    icon = "gg.exe" // optional fall back to exe icon 
    args = "--gg --ff" // optional
}







runtime_pact {

    filesystem {
        
        config {
            path = "%APPDATA%\\pkg.name\\pkg.name"
            policy_on_uninstall = "preserve_prompt"
        }

        logs {}

    }

}












scope {     
    "per_machine" = "%SystemDrive%\\ProramFiles" // main app data storage, this is where the app will store its data can be a path or  engine_managed, ("engine_managed" eveything will be linked against a junction manger will own the junctions and will manage the data) 
    // data = "engine_managed"
    "per_user" = "%APPDATA%"
}





scope {     
    data = "engine_managed" // or an exact path
}





lifecycle {
    update_strategy = "engine_managed" 
}








command  {
        exe  = "cli/wds.exe" // this will create a shim 
    }
command  {
        exe  = "cli/wds.exe" // this will create a shim 
    }
command  {
        exe  = "cli/wds.exe" // this will create a shim 
    }

command  "b"{
        exe  = "cli/wdhhjhjhhs.exe" // this will create a shim 
        args = "f"
    }
command  "t"{
        exe  = "cli/wdhs.exe" // this will create a shim 
    }
*/