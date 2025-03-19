package shell

import (
	"io"
	"os"

	"golang.org/x/term"
)

func cost(x, y byte) int {
	if x == y {
		return 0
	}

	return 1
}

func editDistance(a, b string) int {
	l1 := len(a)
	l2 := len(b)

	dp := make([][]int, l1+1)
	for i := range dp {
		dp[i] = make([]int, l2+1)
	}

	for i := 0; i <= l1; i++ {
		dp[i][0] = i
	}

	for i := 0; i <= l2; i++ {
		dp[0][i] = i
	}

	for i := 1; i <= l1; i++ {
		for j := 1; j <= l2; j++ {
			dp[i][j] = min(
				dp[i-1][j]+1,                      // deletion
				dp[i][j-1]+1,                      // insertion
				dp[i-1][j-1]+cost(a[i-1], b[j-1]), // substitution
			)
		}
	}

	return dp[l1][l2]
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}
