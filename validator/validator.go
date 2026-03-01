package validator

import (
	"errors"
	"fmt"
	"os"
)

var (
	errPathNotExist = errors.New("path doesn't exist")
	errNotDirectory = errors.New("not a directory")
	errNotFile      = errors.New("not a file")
)

func ValidatePkg(pkgPath string) error {
	info, err := os.Stat(pkgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errPathNotExist
		}
		return fmt.Errorf("path stat: %w", err)
	}

	if !info.IsDir() {
		return errNotDirectory
	}

	return nil
}

func ValidateMainFile(mainPath string) error {
	info, err := os.Stat(mainPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errPathNotExist
		}
		return fmt.Errorf("path stat: %w", err)
	}

	if info.IsDir() {
		return errNotFile
	}

	return nil
}

func ValidateAPIFile(apiPath string) error {
	info, err := os.Stat(apiPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("path stat: %w", err)
	}

	if info.IsDir() {
		return errNotFile
	}

	if err = os.Remove(apiPath); err != nil {
		return fmt.Errorf("removing: %w", err)
	}

	return nil
}
