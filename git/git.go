package git

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitRepo struct {
	r *git.Repository
	h plumbing.Hash
}

type TagList []*object.Tag

func (tl TagList) Len() int {
	return len(tl)
}

func (tl TagList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

// sorting tags in reverse chronological order
func (tl TagList) Less(i, j int) bool {
	return tl[i].Tagger.When.After(tl[j].Tagger.When)
}

func Open(path string, ref string) (*GitRepo, error) {
	var err error
	g := GitRepo{}
	g.r, err = git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}

	if ref == "" {
		head, err := g.r.Head()
		if err != nil {
			return nil, fmt.Errorf("getting head of %s: %w", path, err)
		}
		g.h = head.Hash()
	} else {
		hash, err := g.r.ResolveRevision(plumbing.Revision(ref))
		if err != nil {
			return nil, fmt.Errorf("resolving rev %s for %s: %w", ref, path, err)
		}
		g.h = *hash
	}
	return &g, nil
}

func (g *GitRepo) Commits() ([]*object.Commit, error) {
	ci, err := g.r.Log(&git.LogOptions{From: g.h})
	if err != nil {
		return nil, fmt.Errorf("commits from ref: %w", err)
	}

	commits := []*object.Commit{}
	err = ci.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("commits from ref: %w", err)
	}

	return commits, nil
}

func (g *GitRepo) LastCommit() (*object.Commit, error) {
	c, err := g.r.CommitObject(g.h)
	if err != nil {
		return nil, fmt.Errorf("last commit: %w", err)
	}
	return c, nil
}

func (g *GitRepo) FileContent(path string) (string, error) {
	c, err := g.r.CommitObject(g.h)
	if err != nil {
		return "", fmt.Errorf("commit object: %w", err)
	}

	tree, err := c.Tree()
	if err != nil {
		return "", fmt.Errorf("file tree: %w", err)
	}

	file, err := tree.File(path)
	if err != nil {
		return "", err
	}

	isbin, _ := file.IsBinary()

	if !isbin {
		return file.Contents()
	} else {
		return "Not displaying binary file", nil
	}
}

func (g *GitRepo) Tags() ([]*object.Tag, error) {
	ti, err := g.r.TagObjects()
	if err != nil {
		return nil, fmt.Errorf("tag objects: %w", err)
	}

	tags := TagList{}

	_ = ti.ForEach(func(t *object.Tag) error {
		for i, existing := range tags {
			if existing.Name == t.Name {
				if t.Tagger.When.After(existing.Tagger.When) {
					tags[i] = t
				}
				return nil
			}
		}
		tags = append(tags, t)
		return nil
	})

	sort.Sort(tags)

	return tags, nil
}

func (g *GitRepo) Branches() ([]*plumbing.Reference, error) {
	bi, err := g.r.Branches()
	if err != nil {
		return nil, fmt.Errorf("branchs: %w", err)
	}

	branches := []*plumbing.Reference{}

	_ = bi.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref)
		return nil
	})

	return branches, nil
}

func (g *GitRepo) FindMainBranch(branches []string) (string, error) {
	for _, b := range branches {
		_, err := g.r.ResolveRevision(plumbing.Revision(b))
		if err == nil {
			return b, nil
		}
	}
	return "", fmt.Errorf("unable to find main branch")
}
