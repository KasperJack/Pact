



def install():
    #path = 20
    """
    print(identifier)
    print(version)
    print(arch)
    print("------")
    print(staging_dir)
    print(install_dir)
    print("------")
    print(os.arch)
    print(os.is_x64())
    print(os.is_x86())
    print(os.is_arm64())
    """
    path.move_all(path.join(staging_dir,arch,"*"),install_dir)
    
