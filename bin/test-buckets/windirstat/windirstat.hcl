identifier  = "windirstat"
name        = "WinDirStat"
versioning  = "semver"
description = "Disk usage viewer"
homepage    = "https://windirstat.net"
license     = "GPL-2.0"


//archfallback = true

//archpolicy = "interchangeable" 
//archpolicy = "strict"
// only one arch at a time for a version 
// pkgs are idendfied with pkgId/verion 

// pact list windirstat

// windirstat 2.7.0 [64bit] (active)
// windirstat 2.6.1 [64bit]
// windirstat 2.6.0 [32bit]

// pact use windirstat@2.6.1 #changes active version to 2.6.1

// pact update windirstat #updatest the active version to latest 
// error: alreadty found windirstat 2.7.0 installed 

// pact install windirstat@2.7.0 --arch 32  
// alreadty found windirstat 2.7.0 installed
// # remove windirstat@2.7.0 and install witrh same command 
// updates will look to satisfy latest version regaurdless of arch 
// pakes can from 32 --> 64 or vice versa to update to last version 








//index
type = "portable"
latest = "2.7.0"
versions = ["2.6.1","2.7.0"]






