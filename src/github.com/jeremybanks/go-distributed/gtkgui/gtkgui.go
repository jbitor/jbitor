package gtkgui

import (
	"fmt"
	"github.com/jeremybanks/go-distributed/dht"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"time"
)

const targetNoteCount = 32

func Main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("distributedgtk")

	// XXX: We should probably confirm that $HOME is not empty.
	dhtClientStatePath := os.Getenv("HOME") + "/.distributed-dht.benc"
	dhtClient, err := dht.OpenClient(dhtClientStatePath, false)

	_ = err // XXX

	_ = dhtClient

	notebook := gtk.NewNotebook()
	window.Add(notebook)

	dhtPage := gtk.NewAlignment(0.1, 0.1, 0.1, 0.1)
	notebook.AppendPage(dhtPage, gtk.NewLabel("DHT"))

	progress := gtk.NewProgressBar()
	dhtPage.Add(progress)

	go func() {
		for {
			connectionInfo := dhtClient.ConnectionInfo()
			progress.SetFraction(float64(connectionInfo.GoodNodes) / targetNoteCount)
			progress.SetText(fmt.Sprintf(
				"%v Good Known Nodes\n(%v unknown, %v bad)",
				connectionInfo.GoodNodes,
				connectionInfo.UnknownNodes,
				connectionInfo.BadNodes,
			))
			time.Sleep(5 * time.Second)
		}
	}()

	window.Connect("destroy", func() {
		dhtClient.Save()
		gtk.MainQuit()
	})

	window.ShowAll()
	window.SetSizeRequest(600, 400)

	gtk.Main()
}
