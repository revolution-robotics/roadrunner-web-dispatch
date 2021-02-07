package main

import (
	"flag"
	"fmt"
	"net/http"
	"github.com/vharitonsky/iniflags"
	"os"
	"os/exec"
	"github.com/alessio/shellescape"
	"strconv"
)

var (
	// iniflags() defines:
	//     --allowMissingConfig
	//     --allowUnknownFlags
	//     --config=configPath
	//     --configUpdateInterval=duration
	//     --dumpflags
	cgiFlag  = flag.String("cgi", "/usr/bin/status.py", "CGI path")
	portFlag = flag.Int("port", 80, "Port to listen on")
	uriFlag  = flag.String("uri", "/status.json", "URI of CGI trigger")
	wwwFlag  = flag.String("www", "/var/www/html", "HTML directory")
)

func execHandler(w http.ResponseWriter, r *http.Request) {
	reqURI := shellescape.Quote(r.URL.Path)
	cmd := &exec.Cmd {
		Path: *cgiFlag,
		Args: []string{ *cgiFlag,  reqURI },
		Stdout: w,
		Stderr: os.Stderr,
	}
	fmt.Printf("Exec: %s\n", *cgiFlag)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	newURI := "https://" + r.Host + ":9090"
	fmt.Printf("Redirect: %s\n", newURI)
	http.Redirect(w, r, newURI, http.StatusSeeOther)
}

func main() {
	iniflags.Parse()

	port := ":" + strconv.Itoa(*portFlag)
	fmt.Printf("Listening on port %s\n", port)

	http.Handle("/", http.FileServer(http.Dir(*wwwFlag)))
	http.HandleFunc(*uriFlag, execHandler)
	http.HandleFunc("/cockpit", redirectHandler)
	http.ListenAndServe(port , nil)
}
