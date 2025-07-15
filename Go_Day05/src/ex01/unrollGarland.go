package garland

type TreeNode struct {
	Val   bool
	Left  *TreeNode
	Right *TreeNode
}

func unrollGarland(root *TreeNode) []bool {
	result := []bool{}
	if root == nil {
		return result
	}

	queue := []*TreeNode{root}
	level := 0

	for len(queue) > 0 {
		n := len(queue)
		current := make([]bool, n)

		for i := range n {
			node := queue[0]
			queue = queue[1:]

			idx := i
			if level%2 == 0 {
				idx = n - 1 - i
			}
			current[idx] = node.Val

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		result = append(result, current...)
		level++
	}

	return result
}
