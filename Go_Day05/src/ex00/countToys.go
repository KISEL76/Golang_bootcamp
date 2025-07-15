package toycounter

type TreeNode struct {
    HasToy bool
    Left *TreeNode
    Right *TreeNode
}

func dfs(node *TreeNode) int {
	if node == nil {
		return 0
	}

	toy := 0
	if node.HasToy {
		toy = 1
	}
	return dfs(node.Left) + dfs(node.Right) + toy
}

func areToysBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	leftCount := dfs(root.Left)
	rightCount := dfs(root.Right)

	return leftCount == rightCount
}