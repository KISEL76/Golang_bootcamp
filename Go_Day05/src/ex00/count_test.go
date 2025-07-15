package toycounter

import (
	"testing"
)

func TestAreToysBalanced(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		expected bool
	}{
		/* 
		Example 1:
				  1
				/	\
			   1	 0
			  / \  	/ \
			 1	 0  1  1  
	    */
		{
		name: "Case 1: toys are balanced (2 vs 2)",
		root: &TreeNode {
			HasToy: true,
			Left: &TreeNode {
				HasToy: true,
				Left: &TreeNode {HasToy: true},
				Right: &TreeNode {HasToy: false},
			},
			Right: &TreeNode {
				HasToy: false,
				Left: &TreeNode {HasToy: true},
				Right: &TreeNode {HasToy: true},
				},
			},
		expected: true,
		},
		{
		/* 
		Example 2:
				  0
				/	\
			   1	 0
			  / \  	
			 1	 1  
	    */
		name: "Case 2: toys are disbalanced (3 from left vs 1 from right)",
		root: &TreeNode {
			HasToy: false,
			Left: &TreeNode {
				HasToy: true,
				Left: &TreeNode {HasToy: true},
				Right: &TreeNode {HasToy: true},
			},
			Right: &TreeNode {
				HasToy: true,
			},
		},
		expected: false,
		},
		{
		/* 
		Example 3:
				  0
				/	\
			 nil    nil
	    */
			name: "Case 3: only root exists",
			root: &TreeNode {HasToy: true},
			expected: true,
		},
		{
		/* 
		Example 4:
				  1
				/	\
			   0	 0
			    \  	/
			 	 0 0
	    */
			name: "Case 4: no toys",
			root: &TreeNode {
				HasToy: true,
				Left: &TreeNode {
					HasToy: false,
					Right: &TreeNode {HasToy: false},
				},
				Right: &TreeNode {
					HasToy: false,
					Left: &TreeNode {HasToy: false},
				},
			},
			expected: true,
		},
		{
		/* 
		Example 5:
				  1
				/	\
			   1	 0
		*/
			name: "Case 5: only 1 toy",
			root: &TreeNode {
				HasToy: true,
				Left: &TreeNode {HasToy: true},
				Right: &TreeNode {HasToy: false},
			},
			expected: false,
		},			
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := areToysBalanced(test.root)
			if got != test.expected {
				t.Errorf("Test '%s' wasn't completed right: expected %v, got %v", test.name, test.expected, got)
			}
		})
	}
}