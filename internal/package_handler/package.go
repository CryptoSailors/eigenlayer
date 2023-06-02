package package_handler

import (
	"errors"
	"fmt"
	"path/filepath"
)

const (
	pkgDirName       = "pkg"
	checksumFileName = "checksum.txt"
)

// PackageHandler is used to interact with an AVS node software package at the given
// path.
type PackageHandler struct {
	path string
}

// NewPackageHandler creates a new PackageHandler instance for the given package path.
func NewPackageHandler(path string) *PackageHandler {
	return &PackageHandler{path: path}
}

// Check validates a package. It returns an error if the package is invalid.
// It checks the existence of some required files and directories and computes the
// checksums comparing them with the ones listed in the checksum.txt file.
func (p *PackageHandler) Check() error {
	if err := checkPackageDirExist(p.path, pkgDirName); err != nil {
		return err
	}
	err := checkPackageFileExist(p.path, checksumFileName)
	if err != nil {
		var fileNotFoundErr PackageFileNotFoundError
		if errors.As(err, &fileNotFoundErr) {
			return nil
		}
		return err
	} else {
		return p.checkSum()
	}
}

func (p *PackageHandler) checkSum() error {
	currentChecksums, err := parseChecksumFile(filepath.Join(p.path, checksumFileName))
	if err != nil {
		return err
	}
	computedChecksums, err := packageHashes(p.path)
	if err != nil {
		return err
	}
	if len(currentChecksums) != len(computedChecksums) {
		return fmt.Errorf("%w: expected %d files, got %d", ErrInvalidChecksum, len(currentChecksums), len(computedChecksums))
	}
	for file, hash := range currentChecksums {
		if computedChecksums[file] != hash {
			return fmt.Errorf("%w: checksum mismatch for file %s, expected %s, got %s", ErrInvalidChecksum, file, hash, computedChecksums[file])
		}
	}
	return nil
}