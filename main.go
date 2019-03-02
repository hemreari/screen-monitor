package main

import(
	"os/exec"
	"fmt"
	"strings"
	"log"
	"encoding/json"
	"os"
	"path/filepath"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLStruct struct {
	Storage StorageConfig `json:"db"`
}

type StorageConfig struct {
	Driver string `json:"driver"`
	Name string `json:"name"`
}

var config SQLStruct

func readConfig(cfg *SQLStruct, configFileName string) {
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}

}
func main() {
	readConfig(&config, "config.json")

	db, err := sql.Open(config.Storage.Driver, config.Storage.Name)
	if err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command("screen", "-ls").Output()

	//xregexString, _ := regexp.Compile("p([a-z]+)ch")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(out)
	screenOut := string(out)
	//screenOut = strings.TrimSpace(screenOut)
	screenOutArr := strings.Split(screenOut, "\n")
	screenOut = screenOutArr[1]
	screenOut = strings.TrimSpace(screenOut)
	screenOutArr = strings.Split(screenOut, "\t")
	screenOut = screenOutArr[0]
	fmt.Println(screenOut)

	screenOutArr = strings.Split(screenOut, ".")
	screenPID := screenOutArr[0]
	screenName := screenOutArr[1]

	fmt.Println("screenPID:" + screenPID)
	fmt.Println("screenName:" + screenName)

	stmt, err := db.Prepare("INSERT INTO ScreenInfo(PID, screen_name, up) VALUES(?, ?, ?)")
	if err != nil {
		log.Println(err)
	}

	res, err := stmt.Exec(screenPID, screenName, 1)
	if err != nil {
		log.Println(err)
	}

	id, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Println(id)
}
