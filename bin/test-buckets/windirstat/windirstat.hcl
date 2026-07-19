identifier  = "windirstat"
name        = "WinDirStat"
versioning  = "semver"
description = "Disk usage viewer"
homepage    = "https://windirstat.net"
license     = "GPL-2.0"


## windows 
// name for shortcut if not avalable pcat will fall to pcakge name sortcut if not defined pact will use pkge identifier or the name from the exe

//TODO: assert exported entry execuatbles 

 //valid   

// install script for each version is granteed to produce a dir with the help of the script script should own delevery and installation 




match "<=2.0.3" {


    link {

    current ="install_dir" // should be done buy manger 


    }
  
    announce {

        shortcut {  
        name = "hole" // optional fall back to name from exe  
        exe  = "WinDirStat.exe"
        icon = "WinDirStat.exe" // optional fall back to exe icon 
        args = "--ass --hole" // optional
        }

        command  {
            exe  = "cli/wds.exe" // this will create a shim 
        }

    }        


}





//valid  
/*
shortcut  { 
    exe  = "WinDirStat.exe"  
}





/*
detect {
    registry = "HKLM\\Software\\WinDirStat"
    path     = "C:\\Program Files\\WinDirStat\\WinDirStat.exe"
    env      = "WINDIRSTAT_HOME"
}*/