package cmd

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/spf13/cobra"
)

//go:embed control
var uipath embed.FS

func openuri(uri string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", uri).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "darwin":
		err = exec.Command("open", uri).Start()
	default:
		err = fmt.Errorf("unsupported platform, open manually: %s", uri)
	}

	if err != nil {
		log.Fatal(err)
	}
}

var UICmd = &cobra.Command{
	Use: "ui",
	Run: func(cmd *cobra.Command, args []string) {
		addr := "localhost:3078"

		handler := http.FileServer(http.FS(uipath))
		fmt.Println(uipath)

		http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			isFile := regexp.MustCompile(`\.[a-z]{2,}$`)
			path := r.URL.Path
			if !isFile.MatchString(path) {
				path = "/"
			}
			r.URL.Path = fmt.Sprintf("/control%s", path)

			handler.ServeHTTP(rw, r)
		})

		openuri(addr)
		fmt.Printf("Numary control is live on http://%s\n", addr)

		http.ListenAndServe(addr, nil)
	},
}
