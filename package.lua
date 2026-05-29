
package {
    -- controled outer context 
    -- pckage called again insde can be controled 
    package_identifier = "python",
    name               = "Python",
    versioning         = "this" .. "that" ,

    description = "Python programming language",
    homepage    = "https://python.org",
    license     = "PSF",
    
    
    --persist = {
     --   "data/",
    --    "config/"
    --}

    function install ()
        
        ctx.fs.extract("/stage","/sofware-location")
        ctx.fs.glob()
    end






}


