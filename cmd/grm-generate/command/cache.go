// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package command

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/anydotcloud/grm-generate/pkg/git"
)

const (
	defaultGitCloneTimeout = 180 * time.Second
	defaultGitFetchTimeout = 30 * time.Second
)

// ensureDir makes sure that a supplied directory exists and
// returns whether the directory already existed.
func ensureDir(fp string) (bool, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return false, os.MkdirAll(fp, os.ModePerm)
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("expected %s to be a directory", fp)
	}
	if !isDirWriteable(fp) {
		return true, fmt.Errorf("%s is not a writeable directory", fp)
	}

	return true, nil
}

// isDirWriteable returns true if the supplied directory path is writeable,
// false otherwise
func isDirWriteable(fp string) bool {
	testPath := filepath.Join(fp, "test")
	f, err := os.Create(testPath)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(testPath)
	return true
}

// cacheRepo ensures that we have a git clone'd copy of the supplied source code
// repository
func cacheRepo(
	ctx context.Context,
	cachePath string,
	repoURL string,
	tag string, // optional Git tag to checkout
) error {
	var err error
	var repo *git.Repository

	if _, err = ensureDir(cachePath); err != nil {
		return err
	}

	repoName := path.Base(repoURL)

	// Clone repository if it doesn't exist
	repoPath := filepath.Join(cachePath, repoName)
	if _, err = os.Stat(repoPath); os.IsNotExist(err) {
		ctx, cancel := context.WithTimeout(ctx, defaultGitCloneTimeout)
		defer cancel()
		repo, err = git.Clone(ctx, repoPath, repoURL)
		if err != nil {
			return fmt.Errorf("cannot clone repository: %v", err)
		}
	} else {
		if repo, err = git.Open(repoPath); err != nil {
			return fmt.Errorf("could not open repository: %v", err)
		}
	}

	fctx, cancel := context.WithTimeout(ctx, defaultGitFetchTimeout)
	defer cancel()
	if err = git.FetchTags(fctx, repo); err != nil {
		return fmt.Errorf("cannot fetch tags: %v", err)
	}

	if tag != "" {
		if err = git.CheckoutTag(ctx, repo, tag); err != nil {
			return fmt.Errorf("cannot checkout tag: %v", err)
		}
	}

	return nil
}
