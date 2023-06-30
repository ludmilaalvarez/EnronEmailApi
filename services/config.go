package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	//"encoding/base64"
	"EnronEmailApi/models"
	"EnronEmailApi/zinc"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

var MutexJson sync.Mutex
var JSonGeneral models.JsonFinal

func ListAllFolders(folderName string) ([]string, []string) {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var listFolders []string
	var listFiles []string
	for _, f := range files {
		if f.IsDir() {
			listFiles = append(listFiles, f.Name())
		}
		listFolders = append(listFolders, f.Name())
	}
	return listFiles, listFolders
}

func DivideFolders(list []string, numParts int) [][]string {
	var divided [][]string
	partSize := len(list) / numParts
	remaining := len(list) % numParts

	index := 0
	for i := 0; i < numParts; i++ {
		size := partSize
		if i < remaining {
			size++
		}
		divided = append(divided, list[index:index+size])
		index += size
	}

	return divided
}
func ListFiles(folderName string) []string {
	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var filesNames []string
	for _, file := range files {
		filesNames = append(filesNames, file.Name())
	}
	return filesNames
}

func Algodeaca(folderList []string, path string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, user := range folderList {

		var jsonForBulk models.JsonBulk
		jsonForBulk.Index = "email"

		folders, files := ListAllFolders(path + user)
		for _, file := range files {
			var wg sync.WaitGroup
			wg.Add(1)
			go ProcessMailFile(path+user+"/", file, &wg, &jsonForBulk)
			wg.Wait()
		}

		for _, folder := range folders {

			mailFiles := ListFiles(path + user + "/" + folder + "/")

			for _, mailFile := range mailFiles {
				var wg sync.WaitGroup

				wg.Add(1)
				go ProcessMailFile(path+user+"/"+folder+"/", mailFile, &wg, &jsonForBulk)
				wg.Wait()
			}

		}

		IndexDataBulk(jsonForBulk)
		jsonForBulk.Records = []models.Email{}

	}
}

func ProcessMailFile(path, mailFile string, wg *sync.WaitGroup, jsonForBulk *models.JsonBulk) {
	defer wg.Done()

	//var JSonGeneral models.JsonFinal

	sysFile, err := os.Open(path + mailFile)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}
	defer sysFile.Close()

	lines := bufio.NewScanner(sysFile)

	data := ParseData(lines)

	IndexData(data, jsonForBulk)
	JSonGeneral.Emails = append(JSonGeneral.Emails, data)

}
func ParseData(dataLines *bufio.Scanner) models.Email {
	var data models.Email
	for dataLines.Scan() {
		//data.ID = id
		/* fileContent, err := os.ReadFile(email)
		if err != nil {
			panic(err.Error())
		}
		r := bytes.NewReader(fileContent)
		m, err := mail.ReadMessage(r)
		if err != nil {

			return
		}
		body, err := io.ReadAll(m.Body)
		if err != nil {

		}
		zincData <- fmt.Sprintf(`{"_id": "%s", "from": "%s", "to": "%s", "subject": "%s", "content": "%s"}`,
		email, fmt.Sprintf("%q", m.Header.Get("From")), fmt.Sprintf("%q", m.Header.Get("To")),
		fmt.Sprintf("%q", m.Header.Get("Subject")), fmt.Sprintf("%q", string(body)))*/

		line := dataLines.Text()
		switch {
		case strings.Contains(line, "Message-ID:"):
			data.Message_ID = line[11:]
		case strings.Contains(line, "Date:"):
			data.Date = line[5:]
		case strings.Contains(line, "From:"):
			data.From = line[5:]
		case strings.Contains(line, "To:"):
			data.To = line[3:]
		case strings.Contains(line, "Subject:"):
			data.Subject = line[8:]
		case strings.Contains(line, "Cc:"):
			data.Cc = line[3:]
		case strings.Contains(line, "Mime-Version:"):
			data.Mime_Version = line[13:]
		case strings.Contains(line, "Content-Type:"):
			data.Content_Type = line[13:]
		case strings.Contains(line, "Content-Transfer-Encoding:"):
			data.Content_Transfer_Encoding = line[26:]
		case strings.Contains(line, "X-cc:"):
			data.X_cc = line[5:]
		case strings.Contains(line, "X-bcc:"):
			data.X_bcc = line[6:]
		case strings.Contains(line, "X-Folder:"):
			data.X_Folder = line[9:]
		case strings.Contains(line, "X-Origin:"):
			data.X_Origin = line[9:]
		case strings.Contains(line, "X-FileName:"):
			data.X_FileName = line[11:]
		default:
			data.Body += line
		}
	}
	return data
}

func IndexData(data models.Email, jsonForBulk *models.JsonBulk) {

	jsonForBulk.Records = append(jsonForBulk.Records, data)

}

func IndexDataBulk(data models.JsonBulk) {

	index := "emails"

	data.Index = index

	body, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", zinc.Zinc_url, bytes.NewBuffer(body))
	fmt.Println(zinc.Zinc_url)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+zinc.Bas64encoded_creds)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func JSONfinal(jsonData models.JsonFinal) {
	file, err := os.Create("jSonFinal.json")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "   ")
	err = enc.Encode(map[string]models.JsonFinal{"Enron-email": jsonData})
	if err != nil {
		log.Fatalf("failed encoding JSON: %s", err)
	}

	fmt.Println("JSON File successfully created")
}
