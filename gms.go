package main

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
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

func isDirectory(path string) (bool, error) {
    // Returns true if the given path is a directory
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("Cannot open file %s\n", path)
        return false, err
        // Dunno? Return 404?
    }

    fi, err := f.Stat()
    defer f.Close()
    if err != nil {
        return false, err
        fmt.Printf("Error Stating file %s\n", path)
        // How could this error???
    }

    var dir bool = false
    if fi.Mode().IsDir() {
        return true, nil
    }

    return dir, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

    // Get a full path to the destanation file
    wd, _ := os.Getwd()
    working_dir, _ := filepath.Abs(wd)
    var request string = working_dir + r.URL.Path
    fmt.Printf("%s\n",request)

    mode, err := isDirectory(request)
    if err != nil {
        // 404 Should occur here somehow
        return
    }

    if mode {
        // ====================================================================
        // ==== Stuff for a Directory

        // Check if there is a trailing slash.  If not, redirect

        // At this point you can assume there is a slash

        // Get a directory listingp
        files, dir_read_error := ioutil.ReadDir(request)
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
            dir, err:=isDirectory(request+f.Name())
            if err != nil {
                fmt.Printf("Error reading file: %s", request+f.Name())
                // How can this happen?
            }
            if dir {
                fmt.Fprintf(w,"    %s\n",tagify(linkify(f.Name()+"/"),"li", ""))
            } else {
                fmt.Fprintf(w,"    %s\n",tagify(linkify(f.Name()),"li", ""))
            }
        }
        fmt.Fprintf(w,"  </ul>\n")
        fmt.Fprintf(w,"</body>\n</html>")
    } else {
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
