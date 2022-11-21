package outputs

import (
	"context"
	"fmt"
	"main/core"
	"time"

	"github.com/rivo/tview"
)

type Table struct {
	app   *tview.Application
	table *tview.Table
}

func NewTable() (*Table, error) {
	return &Table{}, nil
}

func (c *Table) receiveUpdates(s core.Session, r core.Reading) {
	if c.table == nil || c.app == nil {
		return
	}
	row := c.table.GetRowCount()
	c.table.SetCell(row, 0, tview.NewTableCell(r.Received.Format(time.Stamp)))
	c.table.SetCell(row, 1, tview.NewTableCell(fmt.Sprint(r.Temp1)))
	c.table.SetCell(row, 3, tview.NewTableCell(fmt.Sprint(r.Temp2)))
	c.app.Draw()
}

func (c *Table) Listener() core.Listener {
	return c.receiveUpdates
}

func (c *Table) Start(ctx context.Context, outChan chan error) {
	c.app = tview.NewApplication()
	c.table = tview.NewTable().SetBorders(true)
	err := c.app.SetRoot(c.table, true).EnableMouse(true).Run()
	outChan <- err
}
