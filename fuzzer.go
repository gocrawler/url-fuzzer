package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	MAXX_SPEED     int
	FILE_NAME      string
	TotalCollected []string
	Fuzzing        int32
	Collection     int32
	OpenConns      int32
	Errs           []string
	mu             sync.Mutex
	file           *os.File
)

func main() {
	fmt.Println(`URL FUZZER!
By AnikHasibul (@AnikHasibul)

For Live View 

>>> http://localhost:1339`)

	// parse flags

	flag.IntVar(
		&MAXX_SPEED,
		"max",
		50,
		"Max crawling speed",
	)

	//

	flag.StringVar(
		&FILE_NAME,
		"o",
		"results.urlfuzzer.txt",
		"Output to file",
	)
	//
	flag.Parse()

	//
	//

	var err error
	file, err = os.OpenFile(
		FILE_NAME,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0600,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Let it be fuzzy
	go func() {
		var wg sync.WaitGroup
		sites := ReadSite()
		paths := ReadPath()
		log.Println(
			" %d SITES with %d PATHS...\n",
			len(sites),
			len(paths),
		)
		for _, site := range sites {
			if site == "" {
				continue
			}
			for _, path := range paths {
				if path == "" {
					continue
				}
				wg.Add(1)
				conn := atomic.LoadInt32(
					&OpenConns,
				)
				for conn > int32(
					MAXX_SPEED,
				) {
					log.Println(
						"Max Speed Reached",
						"REQUEST:",
						conn,
						"OPEN CONN:",
						atomic.LoadInt32(
							&Fuzzing,
						),
					)
					time.Sleep(
						3 * time.Second,
					)

					if conn < int32(
						MAXX_SPEED,
					) {
						break
					}
				}

				go Fuzz(
					site,
					path,
					&wg,
				)

			}

		}
		wg.Wait()
		log.Println("Finished!")
		log.Println("Saving index.html")
		log.Println("http://localhost:1339")
		time.Sleep(5 * time.Second)
	}()
	http.HandleFunc("/", live)
	http.ListenAndServe(":1339", nil)

}

// Fuzz checks http status of given url
func Fuzz(site, path string, wg *sync.WaitGroup) {
	conn := atomic.LoadInt32(
		&OpenConns,
	)
	atomic.StoreInt32(
		&OpenConns,
		conn+1,
	)
	defer func() {
		// remove 1 from conn queue
		running := atomic.LoadInt32(
			&Fuzzing,
		)
		atomic.StoreInt32(
			&Fuzzing,
			running-1,
		)
		atomic.StoreInt32(
			&OpenConns,
			conn-1,
		)
		wg.Done()
	}()

	// Add 1 to connection queue
	atomic.AddInt32(
		&Collection,
		1,
	)
	atomic.AddInt32(
		&Fuzzing,
		1,
	)

	// Validate url (host)
	if !strings.HasPrefix(site, "http") {
		fmt.Println(
			site,
			"is not a valid url",
		)
		return
	}
	// validate url (path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Open a request
	resp, err := http.Get(site + path)
	if err != nil {
		mu.Lock()
		Errs = append(Errs, fmt.Sprint(err))
		mu.Unlock()
		return
	}
	defer resp.Body.Close()

	// url OK?
	// Add to the collection!
	if resp.StatusCode == 200 {
		mu.Lock()
		TotalCollected = append(
			TotalCollected,
			site+path,
		)
		file.WriteString(
			strings.Join(
				TotalCollected,
				"\n",
			),
		)
		err := file.Sync()
		if err != nil {
			log.Fatal(err)
		}
		mu.Unlock()
	}
}

// ReadSite reads site list from file
func ReadSite() []string {
	sites, err := ioutil.ReadFile("sites.txt")
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(sites), "\n")
}

// ReadPath reads path list from file
func ReadPath() []string {
	paths, err := ioutil.ReadFile("lists.txt")
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(paths), "\n")
}

// live generates webpage for web interface
func live(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	// Template model
	type gen struct {
		PageTitle string
		Links     []string
	}

	// Fill struct
	ser := gen{
		PageTitle: fmt.Sprintf(
			`UrlFuzzer :: By Anik Hasibul 
			[Tried for %d Times] 
			RUNNING FOR [%s] Found (%d)
			in %d Tries!\n`,
			atomic.LoadInt32(
				&Collection,
			),
			time.Since(start),
			len(TotalCollected),
			atomic.LoadInt32(
				&Collection,
			),
		),
		Links: TotalCollected,
	}

	// Fill template
	t, ert := template.New("View").Parse(`

	<h1>{{.PageTitle}}<h1>
<ul>
    {{range .Links}}
           <li> <a href="{{.}}">{{.}}</a></li><br>
    {{end}}
</ul>`)
	// Hane error!
	if ert != nil {
		log.Println(
			"Parse error:",
			ert,
		)
	}
	err := t.ExecuteTemplate(
		w,
		"View",
		ser,
	)
	if err != nil {
		log.Println(
			"Execution error:",
			err,
		)
	}
}
