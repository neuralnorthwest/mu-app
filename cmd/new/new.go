// Copyright 2023 Scott M. Long
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// newImpl is the implementation for the new command.
type newImpl struct {
	// Name is the name of the new project.
	Name string
	// Path is the path to the new project.
	Path string
	// FromRepo is the repository to use as the template.
	FromRepo string
	// FromRef is the reference to use from the repository.
	FromRef string
	// Target is the GitHub target repository to use, in the form of owner/repo.
	Target string
}

// newCommand returns a new command for creating a new project.
func newCommand() *cobra.Command {
	n := &newImpl{}
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Creates a new project",
		Long:  "Creates a new project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n.Name = args[0]
			return n.run()
		},
	}
	cmd.Flags().StringVar(&n.Path, "path", "", "The path to the new project")
	cmd.Flags().StringVar(&n.FromRepo, "from-repo", "https://github.com/neuralnorthwest/mu-app", "The repository to use as the template")
	cmd.Flags().StringVar(&n.FromRef, "from-ref", "main", "The reference to use from the repository")
	cmd.Flags().StringVar(&n.Target, "target", "", "The GitHub target repository to use, in the form of owner/repo")
	return cmd
}

// runNew creates a new project.
func (n *newImpl) run() error {
	if n.validateName() != nil {
		return fmt.Errorf("invalid name %q", n.Name)
	}
	if n.Path == "" {
		n.Path = n.Name
	}
	n.Path = filepath.Clean(n.Path)
	if _, err := os.Stat(n.Path); err == nil {
		return fmt.Errorf("path %q already exists", n.Path)
	}
	// create parent of path if it doesn't exist using dirname
	if err := os.MkdirAll(filepath.Dir(n.Path), 0755); err != nil {
		return fmt.Errorf("unable to create directory %q: %w", filepath.Dir(n.Path), err)
	}
	// clone the repo
	if err := n.cloneRepo(); err != nil {
		return fmt.Errorf("unable to clone repo: %w", err)
	}
	// remove the .git directory
	if err := os.RemoveAll(filepath.Join(n.Path, ".git")); err != nil {
		return fmt.Errorf("unable to remove .git directory: %w", err)
	}
	// rename the files
	if err := n.renameFiles(); err != nil {
		return fmt.Errorf("unable to rename files: %w", err)
	}
	// rename the contents
	if err := n.renameContents(); err != nil {
		return fmt.Errorf("unable to rename contents: %w", err)
	}
	// process the README.md
	if err := n.processReadme(); err != nil {
		return fmt.Errorf("unable to process README.md: %w", err)
	}
	return nil
}

// validateName validates the name.
func (n *newImpl) validateName() error {
	if n.Name == "" {
		return fmt.Errorf("name is required")
	}
	// must be lowercase
	if n.Name != strings.ToLower(n.Name) {
		return fmt.Errorf("name must be lowercase")
	}
	// must begin with a letter
	if !regexp.MustCompile(`^[a-z]`).MatchString(n.Name) {
		return fmt.Errorf("name must begin with a letter")
	}
	// must end with a letter or number
	if !regexp.MustCompile(`[a-z0-9]$`).MatchString(n.Name) {
		return fmt.Errorf("name must end with a letter or number")
	}
	// must contain only letters, numbers, and dashes
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(n.Name) {
		return fmt.Errorf("name must contain only letters, numbers, and dashes")
	}
	// must not contain consecutive dashes
	if regexp.MustCompile(`--`).MatchString(n.Name) {
		return fmt.Errorf("name must not contain consecutive dashes")
	}
	return nil
}

// cloneRepo clones the repository.
func (n *newImpl) cloneRepo() error {
	return exec.Command("git", "clone", "--depth", "1", "--branch", n.FromRef, n.FromRepo, n.Path).Run()
}

// renameFiles renames the files and directories.
func (n *newImpl) renameFiles() error {
	// get a list of all files and directories
	all := []string{}
	if err := filepath.WalkDir(n.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		if base == "mu-app" || base == "mu_app" || base == "muapp" || base == "muApp" {
			all = append(all, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("unable to walk directory %q: %w", n.Path, err)
	}
	// sort the paths by decreasing depth
	sort.Slice(all, func(i, j int) bool {
		return strings.Count(all[i], string(filepath.Separator)) > strings.Count(all[j], string(filepath.Separator))
	})
	// get the names to use in renaming
	dashName, underName, packName, camelName := n.formatNames()
	// rename the files and directories
	for _, path := range all {
		// get the new path
		newPath := n.rename(path, dashName, underName, packName, camelName)
		// rename the file or directory
		if err := os.Rename(path, newPath); err != nil {
			return fmt.Errorf("unable to rename %q to %q: %w", path, newPath, err)
		}
	}
	return nil
}

// renameContents renames the contents of the files.
func (n *newImpl) renameContents() error {
	// get the names to use in renaming
	dashName, underName, packName, camelName := n.formatNames()
	if err := filepath.WalkDir(n.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// read the file
		contents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", path, err)
		}
		// rename the contents
		newContents := n.rename(string(contents), dashName, underName, packName, camelName)
		// write the file
		if err := os.WriteFile(path, []byte(newContents), 0644); err != nil {
			return fmt.Errorf("unable to write file %q: %w", path, err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("unable to walk directory %q: %w", n.Path, err)
	}
	return nil
}

// processReadme processes the README.md file.
func (n *newImpl) processReadme() error {
	// load the README.md
	path := filepath.Join(n.Path, "README.md")
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file %q: %w", path, err)
	}
	// split into lines
	lines := strings.Split(string(contents), "\n")
	outLines := []string{}
	// state machine:
	// 0: copying lines, have not reached a badge
	// 1: have reached a badge and are skipping lines (target is "")
	// 2: have reached a badge and are copying lines (target is not "")
	// 3: have passed the badges and are done
	state := 0
COPY_LOOP:
	for _, line := range lines {
		switch state {
		case 0:
			// copying lines, have not reached a badge
			if strings.HasPrefix(line, "![") && n.Target == "" {
				state = 1
			} else if strings.HasPrefix(line, "![") && n.Target != "" {
				outLines = append(outLines, strings.ReplaceAll(line, "neuralnorthwest/mu-app", n.Target))
				state = 2
			} else {
				outLines = append(outLines, line)
			}
		case 1:
			// have reached a badge and are skipping lines (target is "")
			if !strings.HasPrefix(line, "![") {
				state = 3
			}
		case 2:
			// have reached a badge and are copying lines (target is not "")
			if !strings.HasPrefix(line, "![") {
				outLines = append(outLines, "")
				state = 3
			} else {
				outLines = append(outLines, strings.ReplaceAll(line, "neuralnorthwest/mu-app", n.Target))
			}
		case 3:
			// have passed the badges and are done
			break COPY_LOOP
		}
	}
	outLines = append(outLines, "PLACEHOLDER")
	// write the file
	if err := os.WriteFile(path, []byte(strings.Join(outLines, "\n")), 0644); err != nil {
		return fmt.Errorf("unable to write file %q: %w", path, err)
	}
	return nil
}

// formatNames formats the names for use in renaming files and directories.
func (n *newImpl) formatNames() (dashName, underName, packName, camelName string) {
	dashName = n.Name
	dashComps := strings.Split(dashName, "-")
	underName = strings.Join(dashComps, "_")
	packName = strings.Join(dashComps, "")
	camelName = strings.ToLower(dashComps[0])
	titler := cases.Title(language.AmericanEnglish)
	if len(dashComps) > 1 {
		for _, comp := range dashComps[1:] {
			camelName += titler.String(comp)
		}
	}
	return
}

// rename renames a string with the given replacements.
func (n *newImpl) rename(s, dashName, underName, packName, camelName string) string {
	titler := cases.Title(language.AmericanEnglish)
	s = strings.ReplaceAll(s, "mu-app", dashName)
	s = strings.ReplaceAll(s, "mu_app", underName)
	s = strings.ReplaceAll(s, "muapp", packName)
	s = strings.ReplaceAll(s, "muApp", camelName)
	s = strings.ReplaceAll(s, "MuApp", titler.String(camelName))
	return s
}
