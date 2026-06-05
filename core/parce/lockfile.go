package parce

import (
	
	"github.com/kasperjack/pact/core/model"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/hashicorp/hcl/v2/hclwrite"
)



func LockFile(lockData []byte) (model.LockFile,error) {


	var config model.LockFile
	
    err := hclsimple.Decode("package.hcl", lockData, nil, &config)
    if err != nil {
        return  model.LockFile{},err
    }

	return config,nil
}



func WriteLockFile(lf model.LockFile) []byte {
    f := hclwrite.NewEmptyFile()
    body := f.Body()

    for _, pkg := range lf.Packages {
        block := body.AppendNewBlock("package", []string{pkg.Name})
        b := block.Body()
        b.SetAttributeValue("version",      cty.StringVal(pkg.Version))
        b.SetAttributeValue("installed_at", cty.StringVal(pkg.InstalledAt))
        b.SetAttributeValue("install_dir",  cty.StringVal(pkg.InstallDir))
    }

    return f.Bytes()
}