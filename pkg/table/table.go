package table

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/logrusorgru/aurora/v3"
)

type Row []string

type Table struct {
	head Row
	rows []Row
}

func NewTable(cols []string) *Table {
	return &Table{
		head: cols,
	}
}

func (t *Table) SetRows(rows []Row) {
	t.rows = rows
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)

	// Print head table
	head := strings.Join(t.head, "\t")
	fmt.Fprintln(w, aurora.Gray(20, head).Bold())

	// Print all data table
	for _, item := range t.rows {
		row := strings.Join(item, "\t")
		fmt.Fprintln(w, aurora.Gray(10, row).Faint())
	}

	w.Flush()
	fmt.Print("\n")
}
