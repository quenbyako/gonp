package gonp_test

import (
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"log"
	"os"
	"time"
	"unicode/utf8"

	"github.com/quenbyako/gonp"
)

func ExampleDiff_Patch_filePatch() {
	if len(os.Args) < 3 {
		log.Fatal("./filepatch filename1 filename2")
	}

	f1 := os.Args[1]
	f2 := os.Args[2]

	var (
		a   []string
		b   []string
		err error
	)

	a, err = getLines(f1)
	if err != nil {
		log.Fatalf("%s: %s", f1, err)
	}

	b, err = getLines(f2)
	if err != nil {
		log.Fatalf("%s: %s", f2, err)
	}

	diff := gonp.New(a, b)
	diff.Compose()

	patchedSeq := diff.Patch(a)
	fmt.Printf("success:%v, applying SES between '%s' and '%s'\n", equalsStringSlice(b, patchedSeq), f1, f2)

	uniPatchedSeq, err := diff.UniPatch(a, diff.UnifiedHunks())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("success:%v, applying unified format difference between '%s' and '%s'\n", equalsStringSlice(b, uniPatchedSeq), f1, f2)
}

func ExampleDiff_PrintSes_intDiff() {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 9, 4, 5}
	diff := gonp.New(a, b)
	diff.Compose()
	fmt.Printf("diff %v %v\n", a, b)
	fmt.Printf("EditDistance: %d\n", diff.EditDistance())
	fmt.Printf("LCS: %v\n", diff.Lcs())
	fmt.Println("SES:")
	diff.PrintSes()
}

func ExampleDiff_PrintSes_strDiff() {
	if len(os.Args) < 3 {
		log.Fatal("./strdiff arg1 arg2")
	}
	if !utf8.ValidString(os.Args[1]) {
		log.Fatalf("arg1 contains invalid rune")
	}

	if !utf8.ValidString(os.Args[2]) {
		log.Fatalf("arg2 contains invalid rune")
	}
	a := []rune(os.Args[1])
	b := []rune(os.Args[2])
	diff := gonp.New(a, b)
	diff.Compose()
	fmt.Printf("EditDistance: %d\n", diff.EditDistance())
	fmt.Printf("LCS: %s\n", string(diff.Lcs()))
	fmt.Println("SES:")

	var buf bytes.Buffer
	ses := diff.Ses()
	for _, e := range ses {
		ee := e.GetElem()
		switch e.GetType() {
		case gonp.SesDelete:
			fmt.Fprintf(&buf, "-%c\n", ee)
		case gonp.SesAdd:
			fmt.Fprintf(&buf, "+%c\n", ee)
		case gonp.SesCommon:
			fmt.Fprintf(&buf, " %c\n", ee)
		}
	}
	fmt.Print(buf.String())
}

func ExampleDiff_UniPatch_strPatch() {
	if len(os.Args) < 3 {
		log.Fatal("./strpatch arg1 arg2")
	}
	if !utf8.ValidString(os.Args[1]) {
		log.Fatal("arg1 contains invalid rune")
	}

	if !utf8.ValidString(os.Args[2]) {
		log.Fatal("arg2 contains invalid rune")
	}
	a := []rune(os.Args[1])
	b := []rune(os.Args[2])
	diff := gonp.New(a, b)
	diff.Compose()

	patchedSeq := diff.Patch(a)
	fmt.Printf("success:%v, applying SES between '%s' and '%s' to '%s' is '%s'\n",
		string(b) == string(patchedSeq),
		string(a), string(b),
		string(a), string(patchedSeq))

	uniPatchedSeq, err := diff.UniPatch(a, diff.UnifiedHunks())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("success:%v, applying unified format difference between '%s' and '%s' to '%s' is '%s'\n",
		string(b) == string(uniPatchedSeq),
		string(a), string(b),
		string(a), string(uniPatchedSeq))
}

func ExampleDiff_PrintUniHunks_uniFileDiff() {
	if len(os.Args) < 3 {
		log.Fatal("./unifilediff filename1 filename2")
	}

	f1 := os.Args[1]
	f2 := os.Args[2]

	var (
		a   []string
		b   []string
		err error
	)

	a, err = getLines(f1)
	if err != nil {
		log.Fatalf("%s: %s", f1, err)
	}

	b, err = getLines(f2)
	if err != nil {
		log.Fatalf("%s: %s", f2, err)
	}

	th, err := buildTargetHeader(f1, f2)
	if err != nil {
		log.Fatal(err)
	}

	diff := gonp.New(a, b)
	diff.Compose()

	fmt.Printf(th.String())
	diff.PrintUniHunks(diff.UnifiedHunks())
}

func ExampleDiff_Patch_uniIntDiff() {
	a := []Row{{1, "Pupa"}, {2, "Lupa"}, {3, "Popa"}}
	b := []Row{{1, "Pupa"}, {2, "Lupa"}, {3, "Zhopa"}}
	diff := gonp.NewCmp(a, b, cmpRow)
	diff.Compose()
	fmt.Printf("diff %v %v\n", a, b)
	fmt.Printf("EditDistance: %d\n", diff.EditDistance())
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

func ExampleDiff_UniPatch_uniStrDiff() {
	if len(os.Args) < 3 {
		log.Fatal("./unistrdiff arg1 arg2")
	}
	if !utf8.ValidString(os.Args[1]) {
		log.Fatalf("arg1 contains invalid rune")
	}

	if !utf8.ValidString(os.Args[2]) {
		log.Fatalf("arg2 contains invalid rune")
	}
	a := []rune(os.Args[1])
	b := []rune(os.Args[2])
	diff := gonp.New(a, b)
	diff.Compose()
	fmt.Printf("EditDistance:%d\n", diff.EditDistance())
	fmt.Printf("LCS:%s\n", string(diff.Lcs()))
	//diff.PrintUniHunks(diff.UnifiedHunks())

	fmt.Println("Unified format difference:")
	uniHunks := diff.UnifiedHunks()
	var w bytes.Buffer
	for _, uniHunk := range uniHunks {
		fmt.Fprintf(&w, uniHunk.SprintDiffRange())
		for _, e := range uniHunk.GetChanges() {
			switch e.GetType() {
			case gonp.SesDelete:
				fmt.Fprintf(&w, "-%c\n", e.GetElem())
			case gonp.SesAdd:
				fmt.Fprintf(&w, "+%c\n", e.GetElem())
			case gonp.SesCommon:
				fmt.Fprintf(&w, " %c\n", e.GetElem())
			}
		}
	}
	fmt.Print(w.String())
}

func equalsStringSlice(a, b []string) bool {
	m, n := len(a), len(b)
	if m != n {
		return false
	}
	for i := 0; i < m; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Target consists of a path and mtime of file.
type Target struct {
	fname string
	mtime time.Time
}

// TargetHeader has 2 targets based on pathes and mtimes based on 2 files
type TargetHeader struct {
	targets []Target
}

// getLines returns a file contents as string array
func getLines(f string) ([]string, error) {
	fp, err := os.Open(f)
	if err != nil {
		return []string{}, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

// builderTargetHeader returns TargetHeader constructed based on 2 files given as arguments
func buildTargetHeader(f1, f2 string) (TargetHeader, error) {
	fi1, err := os.Stat(f1)
	if err != nil {
		return TargetHeader{}, err
	}
	fi2, err := os.Stat(f2)
	if err != nil {
		return TargetHeader{}, err
	}
	return TargetHeader{
		targets: []Target{
			{fname: f1, mtime: fi1.ModTime()},
			{fname: f2, mtime: fi2.ModTime()},
		},
	}, nil
}

// String returns a content of TargetHeader as a string
func (th *TargetHeader) String() string {
	if len(th.targets) != 2 {
		return ""
	}

	var b bytes.Buffer
	fmt.Fprintf(&b, "--- %s\t%s\n", th.targets[0].fname, th.targets[0].mtime.Format(time.RFC3339Nano))
	fmt.Fprintf(&b, "+++ %s\t%s\n", th.targets[1].fname, th.targets[1].mtime.Format(time.RFC3339Nano))
	return b.String()
}
