// Parts of this code modified from aws-controllers-k8s/code-generator/pkg/util
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

package git

import (
	"context"

	gogit "gopkg.in/src-d/go-git.v4"
	gogitplumbing "gopkg.in/src-d/go-git.v4/plumbing"

	"errors"
	"fmt"
	"io"
)

type Repository = gogit.Repository

var Open = gogit.PlainOpen

// getRepositoryTagRef returns the git reference (commit hash) of a given tag.
// NOTE: It is not possible to checkout a tag without knowing it's reference.
//
// Calling this function is equivalent to executing `git rev-list -n 1 $tagName`
func getRepositoryTagRef(repo *Repository, tagName string) (*gogitplumbing.Reference, error) {
	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	for {
		tagRef, err := tagRefs.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error finding tag reference: %v", err)
		}
		if tagRef.Name().Short() == tagName {
			return tagRef, nil
		}
	}
	return nil, errors.New("tag reference not found")
}

// FetchTags fetches a repository's remote tags.
//
// Calling this function is equivalent to executing `git -C $path fetch --all
// --tags`
func FetchTags(
	ctx context.Context,
	repo *Repository,
) error {
	err := repo.FetchContext(ctx, &gogit.FetchOptions{
		Progress: nil,
		Tags:     gogit.AllTags,
	})
	// weirdly go-git returns a error "Already up to date" when all tags
	// are already fetched. We should ignore this error.
	if err == gogit.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

// CheckoutTag checkouts a repository tag by looking for the tag
// reference then calling the checkout function.
//
// Calling This function is equivalent to executing `git checkout tags/$tag`
func CheckoutTag(
	ctx context.Context,
	repo *Repository,
	tag string,
) error {
	tagRef, err := getRepositoryTagRef(repo, tag)
	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = wt.Checkout(&gogit.CheckoutOptions{
		// Checkout only take hashes or branch names.
		Hash: tagRef.Hash(),
	})
	return err
}

// Clone clones a git repository into a given directory and returns a
// Repository object that can be used to manipulate that clone'd repo.
//
// Calling his function is equivalent to executing `git clone $repositoryURL
// $path`
func Clone(
	ctx context.Context,
	path,
	repositoryURL string,
) (*Repository, error) {
	if _, err := gogit.PlainCloneContext(ctx, path, false, &gogit.CloneOptions{
		URL:      repositoryURL,
		Progress: nil,
		// Clone and fetch all tags
		Tags: gogit.AllTags,
	}); err != nil {
		return nil, err

	}
	return Open(path)
}
