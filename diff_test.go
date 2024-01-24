package gonp

import (
	"cmp"
	"slices"
	"testing"
)

func equalsSesElemSlice[T Elem](ses1, ses2 []SesElem[T], cmp func(SesElem[T], SesElem[T]) int) bool {
	m, n := len(ses1), len(ses2)
	if m != n {
		return false
	}
	for i := 0; i < m; i++ {
		if cmp(ses1[i], ses2[i]) != 0 {
			return false
		}
	}
	return true
}

func equalsUniHunks[T Elem](uniHunks1, uniHunks2 []UniHunk[T], cmp func(T, T) int) bool {
	m, n := len(uniHunks1), len(uniHunks2)
	if m != n {
		return false
	}
	for i := 0; i < m; i++ {
		if uniHunks1[i].a != uniHunks2[i].a {
			return false
		}
		if uniHunks1[i].b != uniHunks2[i].b {
			return false
		}
		if uniHunks1[i].c != uniHunks2[i].c {
			return false
		}
		if uniHunks1[i].d != uniHunks2[i].d {
			return false
		}
		if !equalsSesElemSlice(uniHunks1[i].changes, uniHunks2[i].changes, func(se1, se2 SesElem[T]) int { return se1.Cmp(se2, cmp) }) {
			return false
		}
	}
	return true
}

func TestStringDiff(t *testing.T) {

	tests := []struct {
		name     string
		a        string
		b        string
		ed       int
		lcs      string
		ses      []SesElem[rune]
		uniHunks []UniHunk[rune]
	}{
		{
			name: "string diff1",
			a:    "abc",
			b:    "abd",
			ed:   2,
			lcs:  "ab",
			ses: []SesElem[rune]{
				{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
				{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 2},
				{elem: 'c', typ: SesDelete, aIdx: 3, bIdx: 0},
				{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 3},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 3, c: 1, d: 3,
					changes: []SesElem[rune]{
						{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
						{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 2},
						{elem: 'c', typ: SesDelete, aIdx: 3, bIdx: 0},
						{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 3},
					},
				},
			},
		},
		{
			name: "string diff2",
			a:    "abcdef",
			b:    "dacfea",
			ed:   6,
			lcs:  "acf",
			ses: []SesElem[rune]{
				{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 1},
				{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 2},
				{elem: 'b', typ: SesDelete, aIdx: 2, bIdx: 0},
				{elem: 'c', typ: SesCommon, aIdx: 3, bIdx: 3},
				{elem: 'd', typ: SesDelete, aIdx: 4, bIdx: 0},
				{elem: 'e', typ: SesDelete, aIdx: 5, bIdx: 0},
				{elem: 'f', typ: SesCommon, aIdx: 6, bIdx: 4},
				{elem: 'e', typ: SesAdd, aIdx: 0, bIdx: 5},
				{elem: 'a', typ: SesAdd, aIdx: 0, bIdx: 6},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 6, c: 1, d: 6,
					changes: []SesElem[rune]{
						{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 1},
						{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 2},
						{elem: 'b', typ: SesDelete, aIdx: 2, bIdx: 0},
						{elem: 'c', typ: SesCommon, aIdx: 3, bIdx: 3},
						{elem: 'd', typ: SesDelete, aIdx: 4, bIdx: 0},
						{elem: 'e', typ: SesDelete, aIdx: 5, bIdx: 0},
						{elem: 'f', typ: SesCommon, aIdx: 6, bIdx: 4},
						{elem: 'e', typ: SesAdd, aIdx: 0, bIdx: 5},
						{elem: 'a', typ: SesAdd, aIdx: 0, bIdx: 6},
					},
				},
			},
		},
		{
			name: "string diff3",
			a:    "acbdeacbed",
			b:    "acebdabbabed",
			ed:   6,
			lcs:  "acbdabed",
			ses: []SesElem[rune]{
				{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
				{elem: 'c', typ: SesCommon, aIdx: 2, bIdx: 2},
				{elem: 'e', typ: SesAdd, aIdx: 0, bIdx: 3},
				{elem: 'b', typ: SesCommon, aIdx: 3, bIdx: 4},
				{elem: 'd', typ: SesCommon, aIdx: 4, bIdx: 5},
				{elem: 'e', typ: SesDelete, aIdx: 5, bIdx: 0},
				{elem: 'a', typ: SesCommon, aIdx: 6, bIdx: 6},
				{elem: 'c', typ: SesDelete, aIdx: 7, bIdx: 0},
				{elem: 'b', typ: SesCommon, aIdx: 8, bIdx: 7},
				{elem: 'b', typ: SesAdd, aIdx: 0, bIdx: 8},
				{elem: 'a', typ: SesAdd, aIdx: 0, bIdx: 9},
				{elem: 'b', typ: SesAdd, aIdx: 0, bIdx: 10},
				{elem: 'e', typ: SesCommon, aIdx: 9, bIdx: 11},
				{elem: 'd', typ: SesCommon, aIdx: 10, bIdx: 12},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 10, c: 1, d: 12,
					changes: []SesElem[rune]{
						{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
						{elem: 'c', typ: SesCommon, aIdx: 2, bIdx: 2},
						{elem: 'e', typ: SesAdd, aIdx: 0, bIdx: 3},
						{elem: 'b', typ: SesCommon, aIdx: 3, bIdx: 4},
						{elem: 'd', typ: SesCommon, aIdx: 4, bIdx: 5},
						{elem: 'e', typ: SesDelete, aIdx: 5, bIdx: 0},
						{elem: 'a', typ: SesCommon, aIdx: 6, bIdx: 6},
						{elem: 'c', typ: SesDelete, aIdx: 7, bIdx: 0},
						{elem: 'b', typ: SesCommon, aIdx: 8, bIdx: 7},
						{elem: 'b', typ: SesAdd, aIdx: 0, bIdx: 8},
						{elem: 'a', typ: SesAdd, aIdx: 0, bIdx: 9},
						{elem: 'b', typ: SesAdd, aIdx: 0, bIdx: 10},
						{elem: 'e', typ: SesCommon, aIdx: 9, bIdx: 11},
						{elem: 'd', typ: SesCommon, aIdx: 10, bIdx: 12},
					},
				},
			},
		},
		{
			name: "string diff4",
			a:    "abcbda",
			b:    "bdcaba",
			ed:   4,
			lcs:  "bcba",
			ses: []SesElem[rune]{
				{elem: 'a', typ: SesDelete, aIdx: 1, bIdx: 0},
				{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 1},
				{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 2},
				{elem: 'c', typ: SesCommon, aIdx: 3, bIdx: 3},
				{elem: 'a', typ: SesAdd, aIdx: 0, bIdx: 4},
				{elem: 'b', typ: SesCommon, aIdx: 4, bIdx: 5},
				{elem: 'd', typ: SesDelete, aIdx: 5, bIdx: 0},
				{elem: 'a', typ: SesCommon, aIdx: 6, bIdx: 6},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 6, c: 1, d: 6,
					changes: []SesElem[rune]{
						{elem: 'a', typ: SesDelete, aIdx: 1, bIdx: 0},
						{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 1},
						{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 2},
						{elem: 'c', typ: SesCommon, aIdx: 3, bIdx: 3},
						{elem: 'a', typ: SesAdd, aIdx: 0, bIdx: 4},
						{elem: 'b', typ: SesCommon, aIdx: 4, bIdx: 5},
						{elem: 'd', typ: SesDelete, aIdx: 5, bIdx: 0},
						{elem: 'a', typ: SesCommon, aIdx: 6, bIdx: 6},
					},
				},
			},
		},
		{
			name: "string diff5",
			a:    "bokko",
			b:    "bokkko",
			ed:   1,
			lcs:  "bokko",
			ses: []SesElem[rune]{
				{elem: 'b', typ: SesCommon, aIdx: 1, bIdx: 1},
				{elem: 'o', typ: SesCommon, aIdx: 2, bIdx: 2},
				{elem: 'k', typ: SesCommon, aIdx: 3, bIdx: 3},
				{elem: 'k', typ: SesCommon, aIdx: 4, bIdx: 4},
				{elem: 'k', typ: SesAdd, aIdx: 0, bIdx: 5},
				{elem: 'o', typ: SesCommon, aIdx: 5, bIdx: 6},
			},
			uniHunks: []UniHunk[rune]{
				{a: 2, b: 4, c: 2, d: 5,
					changes: []SesElem[rune]{
						{elem: 'o', typ: SesCommon, aIdx: 2, bIdx: 2},
						{elem: 'k', typ: SesCommon, aIdx: 3, bIdx: 3},
						{elem: 'k', typ: SesCommon, aIdx: 4, bIdx: 4},
						{elem: 'k', typ: SesAdd, aIdx: 0, bIdx: 5},
						{elem: 'o', typ: SesCommon, aIdx: 5, bIdx: 6},
					},
				},
			},
		},
		{
			name: "string diff6",
			a:    "abcaaaaaabd",
			b:    "abdaaaaaabc",
			ed:   4,
			lcs:  "abaaaaaab",
			ses: []SesElem[rune]{
				{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
				{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 2},
				{elem: 'c', typ: SesDelete, aIdx: 3, bIdx: 0},
				{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 3},
				{elem: 'a', typ: SesCommon, aIdx: 4, bIdx: 4},
				{elem: 'a', typ: SesCommon, aIdx: 5, bIdx: 5},
				{elem: 'a', typ: SesCommon, aIdx: 6, bIdx: 6},
				{elem: 'a', typ: SesCommon, aIdx: 7, bIdx: 7},
				{elem: 'a', typ: SesCommon, aIdx: 8, bIdx: 8},
				{elem: 'a', typ: SesCommon, aIdx: 9, bIdx: 9},
				{elem: 'b', typ: SesCommon, aIdx: 10, bIdx: 10},
				{elem: 'd', typ: SesDelete, aIdx: 11, bIdx: 0},
				{elem: 'c', typ: SesAdd, aIdx: 0, bIdx: 11},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 6, c: 1, d: 6,
					changes: []SesElem[rune]{
						{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
						{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 2},
						{elem: 'c', typ: SesDelete, aIdx: 3, bIdx: 0},
						{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 3},
						{elem: 'a', typ: SesCommon, aIdx: 4, bIdx: 4},
						{elem: 'a', typ: SesCommon, aIdx: 5, bIdx: 5},
						{elem: 'a', typ: SesCommon, aIdx: 6, bIdx: 6},
					},
				},
				{a: 8, b: 4, c: 8, d: 4,
					changes: []SesElem[rune]{
						{elem: 'a', typ: SesCommon, aIdx: 8, bIdx: 8},
						{elem: 'a', typ: SesCommon, aIdx: 9, bIdx: 9},
						{elem: 'b', typ: SesCommon, aIdx: 10, bIdx: 10},
						{elem: 'd', typ: SesDelete, aIdx: 11, bIdx: 0},
						{elem: 'c', typ: SesAdd, aIdx: 0, bIdx: 11},
					},
				},
			},
		},
		{
			name:     "empty string diff1",
			a:        "",
			b:        "",
			ed:       0,
			lcs:      "",
			ses:      []SesElem[rune]{},
			uniHunks: []UniHunk[rune]{},
		},
		{
			name: "empty string diff2",
			a:    "a",
			b:    "",
			ed:   1,
			lcs:  "",
			ses: []SesElem[rune]{
				{elem: 'a', typ: SesDelete, aIdx: 1, bIdx: 0},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 1, c: 0, d: 0, changes: []SesElem[rune]{
					{elem: 'a', typ: SesDelete, aIdx: 1, bIdx: 0},
				},
				},
			},
		},
		{
			name: "empty string diff3",
			a:    "",
			b:    "b",
			ed:   1,
			lcs:  "",
			ses: []SesElem[rune]{
				{elem: 'b', typ: SesAdd, aIdx: 0, bIdx: 1},
			},
			uniHunks: []UniHunk[rune]{
				{a: 0, b: 0, c: 1, d: 1, changes: []SesElem[rune]{
					{elem: 'b', typ: SesAdd, aIdx: 0, bIdx: 1},
				},
				},
			},
		},
		{
			name: "multi byte string diff",
			a:    "久保竜彦",
			b:    "久保達彦",
			ed:   2,
			lcs:  "久保彦",
			ses: []SesElem[rune]{
				{elem: '久', typ: SesCommon, aIdx: 1, bIdx: 1},
				{elem: '保', typ: SesCommon, aIdx: 2, bIdx: 2},
				{elem: '竜', typ: SesDelete, aIdx: 3, bIdx: 0},
				{elem: '達', typ: SesAdd, aIdx: 0, bIdx: 3},
				{elem: '彦', typ: SesCommon, aIdx: 4, bIdx: 4},
			},
			uniHunks: []UniHunk[rune]{
				{a: 1, b: 4, c: 1, d: 4, changes: []SesElem[rune]{
					{elem: '久', typ: SesCommon, aIdx: 1, bIdx: 1},
					{elem: '保', typ: SesCommon, aIdx: 2, bIdx: 2},
					{elem: '竜', typ: SesDelete, aIdx: 3, bIdx: 0},
					{elem: '達', typ: SesAdd, aIdx: 0, bIdx: 3},
					{elem: '彦', typ: SesCommon, aIdx: 4, bIdx: 4},
				},
				},
			},
		},
	}

	for _, tt := range tests {
		diff := New([]rune(tt.a), []rune(tt.b))
		diff.Compose()
		ed := diff.EditDistance()
		lcs := string(diff.Lcs())
		ses := diff.Ses()
		uniHunks := diff.UnifiedHunks()
		if tt.ed != ed {
			t.Fatalf(":%s:ed: want: %d, got: %d", tt.name, tt.ed, ed)
		}
		if tt.lcs != lcs {
			t.Fatalf(":%s:lcs: want: %s, got: %s", tt.name, tt.lcs, lcs)
		}
		if !equalsSesElemSlice(tt.ses, ses, func(se1, se2 SesElem[rune]) int { return se1.Cmp(se2, cmp.Compare) }) {
			t.Fatalf(":%s:ses: want: %v, got: %v", tt.name, tt.ses, ses)
		}

		if !equalsUniHunks(tt.uniHunks, uniHunks, cmp.Compare) {
			t.Fatalf(":%s:uniHunks: want: %v, got: %v", tt.name, tt.uniHunks, uniHunks)
		}
	}
}

func TestSliceDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		ed       int
		lcs      []int
		ses      []SesElem[int]
		uniHunks []UniHunk[int]
	}{
		{
			name: "int slice diff",
			a:    []int{1, 2, 3},
			b:    []int{1, 5, 3},
			ed:   2,
			lcs:  []int{1, 3},
			ses: []SesElem[int]{
				{elem: 1, typ: SesCommon, aIdx: 1, bIdx: 1},
				{elem: 2, typ: SesDelete, aIdx: 2, bIdx: 0},
				{elem: 5, typ: SesAdd, aIdx: 0, bIdx: 2},
				{elem: 3, typ: SesCommon, aIdx: 3, bIdx: 3},
			},
			uniHunks: []UniHunk[int]{
				{a: 1, b: 3, c: 1, d: 3,
					changes: []SesElem[int]{
						{elem: 1, typ: SesCommon, aIdx: 1, bIdx: 1},
						{elem: 2, typ: SesDelete, aIdx: 2, bIdx: 0},
						{elem: 5, typ: SesAdd, aIdx: 0, bIdx: 2},
						{elem: 3, typ: SesCommon, aIdx: 3, bIdx: 3},
					},
				},
			},
		},
		{
			name:     "empty slice diff",
			a:        []int{},
			b:        []int{},
			ed:       0,
			lcs:      []int{},
			ses:      []SesElem[int]{},
			uniHunks: []UniHunk[int]{},
		},
	}

	for _, tt := range tests {
		diff := New(tt.a, tt.b)
		diff.Compose()
		ed := diff.EditDistance()
		lcs := diff.Lcs()
		ses := diff.Ses()
		uniHunks := diff.UnifiedHunks()
		if tt.ed != ed {
			t.Fatalf(":%s:ed: want: %d, got: %d", tt.name, tt.ed, ed)
		}
		if !slices.Equal(tt.lcs, lcs) {
			t.Fatalf(":%s:lcs: want: %v, got: %v", tt.name, tt.lcs, lcs)
		}
		if !equalsSesElemSlice(tt.ses, ses, func(se1, se2 SesElem[int]) int { return se1.Cmp(se2, cmp.Compare) }) {
			t.Fatalf(":%s:ses: want: %v, got: %v", tt.name, tt.ses, ses)
		}

		if !equalsUniHunks(tt.uniHunks, uniHunks, cmp.Compare) {
			t.Fatalf(":%s:uniHunks: want: %v, got: %v", tt.name, tt.uniHunks, tt.uniHunks)
		}
	}

}

func TestPluralDiff(t *testing.T) {
	a := []rune("abc")
	b := []rune("abd")
	diff := New(a, b)
	diff.SetRouteSize(1)
	diff.Compose()
	lcs := string(diff.Lcs())
	sesActual := diff.Ses()
	sesExpected := []SesElem[rune]{
		{elem: 'a', typ: SesCommon, aIdx: 1, bIdx: 1},
		{elem: 'b', typ: SesCommon, aIdx: 2, bIdx: 2},
		{elem: 'c', typ: SesDelete, aIdx: 3, bIdx: 0},
		{elem: 'd', typ: SesAdd, aIdx: 0, bIdx: 3},
	}

	if diff.EditDistance() != 2 {
		t.Fatalf("want: 2, actual: %v", diff.EditDistance())
	}

	if lcs != "ab" {
		t.Fatalf("want: ab, actual: %v", lcs)
	}

	if !equalsSesElemSlice(sesActual, sesExpected, func(se1, se2 SesElem[rune]) int { return se1.Cmp(se2, cmp.Compare) }) {
		t.Fatalf("want: %v, actual: %v", sesExpected, sesActual)
	}

	uniHunksActual := diff.UnifiedHunks()
	uniHunksExpected := []UniHunk[rune]{
		{a: 1, b: 3, c: 1, d: 3, changes: sesExpected},
	}

	if !equalsUniHunks(uniHunksActual, uniHunksExpected, cmp.Compare) {
		t.Fatalf(":uniHunks: want: %v, got: %v", uniHunksExpected, uniHunksActual)
	}
}

func TestDiffOnlyEditDistance(t *testing.T) {
	a := []rune("abc")
	b := []rune("abd")
	diff := New(a, b)
	diff.OnlyEd()
	diff.Compose()
	lcs := string(diff.Lcs())
	sesActual := diff.Ses()
	sesExpected := []SesElem[rune]{}
	uniHunksActual := diff.UnifiedHunks()
	uniHunksExpected := []UniHunk[rune]{}

	if diff.EditDistance() != 2 {
		t.Fatalf("want: 2, actual: %v", diff.EditDistance())
	}

	if lcs != "" {
		t.Fatalf("want: \"\", actual: %v", lcs)
	}

	if !equalsSesElemSlice(sesActual, sesExpected, func(se1, se2 SesElem[rune]) int { return se1.Cmp(se2, cmp.Compare) }) {
		t.Fatalf("want: %v, actual: %v", sesExpected, sesActual)
	}

	if !equalsUniHunks(uniHunksActual, uniHunksExpected, cmp.Compare) {
		t.Fatalf(":uniHunks: want: %v, got: %v", uniHunksExpected, uniHunksActual)
	}
}

func TestDiffPluralSubsequence(t *testing.T) {
	a := []rune("abcaaaaaabd")
	b := []rune("abdaaaaaabc")
	diff := New(a, b)
	diff.SetRouteSize(2) // dividing sequences forcibly
	diff.Compose()
	if diff.EditDistance() != 4 {
		t.Fatalf("want: 4, actual: %d", diff.EditDistance())
	}
}

func TestDiffSprintSes(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"a", "1", "c"}
	diff := New(a, b)
	diff.Compose()
	actual := diff.SprintSes()
	expected := ` a
-b
+1
 c
`
	if actual != expected {
		t.Fatalf("want: %v, actual: %v", expected, actual)
	}
}

func TestDiffSprintUniHunks(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"a", "1", "c"}
	diff := New(a, b)
	diff.Compose()
	actual := SprintUniHunks(diff.UnifiedHunks())
	expected := `@@ -1,3 +1,3 @@
 a
-b
+1
 c
`
	if actual != expected {
		t.Fatalf("want: %v, actual: %v", expected, actual)
	}
}

func BenchmarkStringDiffCompose(b *testing.B) {
	s1 := []rune("abc")
	s2 := []rune("abd")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff := New(s1, s2)
		diff.Compose()
	}
}

func BenchmarkStringDiffComposeIfOnlyEd(b *testing.B) {
	s1 := []rune("abc")
	s2 := []rune("abd")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff := New(s1, s2)
		diff.OnlyEd()
		diff.Compose()
	}
}

func BenchmarkStringUnifiledHunks(b *testing.B) {
	s1 := []rune("abc")
	s2 := []rune("abd")
	diff := New(s1, s2)
	diff.Compose()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = diff.UnifiedHunks()
	}
}
