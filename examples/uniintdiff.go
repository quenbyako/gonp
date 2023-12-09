package main

import (
	"cmp"
	"fmt"

	"github.com/quenbyako/gonp"
)

func main() {
	a := []Row{{1, "Pupa"}, {2, "Lupa"}, {3, "Popa"}}
	b := []Row{{1, "Pupa"}, {2, "Lupa"}, {3, "Zhopa"}}
	diff := gonp.NewCmp(a, b, cmpRow)
	diff.Compose()
	fmt.Printf("diff %v %v\n", a, b)
	fmt.Printf("Editdistance: %d\n", diff.Editdistance())
	fmt.Printf("LCS: %v\n", diff.Lcs())
	fmt.Println("Unified format difference:")
	diff.PrintUniHunks(diff.UnifiedHunks())
}

type Row struct {
	ID   int
	Name string
}

func (r *Row) String() string {
	return fmt.Sprintf("%v, %v", r.ID, r.Name)
}

func cmpRow(a, b Row) int {
	for _, v := range []int{
		cmp.Compare(a.ID, b.ID),
		cmp.Compare(a.Name, b.Name),
	} {
		if v != 0 {
			return v
		}
	}

	return 0
}
