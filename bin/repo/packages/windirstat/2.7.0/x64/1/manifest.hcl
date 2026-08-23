
scope {

  user {
    install_path = "%LOCALAPPDATA%\\windirstat"
  }



  system {
    install_path = "%PROGRAMFILES(X86)%\\windirstat"
  }


}


system {

    shortcut {  // in the system the short does not have an id it will run uncontianly 
        display_name = "hole" 
        exe  = "WinDirStat.exe"
    }


    

    add_path "system" {
        dir  = "cli"
        }




}



user {

    shortcut "desktop" {  
        display_name = "hole" 
        exe  = "WinDirStat.exe"
    }


    

    add_path "system" {
        dir  = "cli"
        }

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