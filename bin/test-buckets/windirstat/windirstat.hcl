identifier  = "windirstat"
name        = "WinDirStat"
versioning  = "semver"
description = "Disk usage viewer"
homepage    = "https://windirstat.net"
license     = "GPL-2.0"


## windows 
// name for shortcut if not avalable pcat will fall to pcakge name sortcut if not defined pact will use pkge identifier or the name from the exe


 //valid   
shortcut "WinDirStat" { 
    exe  = "WinDirStat.exe"  
}

//valid  
shortcut  { 
    exe  = "WinDirStat.exe"  
}

command  {
    exe  = "cli/WinDirStat.exe"
}




/*
detect {
    registry = "HKLM\\Software\\WinDirStat"
    path     = "C:\\Program Files\\WinDirStat\\WinDirStat.exe"
    env      = "WINDIRSTAT_HOME"
}*/