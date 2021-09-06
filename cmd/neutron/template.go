// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// cmd/neutron/template.yml
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _cmdNeutronTemplateYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x52\xcb\x8e\xdb\x30\x0c\xbc\xeb\x2b\xf8\x03\x75\xbc\x5d\xb4\x07\xdd\x82\xa2\x2f\x60\xb1\x2d\xd0\x0f\x30\x14\x89\x76\x84\xc8\x92\x2b\x51\x69\x02\xc3\xff\x5e\x50\x76\x63\xa7\x87\x45\x81\x5e\x37\xb9\x88\x33\xa3\xe1\x98\xd4\x70\xb2\x52\x00\x68\x25\x61\x87\xa4\x77\x1e\x0f\xd9\xa9\xdd\x38\x56\xcf\x48\x5e\xf5\x38\x4d\x3b\xad\x2a\x1d\x49\x00\x9c\xf0\xfa\x82\xcc\x63\xa6\x18\x7c\x75\xc2\x2b\x3b\x62\xa4\x7f\x10\xb3\xb1\x48\xa4\xc8\xea\xe6\x18\x12\x35\xbd\x1a\x38\xd0\x38\x46\xe5\x3b\x84\xea\x47\xe1\xbe\x84\x44\x69\x9a\x0a\x01\xd5\x73\x31\xdc\x1b\x13\x31\x25\x98\x26\x09\xe3\x58\x7d\xf4\x66\x08\xd6\xd3\xa2\x42\x6f\xa6\x49\x08\x67\xbb\x23\x1d\x43\x4e\xc8\xae\xaa\x6f\x36\x00\xdf\xda\xf7\x4f\x37\xa0\xdc\xb4\x9e\x30\x9e\x95\x93\xf0\xbe\x16\x00\x9c\x29\xdd\x05\x5a\xf5\x73\x20\x80\x37\x25\x14\xdc\x37\xd6\x76\x38\x62\x2c\x3d\x3e\x94\x23\x83\x43\xf6\xfa\x78\x65\xbb\x72\x92\x40\x31\x23\x87\x4c\x84\x5e\x2e\xed\x24\xd4\x55\xf9\xb3\x2c\xf0\x14\xd9\xfe\xa9\x68\xbe\x87\x48\xdc\x48\x50\x2e\x7a\x63\x93\x3a\x38\x34\x12\x5a\xe5\x12\x32\x82\xe7\xd2\xf4\x36\x6a\xb6\x31\x31\x0c\x8d\x0b\x5a\xb9\xe6\x10\x83\x32\x5a\x71\x9b\xdb\x15\x66\xfb\xec\xc8\xde\xe3\x74\x69\x7e\x66\xcc\x28\xe1\x5d\xcd\x2e\x3d\x65\x09\x0f\x8f\x75\x2d\x84\x0b\x5d\x67\x7d\xc7\x11\x1c\x9e\xd1\x49\xb0\xbe\x0d\x02\xa0\x0d\xb1\x57\x24\x81\xf0\x42\x42\xb4\x36\xe2\x2f\xe5\x5c\x79\x62\xc1\x7b\x8a\x4a\x9f\x64\x19\x1a\xe9\xa1\x21\xdb\x63\xc8\x24\xe1\xe1\x6d\x5f\xc0\x6c\x36\xe0\xe3\x8c\x19\x6c\x55\x76\xb4\x11\xd7\x33\xd1\xab\x4b\xc3\xa6\xa8\xc9\x06\x9f\x98\xe0\x9f\x00\x08\x99\x0e\x21\x7b\x33\x77\x1a\x47\x58\x56\xf7\x69\x89\xf3\x6d\xe1\x61\x59\x60\xd1\xd8\x16\xaa\xbd\xbf\xc2\x6d\xa9\xeb\xec\xff\x4c\x7d\xd6\x0e\x31\x50\x58\x08\x3e\xae\xcc\xbc\x3d\xe5\xaf\xab\x2b\xba\x84\xff\x63\xd9\xc5\x90\x87\x24\x97\x6a\xfb\x31\x9f\x0b\xb3\x2a\xef\xdf\xe1\x1a\xe0\xaf\xcf\xdc\xd4\xdb\xca\xfa\x17\x27\xf6\xd5\xbf\x0e\x6c\x5b\xfd\x0e\x00\x00\xff\xff\x25\x0b\x88\xc0\x37\x05\x00\x00")

func cmdNeutronTemplateYmlBytes() ([]byte, error) {
	return bindataRead(
		_cmdNeutronTemplateYml,
		"cmd/neutron/template.yml",
	)
}

func cmdNeutronTemplateYml() (*asset, error) {
	bytes, err := cmdNeutronTemplateYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/neutron/template.yml", size: 1335, mode: os.FileMode(420), modTime: time.Unix(1621717310, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"cmd/neutron/template.yml": cmdNeutronTemplateYml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"cmd": &bintree{nil, map[string]*bintree{
		"neutron": &bintree{nil, map[string]*bintree{
			"template.yml": &bintree{cmdNeutronTemplateYml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
