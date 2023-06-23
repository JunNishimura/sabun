package diff

import "fmt"

const (
	SesInsert SesType = iota
	SesDelete
	SesCommon
)

type SesType int

type SesElem struct {
	e rune
	t SesType
}

func (e *SesElem) GetElem() rune {
	return e.e
}

func (e *SesElem) GetType() SesType {
	return e.t
}

type Point struct {
	x, y int
}

type PointWithRoute struct {
	x, y, r int // rはDiff.routeスライスの何番目のインデックスから辿り着いたかを格納している
}

type Diff struct {
	a, b    []rune
	m, n    int
	reverse bool
	ed      int
	lcs     []rune
	ses     []*SesElem
	path    []int // 各kのroute idを保持
	route   []*PointWithRoute
}

func NewDiff(a, b []rune) *Diff {
	d := &Diff{}
	m, n := len(a), len(b)
	reverse := false
	if m >= n {
		a, b = b, a
		n, m = m, n
		reverse = true
	}
	d.a = a
	d.b = b
	d.n = n
	d.m = m
	d.reverse = reverse
	d.ed = 0

	return d
}

func (d *Diff) EditDistance() int {
	return d.ed
}

func (d *Diff) Lcs() []rune {
	return d.lcs
}

func (d *Diff) Ses() []*SesElem {
	return d.ses
}

func (d *Diff) PrintSes() {
	for _, el := range d.ses {
		switch el.t {
		case SesInsert:
			fmt.Printf("+%v\n", el.e)
		case SesDelete:
			fmt.Printf("-%v\n", el.e)
		case SesCommon:
			fmt.Printf(" %v\n", el.e)
		}
	}
}

func (d *Diff) Compose() {
	delta := d.n - d.m
	offset := d.m + 1
	fp := make([]int, d.n+d.m+3)
	d.path = make([]int, d.n+d.m+3)
	d.route = make([]*PointWithRoute, 0)
	for i := range fp {
		fp[i] = -1
		d.path[i] = -1
	}
	for p := 0; ; p++ {
		for k := -p; k <= delta-1; k++ {
			fp[k+offset] = d.snake(k, fp[k-1+offset]+1, fp[k+1+offset], offset)
		}
		for k := delta + p; k >= delta+1; k-- {
			fp[k+offset] = d.snake(k, fp[k-1+offset]+1, fp[k+1+offset], offset)
		}
		fp[delta+offset] = d.snake(delta, fp[delta-1+offset]+1, fp[delta+1+offset], offset)
		if fp[delta+offset] == d.n {
			d.ed = delta + 2*p
			break
		}
	}

	r := d.path[delta+offset]
	points := make([]*Point, 0)
	for r != -1 {
		points = append(points, &Point{x: d.route[r].x, y: d.route[r].y})
		r = d.route[r].r
	}

	d.recordSeq(points)
}

func (d *Diff) recordSeq(points []*Point) {
	px, py := 0, 0 // posx, posy
	for i := len(points) - 1; i >= 0; i-- {
		for px < points[i].x || py < points[i].y {
			if points[i].x-px < points[i].y-py {
				if d.reverse {
					d.ses = append(d.ses, &SesElem{e: d.b[py], t: SesDelete})
				} else {
					d.ses = append(d.ses, &SesElem{e: d.b[py], t: SesInsert})
				}
				py++
			} else if points[i].x-px > points[i].y-py {
				if d.reverse {
					d.ses = append(d.ses, &SesElem{e: d.a[px], t: SesInsert})
				} else {
					d.ses = append(d.ses, &SesElem{e: d.a[px], t: SesDelete})
				}
				px++
			} else {
				if d.reverse {
					d.lcs = append(d.lcs, d.b[py])
					d.ses = append(d.ses, &SesElem{e: d.b[py], t: SesCommon})
				} else {
					d.lcs = append(d.lcs, d.a[px])
					d.ses = append(d.ses, &SesElem{e: d.a[px], t: SesCommon})
				}
				px++
				py++
			}
		}
	}
}

func (d *Diff) snake(k, pi, pd, offset int) int {
	r := 0
	if pi > pd {
		r = d.path[k-1+offset]
	} else {
		r = d.path[k+1+offset]
	}

	y := max(pi, pd)
	x := y - k
	for x < d.m && y < d.n && d.a[x] == d.b[y] {
		x++
		y++
	}

	d.path[k+offset] = len(d.route)
	d.route = append(d.route, &PointWithRoute{x: x, y: y, r: r})

	return y
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
