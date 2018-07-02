package index_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fossas/go-resolve/index"
	"github.com/fossas/go-resolve/models"
)

func TestIndexRepository(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	pkgs, err := index.Repository("github.com/stretchr/testify/assert")
	assert.NoError(t, err)
	t.Log(pkgs)
}

func TestFindRepositoryGit(t *testing.T) {
	vcs, path, err := index.FindRepository(".")
	assert.NoError(t, err)
	expectedPath, err := filepath.Abs("..")
	assert.NoError(t, err)
	actualPath, err := filepath.Abs(path)
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
	assert.Equal(t, models.Git, vcs)
}
