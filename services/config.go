package services

import (
	"bufio"
	"bytes"
	"fmt"

	"encoding/base64"
	"io/ioutil"
	"path"
	"strings"

	"github.com/joho/godotenv"

	"EnronEmailApi/models"
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
		} else {
			listFolders = append(listFolders, f.Name())
		}
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

func parseData(dataLines *bufio.Scanner) models.Email {
	var data models.Email
	for dataLines.Scan() {
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
func Algodeaca(folderList []string, path string, wg *sync.WaitGroup) {

	defer wg.Done()

	for _, user := range folderList {

		var jsonForBulk models.JsonBulk
		jsonForBulk.Index = "email"

		fmt.Println(user)

		processDir(path+user, &jsonForBulk)

		IndexDataBulk(jsonForBulk)
		jsonForBulk.Records = []models.Email{}

	}

}

func ProcessMailFile(path string, jsonForBulk *models.JsonBulk, wg *sync.WaitGroup) {
	defer wg.Done()

	sysFile, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}
	defer sysFile.Close()

	lines := bufio.NewScanner(sysFile)

	data := parseData(lines)

	IndexData(data, jsonForBulk)

}

func IndexDataBulk(jsonForBulk models.JsonBulk) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	var (
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
	)

	auth := user + ":" + password
	bas64encoded_creds := base64.StdEncoding.EncodeToString([]byte(auth))

	jsonForBulk.Index = "emails"
	zinc_host := "http://localhost:4080"
	zinc_url := zinc_host + "/api/_bulkv2"

	body, _ := json.Marshal(jsonForBulk)

	req, err := http.NewRequest("POST", zinc_url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+bas64encoded_creds)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func IndexData(data models.Email, jsonForBulk *models.JsonBulk) {

	jsonForBulk.Records = append(jsonForBulk.Records, data)
}

func processDir(name string, jsonForBulk *models.JsonBulk) {
	d, err := os.Open(name)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer d.Close()

	files, err := d.ReadDir(-1)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, f := range files {
		if f.IsDir() {
			processDir(path.Join(name, f.Name()), jsonForBulk)
		} else {
			var wg sync.WaitGroup
			wg.Add(1)
			go ProcessMailFile(path.Join(name, f.Name()), jsonForBulk, &wg)
			wg.Wait()
		}
	}
}
