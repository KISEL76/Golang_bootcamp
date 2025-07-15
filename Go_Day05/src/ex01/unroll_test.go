package garland

import "testing"

func equalSlices(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestUnrollGarland(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []bool
	}{
		{
			/*
			      1
			    /   \
			   1     0
			  / \   / \
			 1   0 1   1
			*/
			name: "Case 1: balanced",
			root: &TreeNode{
				Val: true,
				Left: &TreeNode{
					Val:   true,
					Left:  &TreeNode{Val: true},
					Right: &TreeNode{Val: false},
				},
				Right: &TreeNode{
					Val:   false,
					Left:  &TreeNode{Val: true},
					Right: &TreeNode{Val: true},
				},
			},
			expected: []bool{true, true, false, true, true, false, true},
		},
		{
			/*
			      0
			    /   \
			   1     0
			  / \
			 1   1
			*/
			name: "Case 2: unbalanced left",
			root: &TreeNode{
				Val: false,
				Left: &TreeNode{
					Val:   true,
					Left:  &TreeNode{Val: true},
					Right: &TreeNode{Val: true},
				},
				Right: &TreeNode{
					Val: false,
				},
			},
			expected: []bool{false, true, false, true, true},
		},
		{
			name:     "Case 3: root only",
			root:     &TreeNode{Val: true},
			expected: []bool{true},
		},
		{
			/*
			      1
			    /   \
			   0     0
			    \   /
			     0 0
			*/
			name: "Case 4: all false leaves",
			root: &TreeNode{
				Val: true,
				Left: &TreeNode{
					Val:   false,
					Right: &TreeNode{Val: false},
				},
				Right: &TreeNode{
					Val:  false,
					Left: &TreeNode{Val: false},
				},
			},
			expected: []bool{true, false, false, false, false},
		},
		{
			/*
			      1
			    /   \
			   1     0
			*/
			name: "Case 5: two children",
			root: &TreeNode{
				Val:   true,
				Left:  &TreeNode{Val: true},
				Right: &TreeNode{Val: false},
			},
			expected: []bool{true, true, false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unrollGarland(test.root)
			if !equalSlices(got, test.expected) {
				t.Errorf("FAIL %s:\nExpected: %v\nGot:      %v", test.name, test.expected, got)
			}
		})
	}
}
