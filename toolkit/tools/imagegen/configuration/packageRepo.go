// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

// Parser for the image builder's PackageRepo configuration schemas.

package configuration

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// PackageRepo defines the RPM repo to pull packages from during the installation
// or after the system is installed. The "Install" option indicates that the provided
// repository configuration will be saved in the installed system if specified, and only
// available during the installation process if not
type PackageRepo struct {
	Name    string `json:"Name"`
	BaseUrl string `json:"BaseUrl"`
	Install bool   `json:"Install"`
}

// UnmarshalJSON Unmarshals a PackageRepo entry
func (p *PackageRepo) UnmarshalJSON(b []byte) (err error) {
	// Use an intermediate type which will use the default JSON unmarshal implementation
	type IntermediateTypePackageRepo PackageRepo
	err = json.Unmarshal(b, (*IntermediateTypePackageRepo)(p))
	if err != nil {
		return fmt.Errorf("failed to parse [PackageRepo]: %w", err)
	}

	// Now validate the resulting unmarshaled object
	err = p.IsValid()
	if err != nil {
		return fmt.Errorf("failed to parse [PackageRepo]: %w", err)
	}
	return
}

// IsValid returns an error if the PackageRepo struct is not valid
func (p *PackageRepo) IsValid() (err error) {
	err = p.NameIsValid()
	if err != nil {
		return
	}

	err = p.RepoUrlIsValid()
	if err != nil {
		return
	}

	return
}

// NameIsValid returns an error if the package repo name is empty
func (p *PackageRepo) NameIsValid() (err error) {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("invalid value for package repo name (%s), name cannot be empty", p.Name)
	}
	return
}

// RepoUrlIsValid returns an error if the package url is invalid
func (p *PackageRepo) RepoUrlIsValid() (err error) {
	if strings.TrimSpace(p.BaseUrl) == "" {
		return fmt.Errorf("invalid value for package repo URL (%s), URL cannot be empty", p.BaseUrl)
	}

	_, err = url.ParseRequestURI(p.BaseUrl)
	if err != nil {
		return
	}
	return
}
