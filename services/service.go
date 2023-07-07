package services

import (
	"EnronEmailApi/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

func IndexStart() {
	var wg sync.WaitGroup

	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)

	defer pprof.StopCPUProfile()

	path := "c:/Users/ludmi/OneDrive/Documentos/enron_mail_20110402/maildir/"

	fmt.Println("Indexando...")

	userList, _ := ListAllFolders(path)

	numParts := 5
	dividedFolders := DivideFolders(userList, numParts)

	wg.Add(5)

	for index := range dividedFolders {
		go Algodeaca(dividedFolders[index], path, &wg)
	}

	wg.Wait()

	fmt.Println("Indexing finished!!!!")

	runtime.GC()
	mem, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal(err)
	}

}

func SearchEmails(text *string) models.EmailResponse {

	var respuesta models.EmailResponse
	now := time.Now()
	startTime := now.AddDate(0, 0, -7).Format("2006-01-02T15:04:05Z07:00")
	endTime := now.Format("2006-01-02T15:04:05Z07:00")

	fmt.Println(*text)

	query := fmt.Sprintf(`{
		"search_type": "match",
		"query": {
			"term":       "`+*text+`, ",
			"start_time": "%s",
			"end_time":   "%s"
		},
		"from":        0,
		"max_results": 20,
		"_source":     []
	}`, startTime, endTime)

	req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "password")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &respuesta)

	return respuesta

}
