scope {     
    "per_machine" = "%SystemDrive%\\ProramFiles" // main app data storage, this is where the app will store its data can be a path or  engine_managed, ("engine_managed" eveything will be linked against a junction manger will own the junctions and will manage the data) 
    // data = "engine_managed"
    "per_user" = "%APPDATA%"
}



scope {     
    data = "engine_managed" // or an exact path
}


scope {     
    // this is the case where scope was passed as a path
}


//integrations



runtime_pact {

    filesystem {
        
        config {
            path = "%APPDATA%\\pkg.name\\pkg.name"
            policy_on_uninstall = "preserve_prompt"
        }

        logs {}

    }

}





















shortcut "desktop"{  
    name = "hole" // optional fall back to name from exe  
    exe  = "WinDirStat.exe"
    icon = "WinDirStat.exe" // optional fall back to exe icon 
    args = "--ass --hole" // optional
}



command  {
        exe  = "cli/wds.exe" // this will create a shim 
    }

// add to path 
// add environment variable
// start menu
// servcies 
// secudeld tasks 
// drivers ??
// add/load registry keys (user, machine, classes)
// register file association and custom URLs like myapp://
// add Shell Context Menus (Adds options like "Open with MyEditor" to the Explorer right-click menu)
// add Auto-Start Rules




//state defention



lifecycle {
    update_strategy = "engine_managed" 
}
























/*



package
├── scope
├── options/interface
├── integrations
│   ├── shortcut
│   ├── command
│   ├── registry
│   └── ...
│
└── runtime
    ├── filesystem
    ├── config
    ├── logs
    └── ...


























id: publisher.appname
version: 1.4.2
display_name: "Awesome Editor"

# Execution & Binaries
architecture: x64
entry_points:
  - name: "Awesome Editor"
    executable: "bin/editor.exe"
    shortcut: true
    file_associations: [.txt, .md]

# Application Paths & Storage Contracts
storage:
  binaries: "%SystemDrive%\\Apps\\publisher.appname\\v1.4.2"
  config: "%APPDATA%\\publisher.appname\\config"
  runtime_data: "%LOCALAPPDATA%\\publisher.appname\\data"
  temp: "%LOCALAPPDATA%\\publisher.appname\\temp"

# State & Auto-Update Configuration
lifecycle:
  update_strategy: engine_managed  # or 'app_internal', 'disabled'
  update_check_url: "https://api.example.com/update.json"
  tracking_mode: system_journal    # Tracks any unmanifested changes made during runtime/install

# Registry Declarations (HKCU Preferred)
registry:
  hkcu:
    - path: "Software\\Publisher\\AwesomeEditor"
      values:
        - name: "InstallPath"
          type: REG_SZ
          value: "$INSTALL_DIR"

# Complex Hooks (Only executed inside isolated runner if required)
hooks:
  post_install: "scripts/setup_context_menu.ps1"

  */