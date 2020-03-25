package main

// func (i Index) DownloadBuildForLocal(product, version string) error {
// 	p, ok := i.Products[product]
// 	if !ok {
// 		return fmt.Errorf("product %s not found", product)
// 	}
// 	v, ok := p.Versions[version]
// 	if !ok {
// 		return fmt.Errorf("version %s of %s not found", version, product)
// 	}
// 	build := v.GetBuildForLocal()
// 	if build == nil {
// 		return errors.New("no such build")
// 	}
// 	bts, err := build.Download()
// 	if err != nil {
// 		return err
// 	}
// 	if err = CheckBytes(build.Filename, bts); err != nil {
// 		return err
// 	}
// 	_, err = ExtractZip(product, "", bts)
// 	return err
// }
