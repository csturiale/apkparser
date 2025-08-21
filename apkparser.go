// Package apkparser parses AndroidManifest.xml and resources.arsc from Android APKs.
package apkparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
)

type ApkParser struct {
	apkPath string
	zip     *ZipReader

	encoder   ManifestEncoder
	resources *ResourceTable
}

// Calls ParseApkReader
func ParseApk(path string) (ApkInfo, error) {
	f, zipErr := os.Open(path)
	if zipErr != nil {
		return ApkInfo{}, zipErr
	}
	defer f.Close()
	return ParseApkReader(f)
}

// Parse APK's Manifest, including resolving refences to resource values.
// encoder expects an XML encoder instance, like Encoder from encoding/xml package.
//
// zipErr != nil means the APK couldn't be opened. The manifest will be parsed
// even when resourcesErr != nil, just without reference resolving.
func ParseApkReader(r io.ReadSeeker) (ApkInfo, error) {
	zip, zipErr := OpenZipReader(r)
	if zipErr != nil {
		return ApkInfo{}, zipErr
	}
	defer zip.Close()

	return ParseApkWithZip(zip)
}

// Parse APK's Manifest, including resolving refences to resource values.
// encoder expects an XML encoder instance, like Encoder from encoding/xml package.
//
// Use this if you already opened the zip with OpenZip or OpenZipReader before.
// This method will not Close() the zip.
//
// The manifest will be parsed even when resourcesErr != nil, just without reference resolving.
func ParseApkWithZip(zip *ZipReader) (ApkInfo, error) {
	var apkInfo ApkInfo
	buf := new(bytes.Buffer)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "\t")

	p := ApkParser{
		zip:     zip,
		encoder: enc,
	}

	resourcesErr := p.parseResources()
	if resourcesErr != nil {
		return apkInfo, resourcesErr
	}
	manifestErr := p.ParseXml("AndroidManifest.xml")
	if manifestErr != nil {
		return apkInfo, manifestErr
	}

	var manifest Manifest
	_ = xml.Unmarshal(buf.Bytes(), &manifest)
	iconPath := manifest.App.Icon
	icon, _ := p.ParseIcon(iconPath)

	apkInfo.Manifest = manifest
	apkInfo.Icon = icon

	return apkInfo, nil
}

// Prepare the ApkParser instance, load resources if possible.
// encoder expects an XML encoder instance, like Encoder from encoding/xml package.
//
// This method will not Close() the zip, you are still the owner.
func NewParser(zip *ZipReader, encoder ManifestEncoder) (parser *ApkParser, resourcesErr error) {
	parser = &ApkParser{
		zip:     zip,
		encoder: encoder,
	}
	resourcesErr = parser.parseResources()
	return
}

func (p *ApkParser) parseResources() (err error) {
	if p.resources != nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic: %v\n%s", r, string(debug.Stack()))
		}
	}()

	resourcesFile := p.zip.File["resources.arsc"]
	if resourcesFile == nil {
		return os.ErrNotExist
	}

	if err := resourcesFile.Open(); err != nil {
		return fmt.Errorf("Failed to open resources.arsc: %s", err.Error())
	}
	defer resourcesFile.Close()

	p.resources, err = ParseResourceTable(resourcesFile)
	return
}

func (p *ApkParser) ParseXml(name string) error {
	file := p.zip.File[name]
	if file == nil {
		return fmt.Errorf("Failed to find %s in APK!", name)
	}

	if err := file.Open(); err != nil {
		return err
	}
	defer file.Close()

	var lastErr error
	for file.Next() {
		if err := ParseXml(file, p.encoder, p.resources); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	if lastErr == ErrPlainTextManifest {
		return lastErr
	}

	return fmt.Errorf("Failed to parse %s, last error: %v", name, lastErr)
}

func (p *ApkParser) ParseIcon(name string) (image.Image, error) {
	file := p.zip.File[name]

	if file == nil {
		return nil, fmt.Errorf("Failed to find %s in APK!", name)
	}

	if err := file.Open(); err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	resourceIconId, err := p.checkAdaptiveIcon(name, b)
	if err != nil {
		return nil, err
	}
	if resourceIconId > 0 {
		png, err := p.resources.GetIconPng(resourceIconId)
		if err != nil {
			return nil, err
		}
		data, err := png.GetValue().Data()
		if err != nil {
			return nil, err
		}
		return p.ParseIcon(data.(string))
	}

	icon, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return icon, nil
}

func (p *ApkParser) checkAdaptiveIcon(name string, b []byte) (uint32, error) {
	var resourceID uint32
	resourceID = 0
	if filepath.Ext(name) == ".xml" {
		buf := new(bytes.Buffer)
		enc := xml.NewEncoder(buf)
		enc.Indent("", "    ")
		if err := ParseXml(bytes.NewReader(b), enc, nil); err != nil {
			return resourceID, err
		}
		var ai AdaptiveIcon
		err := xml.Unmarshal(buf.Bytes(), &ai)
		resourceID, err = ParseResID(ai.Foreground.Drawable)
		if err != nil {
			return resourceID, err
		}
	}
	return resourceID, nil
}

func ParseResID(s string) (uint32, error) {
	if strings.Trim(s, " ") == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
