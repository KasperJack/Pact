package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/hashicorp/hcl/v2/hclwrite"
)



func LockFile(lockData []byte) (core.LockFileC, error) {
	var config core.LockFileC

	err := hclsimple.Decode("package.hcl", lockData, nil, &config)
	if err != nil {
		return core.LockFileC{}, err
	}

	return config, nil
}

func WriteLockFile(lf core.LockFileC) []byte {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	for _, pkg := range lf.Packages {
		block := body.AppendNewBlock("package", []string{pkg.Identifier})
		b := block.Body()
		b.SetAttributeValue("version", cty.StringVal(pkg.Version))
		b.SetAttributeValue("upstream_version", cty.StringVal(pkg.UpstreamVersion))
		b.SetAttributeValue("installed_at", cty.StringVal(pkg.InstalledAt))
		b.SetAttributeValue("install_dir", cty.StringVal(pkg.InstallDir))
		b.SetAttributeValue("arch", cty.StringVal(pkg.Arch))
	}

	return f.Bytes()
}