package main

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
	"text/template"
)

var visits int = 0

type directory struct {
	Visits  int
	Path    string
	Files   []string
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
		w.WriteHeader(404)
        // 404 Should occur here somehow
        return
    }

    if mode {
        // ====================================================================
        // ==== Stuff for a Directory

        // Check if there is a trailing slash.  If not, redirect
        if !strings.HasSuffix(request, "/") {
            http.Redirect(w, r, r.URL.Path+"/", 300)
        }

        // At this point you can assume there is a slash

        // Get a directory listingp
        files, dir_read_error := ioutil.ReadDir(request)
        if dir_read_error != nil {
            fmt.Printf("Error reading %s: %d\n", r.URL.Path, dir_read_error)
        }

		// Now this part uses go templates.  Need to work out how a
		// project is supposed to lay out its random files.  For now,
		// all web page related stuff is in http, like html, css, and
		// js
		
		// Parse and check the template
		t := template.Must(template.ParseFiles("http/dirview.html"))
		if err != nil {fmt.Printf("Error opening dirview.html\n")}

		var filenames []string = make([]string, len(files))

		var i int = 0

		// Translate the file objects into file names
		// Todo - Use a dile struct with more data - file size,
		// modification date, etc
		for _, f := range files {
            dir, err:=isDirectory(request+f.Name())
            if err != nil {
                fmt.Printf("Error reading file: %s", request+f.Name())
                // How can this happen?
            }
			filenames[i] = f.Name()
            if dir {
                //fmt.Fprintf(w,"
				filenames[i] = filenames[i]+"/"
			} 
			i++
        }
		
		dirStruct := directory {
			Visits: visits,
			Path: r.URL.Path,
			Files: filenames,
		}

		// Execute the template, write to web server
		err := t.Execute(w, dirStruct)
		if err != nil {fmt.Printf("Error executing the template\n")}
		
    } else {
        http.ServeFile(w, r, request)
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
