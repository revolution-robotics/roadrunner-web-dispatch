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
	"strings"
)

var (
	// iniflags() defines:
	//     --allowMissingConfig
	//     --allowUnknownFlags
	//     --config=configPath
	//     --configUpdateInterval=duration
	//     --dumpflags
	cgiFlag	      = flag.String("cgi", "/usr/bin/status.py", "CGI path")
	httpPortFlag  = flag.Int("httpPort", 80, "HTTP port to listen on")
	httpsPortFlag = flag.Int("httpsPort", 443, "HTTPS port to listen on")
	uriFlag	      = flag.String("uri", "/status.json", "URI of CGI trigger")
	wwwFlag	      = flag.String("www", "/var/www/html", "HTML directory")
	certFlag      = flag.String("cert",
		"/etc/web-dispatch/certs/self-signed.crt",
		"Public TLS certificate")
	keyFlag	      = flag.String("key",
		"/etc/web-dispatch/certs/self-signed.key",
		"Private TLS Key")
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

func cockpitHandler(w http.ResponseWriter, r *http.Request) {
	newURI := "https://" + r.Host + ":9090" +
		strings.Replace(r.RequestURI, "/cockpit", "", 1)
	fmt.Printf("Redirect: %s\n", newURI)
	http.Redirect(w, r, newURI, http.StatusSeeOther)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	newURI := "https://" + r.Host + r.RequestURI
	fmt.Printf("Redirect: %s\n", newURI)
	http.Redirect(w, r, newURI, http.StatusSeeOther)
}

func main() {
	iniflags.Parse()

	httpPort := ":" + strconv.Itoa(*httpPortFlag)
	httpsPort := ":" + strconv.Itoa(*httpsPortFlag)
	fmt.Printf("Listening on ports %s and %s\n", httpPort, httpsPort)

	http.Handle("/", http.FileServer(http.Dir(*wwwFlag)))
	http.HandleFunc(*uriFlag, execHandler)
	http.HandleFunc("/cockpit/", cockpitHandler)

	go func() {
		err := http.ListenAndServe(httpPort,
			http.HandlerFunc(redirectHandler))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
	}()

	err := http.ListenAndServeTLS(httpsPort, *certFlag, *keyFlag, nil)
	if  err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
