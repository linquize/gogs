package main

//go:generate go run generate/generate.go -pkg=main -ignore="\\.DS_Store|README.md|TRANSLATORS|auth.d" -o=conf_gen.go -var=confAsset conf
//go:generate go run generate/generate.go -pkg=main -ignore="\\.DS_Store|less" -o=public_gen.go -var=publicAsset public
//go:generate go run generate/generate.go -pkg=main -ignore="\\.DS_Store" -o=templates_gen.go -var=templatesAsset templates

import (
	"os"

	"gogs.io/gogs/internal/assets/public"
	"gogs.io/gogs/internal/assets/templates"
	"gogs.io/gogs/internal/conf"
)

type ConfAsset struct {
}

func (c *ConfAsset) Asset(name string) ([]byte, error) {
	return confAsset.ReadFile(name)
}

func (c *ConfAsset) AssetDir(name string) ([]string, error) {
	dirs, err := confAsset.ReadDir(name)
	if err != nil {
		return []string{}, err
	}

	var r []string
	for _, dir := range dirs {
		r = append(r, dir.Name())
	}

	return r, nil
}

func (c *ConfAsset) MustAsset(name string) []byte {
	a, err := c.Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

type TemplatesAsset struct {
}

func (c *TemplatesAsset) Asset(name string) ([]byte, error) {
	return templatesAsset.ReadFile("templates/" + name)
}

func (c *TemplatesAsset) AssetNames() []string {
	return c.recurse("templates")
}

func (c *TemplatesAsset) recurse(dir string) []string {
	dirs, err := templatesAsset.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var r []string
	for _, d := range dirs {
		if d.IsDir() {
			for _, subdir := range c.recurse(dir + "/" + d.Name()) {
				r = append(r, subdir)
			}
		} else {
			r = append(r, (dir + "/" + d.Name())[len("templates/"):])
		}
	}

	return r
}

type PublicAsset struct {
}

func (c *PublicAsset) Asset(name string) ([]byte, error) {
	return publicAsset.ReadFile("public/" + name)
}

func (c *PublicAsset) AssetDir(name string) ([]string, error) {
	dirs, err := publicAsset.ReadDir("public/" + name)
	if err != nil {
		return []string{}, err
	}

	var r []string
	for _, dir := range dirs {
		r = append(r, dir.Name())
	}

	return r, nil
}

func (c *PublicAsset) AssetInfo(name string) (os.FileInfo, error) {
	f, err := publicAsset.Open("public/" + name)
	if err != nil {
		return nil, err
	}
	s, err2 := f.Stat()
	if err2 != nil {
		return nil, err2
	}
	return s, nil
}

func init() {
	conf.AssetPtr = &ConfAsset{}
	templates.AssetPtr = &TemplatesAsset{}
	public.AssetPtr = &PublicAsset{}
}
