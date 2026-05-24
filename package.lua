
package {
    -- controled outer context 
    -- pckage called again insde can be controled 
    package_identifier = "python",
    name               = "Python",
    versioning         = "this" .. "that" ,

    description = "Python programming language",
    homepage    = "https://python.org",
    license     = "PSF",
    
    

    install = function()
        printFrlomGo("fd")
  
    end,

    function sayHello()
        print("Hello, World!")
    end
}


install = function()
    print("gg")
    printFromGo("hello fdsdsds")
    printFromGs("hello fdsdsds")
    print("gg")
end