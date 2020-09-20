package main

import (
	"log"
	"os"

	"gioui.org/app"
	"github.com/yarcat/playground-gio/transition-app/res"
)

func main() {
	go func() {
		app := newTransitionApp(
			mustDecodePNG(res.GopherSimplePNG),
			mustDecodePNG(res.GopherPNG),
		)
		if err := app.mainloop(); err != nil {
			log.Fatal("mainloop failed:", err)
		}
		os.Exit(0)
	}()
	app.Main()
}
