// Code generated by go-bindata. DO NOT EDIT.
// sources:
// views/aside.html
// views/base.html
// views/chartjs.html
// views/footer.html
// views/header.html
// views/revel.html
package views

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _asideHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xc1\x6e\xf3\x20\x10\x84\xef\xff\x53\xec\x8f\x72\xc5\xdc\x5d\xcc\xa5\xa7\x3e\xc6\xc6\x5e\xc7\x48\x64\xdd\x02\xae\x54\x21\xde\xbd\x0a\x81\x28\x89\xab\x4a\x3d\x7a\x99\xfd\x66\x3d\x1a\x8d\xc1\x4e\x04\xa3\xc3\x10\x06\x71\x46\xcb\xf2\x32\x38\xa2\x17\xe6\x1f\x00\x80\xfe\x2f\x25\xd4\x51\x0f\x21\x7e\x39\x82\x11\x19\x8e\x04\xf3\xba\xf1\x04\x96\xdb\x73\xe7\x28\x04\x90\xb2\x2e\x06\x1a\xa3\x5d\xb9\xb1\x1b\xf6\xca\x18\xc4\x42\xf6\xb4\xc4\x1e\x70\x8b\xeb\x4b\x35\x7b\x36\x84\x33\xf1\xd6\xc3\x9f\x7c\x0b\x63\x73\x4f\xb6\xf2\x42\xba\x73\x29\x2a\x67\x9b\x6a\x21\x9c\xc8\x3f\xbd\x17\xcd\x48\x1c\xc9\x1b\x1d\xa2\x5f\xf9\x64\xfc\x74\x9c\xad\xa3\xa0\x55\x1d\x68\x55\x15\x8f\x68\xe5\xec\xce\x6c\x4f\x4f\xc9\x23\x9f\x08\x0e\x96\x43\x44\x1e\x09\xfa\x01\x0e\xdd\x5b\xfd\x0a\x39\xef\x0f\x72\x16\x52\xb2\x33\xd0\xc7\xdd\xda\xa1\x7b\xdd\xbc\x27\x8e\x6d\x35\xe7\xfa\x67\x38\x46\xfb\x49\x02\x52\x22\x9e\x72\xde\xdf\x50\xa0\x08\x8b\xa7\x79\x10\xaa\x11\x55\x4a\x37\x7a\xce\xc2\xe8\x5b\x54\x33\xc2\x8c\x32\x2e\xc2\x68\x65\x0d\xe8\xf0\x8e\x6c\x1e\xd4\x5a\x95\x99\x56\xf8\x43\x9e\xbb\x60\xae\x39\x94\xe3\x7e\x89\x50\xab\xcd\xd5\x5e\xa9\x5a\xac\xbb\x7e\xaa\xae\x15\xe6\xd2\x02\xad\x4a\xa9\xcd\x77\x00\x00\x00\xff\xff\x31\x8f\xab\xdb\xdc\x02\x00\x00")

func asideHtmlBytes() ([]byte, error) {
	return bindataRead(
		_asideHtml,
		"aside.html",
	)
}

func asideHtml() (*asset, error) {
	bytes, err := asideHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "aside.html", size: 732, mode: os.FileMode(420), modTime: time.Unix(1636423570, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _baseHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\x51\x8f\xdb\x36\x0c\x7e\xcf\xaf\x60\xf5\x5c\x5b\xc9\x5d\x0b\x74\xab\x9d\xa1\x4b\x3b\x60\x40\x81\x6d\xdd\xbd\xec\x69\x50\x2c\xe6\xcc\x9c\x2c\x19\x22\x9d\x5c\x56\xf4\xbf\x0f\xb2\xe3\x25\xe9\xae\xd7\x5c\x70\x7e\x91\x44\xf2\xe3\x27\x89\xd4\x07\x17\x2f\xde\xff\xb6\xb8\xf9\xeb\xf7\x0f\x50\x4b\xe3\xe6\x93\x22\x0d\xe0\x8c\xbf\x2d\xd5\x3f\x75\x56\x79\x35\x9f\x4c\x8a\x1a\x8d\x9d\x4f\x00\x00\x0a\x21\x71\x38\xff\xf4\xfe\x53\xa1\x87\xe9\x60\x6e\x50\x0c\x54\xc1\x0b\x7a\x29\xd5\x96\xac\xd4\xa5\xc5\x0d\x55\x98\xf5\x8b\x97\x40\x9e\x84\x8c\xcb\xb8\x32\x0e\xcb\xd9\x4b\x68\xcc\x3d\x35\x5d\x73\x30\x74\x8c\xb1\x5f\x99\xa5\xc3\xd2\x07\x05\xde\x34\x58\xaa\x0d\xe1\xb6\x0d\x51\xd4\x31\xd7\xe0\xb2\xc8\x55\xa4\x56\x28\x78\x75\xa0\x7f\x20\xd0\x74\x52\x87\xf8\x40\xcc\x8b\x2c\x83\x9f\x43\x10\x96\x68\x5a\xb8\xce\xaf\xf3\x57\x90\x65\x7b\xa7\x23\x7f\x07\x75\xc4\x55\xa9\x34\x8b\x11\xaa\xf4\x72\x8c\xd5\x15\xf3\x61\x95\x37\xe4\xf3\x8a\x59\x41\x44\x57\x2a\x96\x9d\x43\xae\x11\x45\x81\xec\x5a\x2c\x95\xe0\xbd\xe8\x3e\x40\x1f\x31\xff\x12\xbc\xbc\xdb\x22\x87\x06\xe1\x55\x7e\x9d\x4f\x1f\xe5\x6e\x5d\x77\x4b\x9e\x35\x55\xc1\xb3\x5e\x05\x2f\x99\x19\xc0\x17\xd2\xff\x1a\x7c\x9f\x0b\xae\xf2\xe9\x93\xb8\x69\x0f\xbc\x90\xf7\xa6\xc6\x06\xa1\x8f\x7e\x94\xd4\x12\xf7\x68\xfd\xce\x36\xe4\x3f\xde\x7c\xb8\x80\x6b\x84\xc2\x9f\x77\xe4\x39\x87\x45\x1d\x02\x23\x18\xe0\x3b\xf2\xb0\x8a\xa1\x01\xa9\x11\x12\x4b\xb2\x70\x8f\xec\xbf\x55\x70\x16\x23\x90\x67\x41\x63\x21\xac\xc0\x86\xad\x77\xc1\x58\xf2\xb7\x60\x9c\x4b\x26\xa9\xb1\x01\x09\x10\xd1\x76\x15\xf6\xa9\x52\x44\x7e\xde\xb9\x7a\x46\xfd\xb7\x71\x2e\xeb\xa7\x17\xde\xe7\xfa\x8f\x0e\xe3\x0e\xae\xf2\xd9\x71\xf7\x0e\x4f\x03\x38\x56\xff\x2f\xe3\x80\xd8\x0f\x59\x0f\xec\xb9\xd7\xac\xe6\x85\x1e\x90\x8f\xa4\xa9\x6a\x13\x65\xcd\x7a\x91\xc6\xa7\x82\x3a\x21\xc7\xe7\x82\xd6\xac\x39\x44\x91\xa4\x09\x5f\x61\x0e\xe7\xa7\x45\x8d\xd5\xdd\x79\xfd\xdb\x87\xea\x95\x33\xa2\x97\xae\xc3\xf3\x6f\x7b\x52\xe8\x41\x07\x27\xc5\x32\xd8\x1d\x54\xce\x30\x97\x2a\xd5\x2d\x4b\x99\x80\xc9\xe2\xd2\xc4\xac\x21\x4f\xb0\xdd\x31\x25\x2d\x7d\x9d\x71\xd7\x26\xf5\x42\x3b\x4a\x8e\xa5\xcd\x08\xde\x46\xd3\xb6\x18\xd5\xf0\x16\x4a\xd5\x90\xcf\x6a\xa4\xdb\x5a\x7e\x84\x37\xd3\x69\x7b\xff\x16\xc6\xa5\xe9\x24\xbc\xdd\xa7\x48\xdf\xe7\xcf\x82\x4d\xeb\x8c\x20\xa8\xb4\x2f\x8c\x79\xe2\x53\x90\x7f\xf9\x72\xe2\x34\x69\x5f\xdf\xf0\xed\x4b\x72\xe2\x85\xfc\xa3\xd9\x85\x4e\x16\x83\x54\xc2\x57\x90\x55\x08\x72\xcc\x35\x9c\x49\x5b\xda\x9c\x56\xf1\xe8\x06\xd7\x66\x63\x06\xab\x3a\x2d\xee\x41\x4c\xd7\x47\x5a\x9a\x49\x08\x4e\xa8\xcd\xd7\xfc\xd3\xa6\xbc\x9a\xce\x5e\x4f\xdf\xcc\x7e\x98\xce\xbe\xd5\x2f\x97\x33\x7d\xa7\xe5\xcf\x49\x3c\xf6\xd5\xca\xb0\x54\x8e\x52\x6b\x8d\xb3\x67\xc8\xde\x0b\xc5\x9a\xb5\x69\x9f\x63\xaf\x63\x36\x8b\x4d\x78\xa6\x63\xa7\xae\x5c\x9a\xda\xff\x37\x79\x60\x9b\x85\x4e\xef\x65\x78\x40\xe9\x27\xe3\xdf\x00\x00\x00\xff\xff\x1b\x0f\x4d\xf9\x74\x08\x00\x00")

func baseHtmlBytes() ([]byte, error) {
	return bindataRead(
		_baseHtml,
		"base.html",
	)
}

func baseHtml() (*asset, error) {
	bytes, err := baseHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "base.html", size: 2164, mode: os.FileMode(420), modTime: time.Unix(1636423570, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _chartjsHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\xdf\x8f\x9b\xb8\x16\x7e\xef\x5f\x71\xc4\xad\x44\xa2\x66\x68\xe8\xed\x13\x6d\x2a\xcd\x4c\x6f\x75\xbb\x5b\xb5\x95\x5a\xad\xb6\x8a\xf2\xe0\x84\x13\x62\xd5\xd8\xc8\x36\x93\xd0\x88\xff\x7d\x65\x0c\x04\xc2\x8f\x61\x57\xd5\xbe\xec\xfa\x61\x06\xf0\x77\x3e\x7f\x3e\x3e\xfe\x30\x79\xad\x76\x92\x26\xfa\xcd\x13\x00\x80\xfb\x03\x91\xda\x4b\x58\x1a\x51\xfe\x05\xe5\x03\xdd\xa1\x27\x31\xa2\x4a\xa3\x9c\x9d\x0b\x88\x69\x64\xaf\x51\xbe\x95\xe4\x18\xc0\x3e\xe5\x3b\x4d\x05\x9f\xed\x4c\xe8\x1c\x2e\x20\xd3\xe8\x1e\x6c\x87\xb7\x13\x7c\x4f\x23\x4f\x24\x06\xac\x3c\x64\x18\x23\xd7\xca\xdb\x21\xd7\x28\xaf\xe3\x4c\x7b\x20\x12\x0e\xc8\x12\x94\x0a\x56\xa5\xb2\xf2\xfe\x55\x2f\xd8\x52\xfd\x0e\xab\x7a\x4c\xf3\xf7\x56\x22\xf1\x18\xee\x35\x3c\x83\xeb\xc7\x92\x46\x07\x3d\x87\xe7\xf0\x62\x8c\xf1\x5b\x1f\xa3\x16\x49\x0f\xe1\x56\x68\x2d\xe2\x92\xb1\x9f\x52\x9f\x60\xd5\x8c\xf3\x76\xfa\xd4\x1d\x7d\xa7\x4f\x9e\x22\x0f\x38\x9b\xf7\x2b\xdb\x0b\xae\xbf\xd0\x1f\x08\xab\x2a\x47\x5e\x84\xfa\x37\xc2\x52\xfc\x24\xdf\xe2\x9e\xa4\x4c\x4f\xca\xbc\x57\x31\x2d\xca\x14\x87\x36\x58\x79\x11\x13\x5b\xc2\xaa\xfb\x77\x25\x6c\x4c\x8f\xce\xd8\xcf\x12\x64\xa8\xa6\x28\x32\xb8\x11\x49\xef\x48\x4c\x59\xf6\x73\x34\x59\xae\x09\xa2\x2c\x70\x44\x55\x43\x8f\x9d\xac\xa4\x3c\x9a\x5d\x16\xa2\x91\x81\xfd\x18\x9f\x29\x92\x92\xcf\xfc\x1b\x00\x50\xc6\x7e\xe2\xca\xdc\x0b\x26\xe4\x84\x24\x14\xb8\x01\xcd\x1a\x4f\xfa\x96\xd1\x88\xc3\x0a\x5c\xcb\xed\x0e\x23\xef\x88\x42\x46\xb9\xd1\xef\xc6\x34\x0c\x19\x0e\x80\xcd\x4c\xbf\xe2\x69\xe2\x64\x0c\xf5\xa2\xf2\x8c\xea\xe2\xdb\x80\x62\x89\x4a\x0b\xd9\xd9\x8d\x79\x7d\x97\x2f\x8a\xcb\x7c\x6e\x77\xbd\x59\xe9\x84\xe2\x7d\x21\x01\x56\x0d\x77\xd3\x59\x82\x01\xb8\xa1\x48\xa3\x03\x4f\xb5\xbb\xa8\x7b\x42\xa2\x49\x70\xe5\x83\xe6\x99\x42\xad\x02\x58\x77\x0d\xd2\x06\xac\x37\x8b\x4e\xcf\x96\xec\xbe\x47\x52\xa4\x3c\x2c\x96\x21\x80\x75\x07\x62\xda\x91\xf2\x50\x1c\xad\x0d\x15\x40\xe5\x49\x0c\xbb\x74\x03\x58\x21\x09\x8f\x70\x32\x3c\x43\xc6\xc4\x71\x32\x3c\x92\x88\x7c\x32\x7a\xcb\xd2\x1e\x25\x57\xa9\xc9\xaf\xee\x19\xd9\x22\x33\xb9\xdd\x5c\xaf\xa3\x69\x65\xe1\x5c\x2f\x89\xfe\x8c\xd2\x14\x0b\x89\x30\x80\x97\xcb\x36\xa3\x44\x95\x08\xae\xe8\x03\x06\xa0\xe5\xb5\x24\x86\x11\xf2\xf0\x9a\xd0\xb4\x90\xaa\x84\x91\x2c\x80\x3d\x61\x0a\xdb\xa2\xdb\x1c\x9a\x6a\x86\xa3\x14\xdd\x71\x4d\xab\xf7\x6e\x00\xce\x7f\x96\xcb\xa5\xd3\x0f\x31\x06\x14\x80\xff\xdf\x6e\xaf\xd9\x2f\x01\xb8\xdf\x31\x03\x9e\xc6\x5b\x94\xb0\xcd\x8a\x6a\x56\xee\x98\xde\x6a\xdb\xf5\x49\xb6\x7b\xae\xaf\xa7\x31\xa0\xdb\x5f\x03\x8d\xf9\xb8\x66\x3e\x23\x30\xeb\x9f\x01\x38\xee\xff\x91\x3d\xa0\xa6\x3b\x02\x1f\x31\x45\x77\x01\x97\x27\xe6\xe6\x56\x52\xc2\xdc\x05\x28\xc2\xd5\x8d\x42\x49\xf7\x3d\x59\x6a\x67\xea\xc5\xcb\x11\x84\x71\xdd\x00\x5c\x2e\x64\x4c\x98\xdb\xc1\xe5\x43\x5e\x62\xad\xc4\x5a\x4d\x59\xeb\x82\x33\x41\x42\x63\xf4\xd5\x79\xab\x79\x64\x32\x76\x23\xc9\xf1\x73\xc3\x71\x7e\xf9\xf2\xe9\xa3\xa7\x8a\x37\x0b\xdd\x67\xb3\xda\x8c\x1a\x0e\xd6\x32\xa9\x65\x15\x93\x10\xa9\x70\xd6\x64\x1b\x0a\xf1\x27\x85\x5c\x46\xa8\xad\xb8\xa8\xe1\xc2\x80\x61\x05\x4e\xb7\xa2\x9c\xd1\xe8\x3e\x23\x37\x16\x7b\x3e\xa4\x31\xe1\xf4\x07\xde\x8b\x38\x26\xe0\x7d\x15\x9a\xb0\x8f\x69\x9c\xe7\x5d\x36\x7f\x40\x4b\x8c\xb1\x90\x19\xa4\x0a\xc3\x4a\x0c\xcc\xb6\x99\xc6\x79\x9f\x26\x7f\xba\xa6\xbb\x4c\xa3\x2a\x35\x15\xd7\x79\x7e\x39\x1f\x9e\xcf\x85\x93\xc2\x53\x33\xde\x02\x9e\xf2\x34\x86\x60\x05\xde\xd7\x2c\xc1\x81\x09\x2c\x3d\xe3\xfe\x5e\xf5\x7e\x58\x2f\x37\xf6\x41\x92\xaa\xc3\xec\x7c\x36\x14\x79\xde\xb7\x0a\xfe\x78\x20\xe5\x21\x9e\xe0\x69\x31\xb4\xd5\x5c\x88\xea\xe7\x2a\x45\x58\x1b\xad\x87\x1e\xc6\xfb\xa3\x78\x68\xe4\x03\x79\xd8\x4c\x50\x55\x76\xfa\x64\xca\x34\x14\xbb\xd4\x64\xdb\x1c\x63\xfe\x67\x13\x7f\x97\xbd\x0f\x67\x4e\xf1\x42\xb8\x21\x12\xc9\xcd\xd2\x99\x9b\xee\x7b\xc1\xcd\x52\xcc\x9c\x17\xa1\xd3\x50\x54\x6e\xaa\x38\xfb\x4c\xd1\x30\x72\x3c\xda\x13\xcd\xac\x1c\x65\xd1\x98\x64\x4f\xf9\xeb\x93\x3f\x51\x87\x3f\x59\x87\xdf\xa3\xc3\x6f\xe8\xf0\xe7\xaf\xda\x09\xd9\x89\xb4\x38\xfa\x9d\xf3\x56\xea\xda\xa5\x84\x5c\x4b\x8a\xaa\x28\xa7\x0f\xc8\x3f\xe0\x03\xb2\x7b\x13\xd8\x28\x2a\x43\x56\xe1\x56\xb0\xde\xf4\xd1\x99\xfe\xcc\xb0\x54\x8c\x8d\xf8\xf2\x49\xb9\xa0\xdd\x37\x6c\x00\xe7\xb3\x25\xf0\x7e\xc5\x2c\xbf\x7a\x41\x98\xdd\xa5\x1a\x90\x72\x77\xb4\x41\x3c\x8d\x1b\x90\xf6\x9e\xe8\x2d\x9d\xea\xbe\x48\xd2\xba\xae\xb2\x0d\xac\x2a\xb9\x63\xf5\xb6\x17\x12\x66\x3a\x4b\x80\x72\xcb\x70\xfd\x7d\xda\x4e\x99\x1d\x44\x67\xc9\xe6\x55\x07\x25\xc9\xf1\x8e\xc8\x21\x5b\xde\x56\x5d\xf3\x6e\x64\xdd\xd7\xb5\xe6\xbb\x29\x61\x5d\x7b\x1e\x0a\x2b\xa6\x4b\xcd\x64\xcb\x39\x0d\x7d\x8e\xdb\x2a\xa8\x53\xb8\xa6\x9b\xee\x71\xf9\x22\x7b\xcc\x69\xec\x42\xf2\x34\xee\x39\x70\x5f\x66\xf0\x38\x43\x51\x3d\x63\x1c\x3d\x1e\xe5\xbc\x71\xe0\x99\x9d\x8b\x7d\x3e\x41\xc3\x9f\x89\xcf\x7b\xd7\xe4\x11\xef\xda\x12\xb9\xbc\x21\x28\xc9\x8d\x21\xd7\x59\x32\xe6\x1b\x50\x7b\xc7\xda\x89\xb3\x3b\x22\x97\x65\xd0\xa6\x65\x22\xe5\xb0\x8b\x46\x36\x06\x0a\xe6\x11\x43\xdb\x12\xe9\xff\x55\x71\x92\xf8\x23\xe2\xfc\x86\x38\xbf\x41\x52\x9f\x81\x5e\x3f\xaf\x7e\x9f\x7a\xd2\xfa\xa5\xaa\x55\xea\x7d\x5f\x58\x5b\x22\xff\xfd\xb8\x1a\x82\xff\x0d\x1f\x57\xa7\xdb\x13\x55\xef\xdf\x06\xe0\x6c\x69\x14\xa1\x04\x7d\x20\xfc\xea\x3c\x9d\xd5\x18\x7b\x00\x74\xfe\x09\xdf\x67\x05\xc5\xe2\xb1\xb3\x7f\x5d\xf7\x7f\x04\x00\x00\xff\xff\x8a\x03\xd5\xa9\xa0\x15\x00\x00")

func chartjsHtmlBytes() ([]byte, error) {
	return bindataRead(
		_chartjsHtml,
		"chartjs.html",
	)
}

func chartjsHtml() (*asset, error) {
	bytes, err := chartjsHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "chartjs.html", size: 5536, mode: os.FileMode(420), modTime: time.Unix(1710821830, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _footerHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\xcb\xcf\x2f\x49\x2d\x52\x48\xce\x49\x2c\x2e\xb6\x55\xca\x4d\xcc\xcc\xd3\x85\x08\x29\xd9\x71\x29\x28\x28\x28\xd8\x24\xa7\xe6\x95\xa4\x16\xd9\xd9\x14\x97\x14\xe5\xe7\xa5\xdb\x39\xe7\x17\x54\x16\x65\xa6\x67\x94\x28\x1c\x5a\xa9\x60\x64\x60\x68\xae\x60\xa3\x0f\x95\x52\x70\xcc\xc9\x51\x00\xcb\x15\x2b\x14\xa5\x16\xa7\x16\x95\xa5\xa6\xe8\xd9\xe8\x43\x0d\xe0\xb2\xd1\x87\x18\x6c\x07\x08\x00\x00\xff\xff\x75\xbb\x17\xa1\x74\x00\x00\x00")

func footerHtmlBytes() ([]byte, error) {
	return bindataRead(
		_footerHtml,
		"footer.html",
	)
}

func footerHtml() (*asset, error) {
	bytes, err := footerHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "footer.html", size: 116, mode: os.FileMode(420), modTime: time.Unix(1636423570, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _headerHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x93\xcf\x6e\xdb\x30\x0c\xc6\xef\x7d\x0a\x4e\x3b\x2b\x4a\x07\xb4\x87\xc1\x31\x30\x6c\x87\x15\x18\x7a\xe8\xb6\x07\xa0\x2d\x45\xe1\x20\x53\x81\x24\x7b\x29\x82\xbc\xfb\xe0\xbf\x71\x9c\x66\xbd\x24\x30\x45\x7e\x1f\xfd\xfd\xac\x6c\x67\x50\x9b\x00\xa5\xc3\x18\x37\xa2\x42\x62\xd9\x97\x44\x7e\x07\x00\x90\x7d\x90\x12\x7e\x78\xeb\x41\xca\xa1\x82\xb0\x0b\x66\xbb\x11\x4a\x8c\x63\xce\x5b\x3f\xf4\x4f\x33\x15\x31\x41\x7b\x00\x5b\x1f\x20\x92\x36\x05\x86\xbe\xfa\xb0\x3e\x3c\xac\x61\x4f\x07\xe3\xe2\x24\xdb\x0d\xc6\x3d\xf2\x5c\x54\xb6\xfd\x22\xcf\x8a\xfc\xe5\xdb\x4b\xa6\x8a\x3c\x53\x6d\xcb\xc2\x6a\x72\x09\xc6\xd6\x0e\x03\xc4\x84\xc9\x00\xb2\x86\xca\x17\xe4\x0c\x68\xd3\x50\x69\xde\x33\x73\xf6\x96\x55\xa6\x70\x16\xc7\xf7\x3e\xb3\x67\x6c\x0a\x0c\x9f\x21\xa6\x57\x67\xa0\x44\x86\xc2\xc0\xd6\xd7\xac\x81\x18\xfa\x14\x57\xce\xc4\xb3\x6d\xc6\xd8\x8c\x8e\xdc\x4d\x43\xff\x27\xdb\x8d\xa9\x94\xc9\xef\x05\x04\xef\x4c\x77\x4e\x16\x13\x79\x5e\x26\xfb\x73\xc8\x32\x79\x6b\x9d\x81\xa2\x4e\xc9\xf3\xc5\x9b\x4d\x84\x88\x63\x42\x2e\x8d\x3a\x1e\x57\x5f\xeb\x10\x0c\xa7\xa7\xa1\x74\x3a\x4d\xf8\x06\x38\xb2\x17\x9c\xd9\x5d\xc5\x14\x83\xf4\xec\x5e\x45\xfe\xab\xf7\x3e\x2f\xb9\x04\x73\x35\xba\x23\xad\x0d\xcb\x43\x14\x39\x1c\x8f\xb0\x5c\x07\x4e\xa7\x2b\xb6\x63\xe8\xdd\x83\xa6\x45\x74\xb2\xac\x63\xf2\x95\xac\x0c\xd7\xcb\x9d\x6b\x37\xeb\x1d\x33\x66\x6c\x16\x7d\x53\xa4\xbf\xa3\x09\xf0\xa5\x2c\x7d\xcd\xe9\x16\x50\x1d\xfc\x5e\xfb\xbf\x7c\x89\xf4\x42\xca\xd1\x68\x3b\x36\x43\xdd\x2a\xb7\x3f\x6f\xed\x79\x05\xec\xa3\x58\x0a\x8c\x50\x40\x63\xc2\xe1\xe1\x7c\x7a\x43\xaf\xd3\xa4\xca\x42\x0c\xe5\x46\xa8\xfe\xe3\x52\x9a\x62\x52\x54\x59\xd5\xae\xf3\x49\xde\x3f\xae\x0f\xf7\x8f\xeb\xd5\x9f\xbd\x9d\x5c\xbb\x45\xa9\x42\x6b\x04\xa0\x4b\x1b\xd1\x05\xf3\xd4\x15\xfe\x63\x75\x83\x73\x8b\xb9\x15\x78\xc6\xea\x2d\xbe\x17\x12\x73\xd6\xe7\xa2\xa3\x05\x58\x55\xbb\xf9\x07\xa2\xa9\x19\x2f\x28\x63\x93\xdf\x65\xaa\xbf\x76\xf9\xbf\x00\x00\x00\xff\xff\xb9\x8a\x81\x0a\xd8\x04\x00\x00")

func headerHtmlBytes() ([]byte, error) {
	return bindataRead(
		_headerHtml,
		"header.html",
	)
}

func headerHtml() (*asset, error) {
	bytes, err := headerHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "header.html", size: 1240, mode: os.FileMode(420), modTime: time.Unix(1636423570, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _revelHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x57\xcd\x6e\xe3\x36\x10\xbe\xf7\x29\x06\x6a\x7a\x2b\x57\x36\x8a\xa2\x80\x23\xfb\xd0\x45\x7b\xc9\xa2\xdb\xc3\xde\x0b\x4a\x1c\x5b\x6c\x28\x4a\x20\x29\x27\xaa\xe0\x77\x2f\x48\x4b\x8e\xa4\x58\x7f\x8e\x1b\x24\x68\x7c\x50\x22\x71\xf8\xcd\x2f\xbf\x19\x06\x8c\xef\x21\x12\x54\xeb\xb5\x17\xa5\xd2\xa0\x34\xe4\x41\xd1\x2c\x43\xe5\x81\x36\x85\xc0\xb5\x97\x70\x49\x62\xe4\xbb\xd8\xac\x60\xb9\x58\x64\x8f\xb7\x50\xbf\xd2\xdc\xa4\xb7\x90\xee\x51\x6d\x45\xfa\xb0\x82\x98\x33\x86\xd2\xdb\x7c\x07\x00\xd0\xc6\x16\x24\x61\xe4\xa7\x6a\xc9\x2d\x6b\x8c\x0c\x4f\x65\x57\x7d\x8c\x94\xa1\x6a\x08\x76\xb1\xc2\xf4\xb1\xb3\x7a\x46\x82\x84\x29\x2b\xce\x88\x9d\x44\x39\x5b\x7b\x11\x95\x7b\xaa\x49\x9c\x0a\xd6\xf0\xf7\x81\x33\x13\xaf\x96\x8b\xc5\x0f\x3d\xfb\x1d\xc6\x71\xef\x11\x26\xa6\xca\x10\xaa\x90\x92\x85\x07\x7e\x8f\x52\x9f\xf1\xfd\x19\xb3\x9f\x7f\xee\x7c\x0a\xfc\x2a\x50\x55\x54\x9f\x56\xdf\x59\x80\xc9\xf2\x2a\x21\x5e\x7e\x84\xd8\x05\x07\xa5\x41\xb5\x09\xb4\x51\xa9\xdc\x6d\x04\x4a\x10\xb8\x47\x01\x51\x9a\x4b\x13\xf8\xd5\xf7\xc0\xaf\x05\x43\x35\x90\xac\x4a\xa9\xa4\x7b\x62\x68\xa8\x49\x94\x6b\x93\x26\x43\xc9\xc9\x85\x4b\x8c\x40\x49\x64\x9e\x38\x8f\xbd\x06\x0c\xd4\x50\x03\x18\xf6\x57\x96\x8a\xca\x1d\xc2\x8d\x29\x32\xfc\x11\x6e\x50\x1a\xc5\x51\xc3\x6a\x0d\x9f\xbe\xa0\xfc\x62\x5d\xfa\x6c\x3d\x3a\x1c\x06\x71\x02\xc1\xa1\x2c\xf9\x16\xb8\xfe\x9d\x2b\x6d\x0e\x87\xda\x18\x1a\x19\xbe\x47\x0f\xca\x12\x25\x3b\x1c\x36\x01\x85\x58\xe1\x76\xed\x7d\x1f\x52\xb5\x20\x65\xe9\x54\x1f\x0e\x1e\x30\x6a\x28\x31\xe9\x6e\x67\x4b\xd4\xd0\xd0\xdb\x9c\x16\x03\x9f\x6e\x02\x5f\xf0\x31\x67\x9c\x0a\x28\xcb\x1b\x99\x0b\x01\xab\x75\x24\x90\xaa\xca\xa0\xfe\x58\xfa\xb9\x18\x88\x74\x23\x3f\x86\x86\xa4\xaa\x31\xc8\xa8\x44\x31\x54\x22\x57\x0f\x71\x7d\xaa\xbb\x71\x6b\x18\x67\x8d\xea\xe4\xe1\x18\xff\x2a\x34\x23\xa6\x3e\x57\x63\xbd\xa5\x5c\xa2\x6a\x2a\x9c\xc3\x22\x2d\xe8\x06\xa3\x38\x74\x8a\x8a\x36\x81\x7b\x88\xa5\x85\x71\x9e\x64\x66\x8a\x5c\x5c\x29\xbd\xc8\x1f\xe4\x67\xc9\x2f\xc1\x04\x72\x8d\xec\x55\xf9\x2f\xc1\x77\xcb\x7f\xcb\x0f\xfe\xbb\x88\xff\x96\xaf\xc3\x7f\xcb\xff\x94\xff\x96\x1f\xfc\x37\x97\xff\x7e\x79\x93\xfc\x67\xd2\xcc\xde\xcc\x40\x50\xb5\x43\x6d\xe0\x1e\x0b\x3d\x87\x01\x0d\x0d\x05\x36\xca\x58\x20\xb8\xa7\xad\x3e\x86\x52\x23\xab\xde\x63\x7b\xd5\x03\x9d\x2a\xf7\xfa\x54\x88\xa9\x62\x24\x54\x48\xef\x57\xee\x49\xa8\x10\xb7\xe0\xbe\xda\xcb\x64\xe3\xe3\x10\xa3\x1a\x1b\xb8\x91\xca\x32\x3d\x2e\xb4\x85\x58\xed\x8b\xb5\xd4\x99\xfa\x17\x15\x59\x4c\xbd\x0d\xdc\x61\x01\x81\x6f\x46\xf4\x8c\xc2\x7c\x2b\x32\xbc\x06\xce\xaf\x85\x41\xfd\x12\x20\x99\x27\xa8\x78\xe4\x6d\xe0\x8f\x3c\xf9\xba\xfd\x4d\x60\x62\xc9\x71\x1c\x31\xf0\x87\x62\x19\xf8\x23\xd9\x08\x8c\xad\xcf\x89\xc4\x6b\x19\xb7\x38\xf2\xed\xb1\x42\xef\xb0\xd0\x63\x6c\x3b\x31\xd5\xb6\x47\x39\xfc\x4f\x77\x58\xd8\x46\x35\x31\x94\x4f\xfb\xbe\x55\x1d\x6e\xfa\xc6\x38\x4f\xa8\xe4\xff\xe0\x31\x7b\x15\x8c\x7b\xb9\x0c\xe7\x73\x9a\x24\xb4\xc6\x39\xe5\x71\x0a\xd6\x70\x16\xe1\x89\x69\x87\x32\xdd\x9f\xc9\xc0\x77\x55\x76\x3d\x5a\xed\xe3\xd5\x9f\xdf\x24\xaf\xba\xab\x34\x84\x85\x25\x54\xc8\x14\x6e\xf9\xe3\xab\x0c\x96\x67\xc6\x48\xf7\xcf\xdf\xb9\x36\x7c\xcb\x91\xbd\x64\xe2\x39\x9d\xc0\x3f\x9d\x43\x38\x7a\x10\x2f\x98\x2c\xad\xde\xff\xfb\x64\x39\x3b\xce\xf5\xdc\xd7\x0d\xde\x95\xc7\xcb\x37\xd1\xec\xdb\x16\x8d\x37\xfe\xb6\xfc\x84\xce\xd0\xde\x30\x3c\x10\x1c\x13\x34\xb1\x0b\x4f\x06\x9e\xd3\xda\x47\x41\x5b\x6d\x3e\x44\xf5\x75\x3b\x7d\x90\x39\x01\x8f\xf6\x8a\x86\xe4\xf4\x94\x4c\x99\x04\x9a\xbf\xe7\x53\x41\x7d\x6e\x46\x0e\x48\x5b\xe9\xfc\x1a\xb8\x68\x52\x78\x8e\xf1\xf2\xe6\xdf\x8f\xd9\x1d\x04\xe6\x22\x4e\x4f\x31\x4c\x18\x0d\x3a\xc8\xd3\xd2\xdc\x3b\x32\xb4\x85\xde\xff\x0d\xf1\xf8\xe7\xdf\x00\x00\x00\xff\xff\x8c\xd4\xc8\x9c\x41\x1b\x00\x00")

func revelHtmlBytes() ([]byte, error) {
	return bindataRead(
		_revelHtml,
		"revel.html",
	)
}

func revelHtml() (*asset, error) {
	bytes, err := revelHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "revel.html", size: 6977, mode: os.FileMode(420), modTime: time.Unix(1710821801, 0)}
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
	"aside.html": asideHtml,
	"base.html": baseHtml,
	"chartjs.html": chartjsHtml,
	"footer.html": footerHtml,
	"header.html": headerHtml,
	"revel.html": revelHtml,
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
	"aside.html": &bintree{asideHtml, map[string]*bintree{}},
	"base.html": &bintree{baseHtml, map[string]*bintree{}},
	"chartjs.html": &bintree{chartjsHtml, map[string]*bintree{}},
	"footer.html": &bintree{footerHtml, map[string]*bintree{}},
	"header.html": &bintree{headerHtml, map[string]*bintree{}},
	"revel.html": &bintree{revelHtml, map[string]*bintree{}},
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

