package main

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
	"text/template"
	"time"
)

var visits int = 0

type dirNode struct {
	Name    string
	Size    int64
	//Mode    FileMode
	ModTime time.Time
	IsDir   bool
}

type fileNode struct {
	Visits  int
	Path    string
	Files   []*dirNode
}

func getHostname() (string) {
    hostname, _ := os.Hostname()
    return hostname
}

func isDirectory(path string) (bool, error) {
    // Returns true if the given path is a directory

    fi, err := os.Stat(path)
	if err != nil {
        return false, err
        fmt.Printf("Error Stating file %s\n", path)
    }

	return fi.IsDir(), nil
}

func makeFileStruct(f os.FileInfo) (*dirNode) {
	r := new(dirNode)
	r.Name = f.Name()
	r.Size = f.Size()
	r.ModTime = f.ModTime()
	r.IsDir = f.IsDir()
	return r
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

        // Get a directory listing
        files, dir_read_error := ioutil.ReadDir(request)
        if dir_read_error != nil {
            fmt.Printf("Error reading %s: %d\n", r.URL.Path, dir_read_error)
        }

		// Parse and check the template
		t := template.Must(template.ParseFiles("http/dirview.html"))
		if err != nil {fmt.Printf("Error opening dirview.html\n")}

		var Nodes []*dirNode = make([]*dirNode, len(files))
		var i int = 0
		for _, f := range files {

			// Delegate the actual processing to another thing
			Nodes[i] = makeFileStruct(f)
			i++
        }
		
		// Construct the dirStruct!
		dirStruct := fileNode {
			Visits: visits,
			Path: r.URL.Path,
			Files: Nodes,
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
