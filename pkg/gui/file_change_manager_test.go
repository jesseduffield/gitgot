package gui

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	scenarios := []struct {
		name           string
		root           *models.FileChangeNode
		collapsedPaths map[string]bool
		expected       []string
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			root: &models.FileChangeNode{
				Path: "",
				Children: []*models.FileChangeNode{
					{File: &models.File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Path: "test"},
				},
			},
			expected: []string{" M test"},
		},
		{
			name: "big example",
			root: &models.FileChangeNode{
				Path: "",
				Children: []*models.FileChangeNode{
					{
						Path: "dir1",
						Children: []*models.FileChangeNode{
							{
								File: &models.File{Name: "dir1/file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir1/file2",
							},
							{
								File: &models.File{Name: "dir1/file3", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir1/file3",
							},
						},
					},
					{
						Path: "dir2",
						Children: []*models.FileChangeNode{
							{
								Path: "dir2/dir2",
								Children: []*models.FileChangeNode{
									{
										File: &models.File{Name: "dir2/dir2/file3", ShortStatus: " M", HasStagedChanges: true},
										Path: "dir2/dir2/file3",
									},
									{
										File: &models.File{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
										Path: "dir2/dir2/file4",
									},
								},
							},
							{
								File: &models.File{Name: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir2/file5",
							},
						},
					},
					{
						File: &models.File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "file1",
					},
				},
			},
			expected:       []string{"dir1 ►", "dir2 ▼", "├─ dir2 ▼", "│  ├─  M file3", "│  └─ M  file4", "└─ M  file5", "M  file1"},
			collapsedPaths: map[string]bool{"dir1": true},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			mngr := &FileChangeManager{Tree: s.root, CollapsedPaths: s.collapsedPaths}
			result := mngr.Render("", nil)
			assert.EqualValues(t, s.expected, result)
		})
	}
}
