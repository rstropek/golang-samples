package gabdownloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
)

const gitHubBaseURL = "https://api.github.com/repos/coding-club-linz/global-azure-bootcamp-2019/contents/"
const sessionsFolder = "sessions"
const sessionListFile = "./session-list.json"
const sessionsFile = "./sessions.json"

type Session struct {
	ID        string `json:"id"`
	Title     string `yaml:"title" json:"title"`
	Speaker   string `yaml:"speaker" json:"speaker"`
	SpeakerID string `yaml:"speaker-id" json:"speakerId"`
	Room      string `yaml:"room" json:"room"`
	Slot      int    `yaml:"slot" json:"slot"`
	Content   string `json:"content"`
	URL       string `json:"url"`
}

type Sessions []Session

type SessionReference struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	DownloadURL string `json:"download_url"`
}

type SessionReferences []SessionReference

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetAndCacheSessionReferencesFromGitHub() SessionReferences {
	var files SessionReferences
	stat, err := os.Stat(sessionListFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Local session list cache in file %s does not exist, creating and filling it...\n", sessionListFile)
			files = GetSessionReferencesFromGitHub()
			filesData, err := json.Marshal(files)
			check(err)
			ioutil.WriteFile(sessionListFile, filesData, 0644)
		} else {
			panic(err)
		}
	} else if !stat.IsDir() {
		fmt.Printf("Found local session list cache in file %s, reading it...\n", sessionListFile)
		filesData, err := ioutil.ReadFile(sessionListFile)
		check(err)
		check(json.Unmarshal(filesData, &files))
	} else {
		panic(fmt.Errorf("%s must not be an existing folder", sessionListFile))
	}

	return files
}

func GetSessionReferencesFromGitHub() SessionReferences {
	resp, err := http.Get(strings.Join([]string{gitHubBaseURL, sessionsFolder}, ""))
	check(err)

	defer resp.Body.Close()

	var files SessionReferences = make([]SessionReference, 0)
	err = json.NewDecoder(resp.Body).Decode(&files)
	check(err)

	if len(files) == 0 {
		panic(errors.New("Did not find any sessions"))
	}

	return files
}

func GetAndCacheSessionsFromGitHub(files SessionReferences) Sessions {
	var sessions []Session

	stat, err := os.Stat(sessionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Local session cache in file %s does not exist, creating and filling it...\n", sessionsFile)
			sessions = GetSessionsFromGitHub(files)
			sessionsData, err := json.Marshal(sessions)
			check(err)
			ioutil.WriteFile(sessionsFile, sessionsData, 0644)
		} else {
			panic(err)
		}
	} else if !stat.IsDir() {
		fmt.Printf("Found local session cache in file %s, reading it...\n", sessionsFile)
		sessionsData, err := ioutil.ReadFile(sessionsFile)
		check(err)
		check(json.Unmarshal(sessionsData, &sessions))
	} else {
		panic(fmt.Errorf("%s must not be an existing folder", sessionsFile))
	}

	return sessions
}

func GetSessionsFromGitHub(files SessionReferences) Sessions {
	sessions := make([]Session, len(files))
	for i, file := range files {
		resp, err := http.Get(file.DownloadURL)
		check(err)

		defer resp.Body.Close()

		mdContentBytes := make([]byte, 0)
		mdContentBytes, err = ioutil.ReadAll(resp.Body)
		check(err)

		mdContent := string(mdContentBytes)
		yamlHeaderStart := strings.Index(mdContent, "---")
		if yamlHeaderStart == (-1) {
			panic(errors.New("Did not find beginning of front matter"))
		}
		yamlheaderEnd := strings.Index(mdContent[(yamlHeaderStart+3):], "---")
		if yamlheaderEnd == (-1) {
			panic(errors.New("Did not find end of front matter"))
		}

		check(yaml.Unmarshal(mdContentBytes[(yamlHeaderStart+3):(yamlheaderEnd+3)], &sessions[i]))
		sessions[i].Content = string(blackfriday.Run(mdContentBytes[(yamlheaderEnd + 6):]))
		sessions[i].URL = file.URL
		sessions[i].ID = file.Name
	}

	return sessions
}
