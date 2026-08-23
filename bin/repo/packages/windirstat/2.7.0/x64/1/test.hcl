
scope {

  user {
    install_path = "%LOCALAPPDATA%\\windirstat"
  }


/*
  system {
    install_path = "%%PROGRAMFILES(X86)%/windirstat"
  }
 */

}







shortcut {  
    display_name = "gsd" // optional fall back to name from exe  
    exe  = "gg.exe"
    icon = "gg.exe" // optional fall back to exe icon 
    args = "--gg --ff" // optional
}



shortcut "desktop"{  
    display_name = "hole" // optional fall back to name from exe  
    exe  = "WinDirStat.exe"
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
