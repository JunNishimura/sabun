package diff

import "fmt"

type Diff struct {
	a, b []rune
	m, n int
	ed   int
}

func NewDiff(a, b []rune) *Diff {
	d := &Diff{}
	m, n := len(a), len(b)
	if m >= n {
		a, b = b, a
		n, m = m, n
	}
	d.a = a
	d.b = b
	d.n = n
	d.m = m
	d.ed = 0

	return d
}

func (d *Diff) EditDistance() int {
	return d.ed
}

func (d *Diff) Compose() {
	delta := d.n - d.m
	offset := d.m + 1
	fp := make([]int, d.n+d.m+3)
	for i := range fp {
		fp[i] = -1
	}
	for p := 0; ; p++ {
		fmt.Println(p)
		for k := -p; k <= delta-1; k++ {
			fp[k+offset] = d.snake(k, max(fp[k-1+offset]+1, fp[k+1+offset]))
		}
		for k := delta + p; k >= delta+1; k-- {
			fp[k+offset] = d.snake(k, max(fp[k-1+offset]+1, fp[k+1+offset]))
		}
		fp[delta+offset] = d.snake(delta, max(fp[delta-1+offset]+1, fp[delta+1+offset]))
		if fp[delta+offset] == d.n {
			d.ed = delta + 2*p
			break
		}
	}
}

func (d *Diff) snake(k, y int) int {
	x := y - k
	for x < d.m && y < d.n && d.a[x] == d.b[y] {
		x++
		y++
	}
	return y
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
