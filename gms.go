package main

import (
        "io/ioutil"
        "fmt"
        "net/http"
        "os"
        "path/filepath"
)

var visits int = 0

func linkify(filename string) (string) {
        var r string = ""
        r += "<a href=\"" + filename + "\">" + filename + "</a>"
        return r
}

func tagify(text string, tag string, attributes string) (string) {
        var r string = "<" + tag + " " + attributes + ">" + text + "</" + tag + ">"
        return r
}

func getHostname() (string) {
        hostname, _ := os.Hostname()
        return hostname
}

func isDirictery(path string) (bool) {
        f, err := os.Open(path)
        if err != nil {
                fmt.Printf("Cannot open file %s\n", request)
                // Dunno? Return 404?
        }

        fi, err := f.Stat()
        if err != nil {
                fmt.Printf("Error Stating file %s\n", request)
                // How could this error???
        }

        switch mode := fi.mode(); {
        case mode.IsDir():
                return true
        case mode.IsRegular():
                return false
        }
	f.close()
}

func handler(w http.ResponseWriter, r *http.Request) {

        // Get a full path to the destanation file
        wd, _ := os.Getwd()
        working_dir, _ := filepath.Abs(wd)
        var request string = working_dir + r.URL.Path
        fmt.Printf("%s\n",request)

        // Detect file type
        f, err := os.Open(request)
        if err != nil {
                fmt.Printf("Cannot open file %s\n", request)
                // Dunno? Return 404?
        }
        defer f.Close()
        fi, err := f.Stat()
        if err != nil {
                fmt.Printf("Error Stating file %s\n", request)
                // How could this error???
        }
        switch mode := fi.Mode(); {
        case mode.IsDir():
                // ====================================================================
                // ==== Stuff for a Directory

                // Get a directory listing
                files, dir_read_error := ioutil.ReadDir(working_dir + r.URL.Path)
                if dir_read_error != nil {
                        fmt.Printf("Error reading %s: %d\n", r.URL.Path, dir_read_error)
                }

                fmt.Fprintf(w,"<html>\n")
                fmt.Fprintf(w,"<head>\n")
                fmt.Fprintf(w,"  <script type='text/javascript' src='script.js'></script>\n")
                fmt.Fprintf(w,"</head>\n")
                fmt.Fprintf(w,"<body>\n")
                fmt.Fprintf(w,"  <span id='count'>%d</span> Visits so far\n", visits)
                fmt.Fprintf(w,"  <ul>\n")
                for _, f := range files {
                        fmt.Fprintf(w,"    %s\n",tagify(linkify(f.Name()),"li", ""))
                }
                fmt.Fprintf(w,"  </ul>\n")
                fmt.Fprintf(w,"</body>\n</html>")
        case mode.IsRegular():
                // TODO: read byte by byte
                fb, _ := ioutil.ReadFile(request)
                fmt.Fprintf(w, string(fb))
        }

        //fmt.Printf("Handling connection number %d\n", visits)
        visits += 1
}

func scriptGetter(w http.ResponseWriter, r *http.Request) {
        fi, _ := ioutil.ReadFile("script.js")
        fmt.Fprintf(w, string(fi))
}

func main() {
        http.HandleFunc("/", handler)
        http.HandleFunc("/script.js", scriptGetter)
        http.ListenAndServe(":8080", nil)
}
