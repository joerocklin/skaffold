/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package buildpacks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
)

// For testing
var (
	Validate = validate
)

// Name is the name of the Buildpack builder
var Name = "Buildpacks"

// ArtifactConfig holds information about a Buildpack project
type ArtifactConfig struct {
	File    string `json:"path,omitempty"`
	Builder string `json:"builder,omitempty"`
}

// Name returns the name of the builder
func (c ArtifactConfig) Name() string {
	return Name
}

// Describe returns the initBuilder's string representation, used when prompting the user to choose a builder.
func (c ArtifactConfig) Describe() string {
	return fmt.Sprintf("%s (%s)", c.Name(), c.File)
}

// CreateArtifact creates an Artifact to be included in the generated Build Config
func (c ArtifactConfig) UpdateArtifact(a *latest.Artifact) {
	a.ArtifactType = latest.ArtifactType{
		BuildpackArtifact: &latest.BuildpackArtifact{
			Builder: c.Builder,
		},
	}
}

// ConfiguredImage returns the target image configured by the builder, or empty string if no image is configured
func (c ArtifactConfig) ConfiguredImage() string {
	// Target image is not configured in buildpacks
	return ""
}

// Path returns the path to the build definition
func (c ArtifactConfig) Path() string {
	return c.File
}

// validate checks if a file is a valid Buildpack configuration.
func validate(path string) bool {
	switch filepath.Base(path) {
	case "package.json":
		return !hasParent(path, "node_modules")
	case "go.mod":
		return !hasParent(path, "vendor")
	default:
		return false
	}
}

func hasParent(path, parent string) bool {
	for _, p := range strings.Split(path, string(os.PathSeparator)) {
		if p == parent {
			return true
		}
	}

	return false
}
