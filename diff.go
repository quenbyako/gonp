// The algorithm implemented here is based on "An O(NP) Sequence Comparison Algorithm"
// by described by Sun Wu, Udi Manber and Gene Myers

package gonp

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
)

// SesType is manipulaton type
type SesType int

const (
	// SesDelete is manipulaton type of deleting element in SES
	SesDelete SesType = iota
	// SesCommon is manipulaton type of same element in SES
	SesCommon
	// SesAdd is manipulaton type of adding element in SES
	SesAdd
)

const (
	// limit of cordinate size
	DefaultRouteSize = 2000000
)

// Point is coordinate in edit graph
type Point struct{ x, y int }

// PointWithRoute is coordinate in edit graph attached route
type PointWithRoute struct{ x, y, r int }

// Type constraints for element in SES
type Elem = any

// SesElem is element of SES
type SesElem[T any] struct {
	elem T
	typ  SesType
	aIdx int
	bIdx int
}

func (e SesElem[T]) Cmp(b SesElem[T], c func(T, T) int) int {
	for _, v := range []int{
		c(e.elem, b.elem),
		cmp.Compare(e.typ, b.typ),
		cmp.Compare(e.aIdx, b.aIdx),
		cmp.Compare(e.bIdx, b.bIdx),
	} {
		if v != 0 {
			return v
		}
	}

	return 0
}

// GetElem is getter of element of SES
func (e *SesElem[T]) GetElem() T { return e.elem }

// GetType is getter of manipulation type of SES
func (e *SesElem[T]) GetType() SesType { return e.typ }

// Diff is context for calculating difference between a and b
type Diff[T Elem] struct {
	a, b           []T
	aLen, bLen     int
	ox, oy         int
	ed             int
	// lsc means Longest Common Subsequence
	lcs            []T
	ses            []SesElem[T]
	reverse        bool
	path           []int
	onlyEd         bool
	pointWithRoute []PointWithRoute
	contextSize    int
	routeSize      int
	cmp            func(T, T) int
}

func New[T cmp.Ordered](a, b []T) *Diff[T] { return NewCmp(a, b, cmp.Compare) }

// NewCmp is initializer of Diff
func NewCmp[T any](a, b []T, cmp func(T, T) int) *Diff[T] {
	var reverse bool
	if len(a) >= len(b) {
		a, b = b, a
		reverse = true
	}

	return &Diff[T]{
		a:           a,
		b:           b,
		aLen:        len(a),
		bLen:        len(b),
		ed:          0,
		reverse:     reverse,
		onlyEd:      false,
		contextSize: DefaultContextSize,
		routeSize:   DefaultRouteSize,
		cmp:         cmp,
	}
}

// OnlyEd enables to calculate only edit distance
func (d *Diff[T]) OnlyEd() *Diff[T] { d.onlyEd = true; return d }

// SetContextSize sets the context size of unified format difference
func (d *Diff[T]) SetContextSize(n int) *Diff[T] { d.contextSize = n; return d }

// SetRouteSize sets the context size of unified format difference
func (d *Diff[T]) SetRouteSize(n int) *Diff[T] { d.routeSize = n; return d }

// EditDistance returns edit distance between a and b
func (d *Diff[T]) EditDistance() int { return d.ed }

// Lcs returns LCS (Longest Common Subsequence) between a and b
func (diff *Diff[T]) Lcs() []T { return diff.lcs }

// Ses return SES (Shortest Edit Script) between a and b
func (diff *Diff[T]) Ses() []SesElem[T] {
	return diff.ses
}

// PrintSes prints shortest edit script between a and b
func (diff *Diff[T]) PrintSes() {
	fmt.Print(diff.SprintSes())
}

// SprintSes returns string about shortest edit script between a and b
func (diff *Diff[T]) SprintSes() string {
	var buf bytes.Buffer
	diff.FprintSes(&buf)
	return buf.String()
}

// FprintSes emit about shortest edit script between a and b to w
func (diff *Diff[T]) FprintSes(w io.Writer) {
	for _, e := range diff.ses {
		switch e.typ {
		case SesDelete:
			fmt.Fprintf(w, "-%v\n", e.elem)
		case SesAdd:
			fmt.Fprintf(w, "+%v\n", e.elem)
		case SesCommon:
			fmt.Fprintf(w, " %v\n", e.elem)
		}
	}
}

// Compose composes diff between a and b
func (diff *Diff[T]) Compose() {
ONP:
	fp := make([]int, diff.aLen+diff.bLen+3)
	diff.path = make([]int, diff.aLen+diff.bLen+3)
	diff.pointWithRoute = make([]PointWithRoute, 0)

	for i := range fp {
		fp[i] = -1
		diff.path[i] = -1
	}

	offset := diff.aLen + 1
	delta := diff.bLen - diff.aLen
	for p := 0; ; p++ {
		for k := -p; k <= delta-1; k++ {
			fp[k+offset] = diff.snake(k, fp[k-1+offset]+1, fp[k+1+offset], offset)
		}

		for k := delta + p; k >= delta+1; k-- {
			fp[k+offset] = diff.snake(k, fp[k-1+offset]+1, fp[k+1+offset], offset)
		}

		fp[delta+offset] = diff.snake(delta, fp[delta-1+offset]+1, fp[delta+1+offset], offset)

		if fp[delta+offset] >= diff.bLen || len(diff.pointWithRoute) > diff.routeSize {
			diff.ed += delta + 2*p
			break
		}
	}

	if diff.onlyEd {
		return
	}

	r := diff.path[delta+offset]
	epc := make([]Point, 0)
	for r != -1 {
		epc = append(epc, Point{x: diff.pointWithRoute[r].x, y: diff.pointWithRoute[r].y})
		r = diff.pointWithRoute[r].r
	}

	if !diff.recordSeq(epc) {
		goto ONP
	}
}

func (diff *Diff[T]) snake(k, p, pp, offset int) int {
	r := 0
	if p > pp {
		r = diff.path[k-1+offset]
	} else {
		r = diff.path[k+1+offset]
	}

	y := max(p, pp)
	x := y - k

	for x < diff.aLen && y < diff.bLen && diff.cmp(diff.a[x], diff.b[y]) == 0 {
		x++
		y++
	}

	if !diff.onlyEd {
		diff.path[k+offset] = len(diff.pointWithRoute)
		diff.pointWithRoute = append(diff.pointWithRoute, PointWithRoute{x: x, y: y, r: r})
	}

	return y
}

func (diff *Diff[T]) recordSeq(epc []Point) bool {
	x, y := 1, 1
	px, py := 0, 0
	for i := len(epc) - 1; i >= 0; i-- {
		for (px < epc[i].x) || (py < epc[i].y) {
			if (epc[i].y - epc[i].x) > (py - px) {
				if diff.reverse {
					diff.ses = append(diff.ses, SesElem[T]{elem: diff.b[py], typ: SesDelete, aIdx: y + diff.oy, bIdx: 0})
				} else {
					diff.ses = append(diff.ses, SesElem[T]{elem: diff.b[py], typ: SesAdd, aIdx: 0, bIdx: y + diff.oy})
				}
				y++
				py++
			} else if epc[i].y-epc[i].x < py-px {
				if diff.reverse {
					diff.ses = append(diff.ses, SesElem[T]{elem: diff.a[px], typ: SesAdd, aIdx: 0, bIdx: x + diff.ox})
				} else {
					diff.ses = append(diff.ses, SesElem[T]{elem: diff.a[px], typ: SesDelete, aIdx: x + diff.ox, bIdx: 0})

				}
				x++
				px++
			} else {
				if diff.reverse {
					diff.lcs = append(diff.lcs, diff.b[py])
					diff.ses = append(diff.ses, SesElem[T]{elem: diff.b[py], typ: SesCommon, aIdx: y + diff.oy, bIdx: x + diff.ox})
				} else {
					diff.lcs = append(diff.lcs, diff.a[px])
					diff.ses = append(diff.ses, SesElem[T]{elem: diff.a[px], typ: SesCommon, aIdx: x + diff.ox, bIdx: y + diff.oy})
				}
				x++
				y++
				px++
				py++
			}
		}
	}

	if x > diff.aLen && y > diff.bLen {
		// all recording succeeded
	} else {
		diff.a = diff.a[x-1:]
		diff.b = diff.b[y-1:]
		diff.aLen = len(diff.a)
		diff.bLen = len(diff.b)
		diff.ox = x - 1
		diff.oy = y - 1
		return false
	}

	return true
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
