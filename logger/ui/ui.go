package ui

import (
	"context"
	"fmt"
	"main/core"
	"main/datautil"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/guptarohit/asciigraph"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type UI struct {
	updateChan chan core.Reading
	session    core.Session
	device     core.Device
	log        *logrus.Logger
	plotBox    *tview.TextView
	app        *tview.Application
	grid       *tview.Grid
	table      *tview.Table
	form       *tview.Form
	LogView    *tview.TextView
}

func New(log *logrus.Logger, sess core.Session, device core.Device) (*UI, error) {
	ui := &UI{
		session:    sess,
		device:     device,
		updateChan: make(chan core.Reading),
		log:        log,
		plotBox: tview.NewTextView().
			SetScrollable(false).
			SetWrap(false),
		app: tview.NewApplication(),
		grid: tview.NewGrid().
			SetRows(0, 0).
			SetColumns(0, 0),
		table: tview.NewTable().
			SetSeparator(tview.Borders.Vertical),
		form:    tview.NewForm(),
		LogView: tview.NewTextView(),
	}

	ui.plotBox.SetTitle("Plot")
	ui.plotBox.SetBorder(true)

	ui.table.SetTitle("Readings")
	ui.table.SetBorder(true)

	ui.LogView.SetTitle("System Log")
	ui.LogView.SetBorder(true)

	ui.form.SetTitle("Controls")
	ui.form.SetBorder(true)

	err := ui.setupForm()
	if err != nil {
		return nil, err
	}

	ui.grid.
		AddItem(ui.plotBox, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.table, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.LogView, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.form, 1, 1, 1, 1, 0, 0, true)

	return ui, nil
}

func (c *UI) setupForm() error {
	md, err := c.session.GetMetadata()
	if err != nil {
		return err
	}
	c.form.Clear(true)

	tempMd := md
	c.form.AddInputField("Food", md.Food, 0, nil, func(text string) {
		tempMd.Food = text
	})
	c.form.AddInputField("Method", md.Method, 0, nil, func(text string) {
		tempMd.Method = text
	})

	var tempTemperatures [2]string
	tempTemperatures[0] = fmt.Sprint(md.Calibrations[0])
	tempTemperatures[1] = fmt.Sprint(md.Calibrations[1])
	c.form.AddInputField("Calibration 1", tempTemperatures[0], 0, nil, func(text string) {
		tempTemperatures[0] = text
	})
	c.form.AddInputField("Calibration 2", tempTemperatures[1], 0, nil, func(text string) {
		tempTemperatures[1] = text
	})

	c.form.AddButton("Save", func() {
		var calib [2]float64

		calib[0], err = strconv.ParseFloat(tempTemperatures[0], 64)
		if err != nil {
			c.log.Error(err)
			return
		}

		calib[1], err = strconv.ParseFloat(tempTemperatures[1], 64)
		if err != nil {
			c.log.Error(err)
			return
		}

		tempMd.Calibrations = calib

		err := c.session.SetMetadata(tempMd)
		if err != nil {
			c.log.Error(err)
			return
		}

	})

	return nil
}

func (c *UI) receiveUpdates(s core.Session, r core.Reading) {
	c.updateChan <- r
}

func (c *UI) Listener() core.Listener {
	return c.receiveUpdates
}

func (c *UI) updateTable(reading core.Reading) {
	row := c.table.GetRowCount()
	c.table.SetCell(row, 0, tview.NewTableCell(reading.Received.Format(time.Stamp)))
	c.table.SetCell(row, 1, tview.NewTableCell(fmt.Sprint(reading.Temperatures[0])))
	c.table.SetCell(row, 3, tview.NewTableCell(fmt.Sprint(reading.Temperatures[1])))
}

func (c *UI) updatePlot() {
	readings, err := c.session.GetReadings()
	if err != nil {
		c.app.QueueEvent(tcell.NewEventError(err))
		return
	}
	_, _, width, height := c.plotBox.GetInnerRect()
	data := datautil.NormalizeTimeDistribution(readings, width)
	graph := asciigraph.PlotMany(data, asciigraph.Width(width), asciigraph.Height(height))
	c.plotBox.SetText(graph)
}

func (c *UI) Start(ctx context.Context, errChan chan error) {
	go func() {
		for {
			select {
			case err := <-errChan:
				c.app.QueueEvent(tcell.NewEventError(err))
			case reading := <-c.updateChan:
				c.updateTable(reading)
				c.updatePlot()
				c.app.Draw()
			}
		}
	}()

	if err := c.app.SetRoot(c.grid, true).Run(); err != nil {
		c.log.Error(err)
	}
}
