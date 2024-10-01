package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	URL      string
	Method   string
	Body     string
	Interval time.Duration
}

var (
	a           fyne.App
	tasks       []Task
	lock        sync.RWMutex
	serverCache string
)

func main() {
	a = app.New()
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("HTTP Load Simulator", fyne.NewMenuItem("Show", showMain))
		desk.SetSystemTrayMenu(m)
	}

	go worker()
	showMain()
	a.Run()
}

func showMain() {
	w := a.NewWindow("HTTP Load Simulator")
	w.Resize(fyne.NewSize(600, 400))
	vbox := container.NewVBox(
		container.NewHBox(
			widget.NewButton("Add", func() {
				addLoad(w)
			}),
		),
		widget.NewTableWithHeaders(
			func() (rows int, cols int) {
				return len(tasks), 1
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("task")
			},
			func(tci widget.TableCellID, co fyne.CanvasObject) {
				co.(*widget.Label).SetText(tasks[tci.Row].URL)
			},
		),
	)
	w.SetContent(vbox)
	w.Show()
}

func worker() {
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			fmt.Println(tasks)
			for _, task := range tasks {
				http.Get(task.URL)
			}
		}
	}
}

func addLoad(w fyne.Window) {
	url := binding.NewString()
	dialog.ShowForm("Add New Request", "Save", "Cancel", []*widget.FormItem{
		widget.NewFormItem("URL", widget.NewEntryWithData(url)),
	}, func(b bool) {
		if b {
			s, e := url.Get()
			if e != nil {
				log.Println(e)
				return
			}
			tasks = append(tasks, Task{
				URL: s,
			})
		}
	}, w)
}
