package main

/*
1. Create simple cli application for finding all works from authors of specific book
2. Application has to find all authors for book and it will print list of all their works
3. Create list of works for each author (name, revision)
4. Print result to stdout in yaml format sorted by author name, count of revision (asc, desc as argument).
Names of authors have to be part of output.
*/

/*for installation of yaml package:
if you are in your default GOPATH then use "go mod init"
else use "go mod init <name>"
then use "go mod tidy"
type in terminal "go get gopkg.in/yaml.v3/" */

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
	Type of struct needed for parsing useful data. "Key" value is ID of work(book).
	"AuthorsKeys" gets slice of ID values for later searching of those authors again.
	"AuthorsNames" are literal names in string format.

Index [0] in AuthorsKeys and AuthorsNames shows data for the same person.
*/
type KeysSearch struct {
	Key          string   `json:"key"`
	AuthorsKeys  []string `json:"author_key"`
	AuthorsNames []string `json:"author_name"`
}

// API response consist of syntax { "docs":[{"key": "XYZ"}]}. Thus its necessary to assign "docs" as slice.
type Docs struct {
	Results []KeysSearch `json:"docs"`
}

// Gives bool value for checking if there is exact match of user input and books in database.
type ExactMatch struct {
	ExctMtch bool `json:"numFoundExact"`
}

func main() {
	var introCli []string = []string{"This is LibraryGo created by Radoslav Serstuk.",
		"For help type '-help'",
		"To search authors of the book type '-search <name of the book>'",
		""}
	for _, element := range introCli {
		fmt.Printf("%s  \n", element)
	}

	limit := "200"
	var input string
	startSearch, _ := regexp.Compile("^-search .*")
	helpinfo, _ := regexp.Compile("^-help.*")
	setlimit, _ := regexp.Compile(`^-limit \d`)
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Waiting for input: ")
		scanner.Scan()
		input = scanner.Text()
		if startSearch.MatchString(input) {
			input = strings.TrimSpace(input[7:])
			input = strings.ReplaceAll(input, " ", "%20")
			break
		}
		if helpinfo.MatchString(input) {
			fmt.Println("Supported commands are -limit, -search, -help.")
			fmt.Println("To reduce number of displayed works for every author write command -limit <number> i.e. -limit 10, default value is 200.")
			fmt.Println("In command -search <book> , the book value is not case sensitive.")
			fmt.Println("Example of usage: -search the lord of the rings")
			fmt.Println("Authors and book titles are ordered by ascending value, revision number is in descending order.")
			fmt.Println("Results are shown in YAML.")
		}
		if setlimit.MatchString(input) {
			limit = strings.TrimSpace(input[7:])
			fmt.Println("Setting new limit for works: ", limit)
		}

	}
	fmt.Println(input)

	//search_query := "https://openlibrary.org/search.json?q=has_fulltext:true%20AND%20title:the%20lord%20of%20the%20rings&fields=key,author_key,author_name,availability&limit=1"
	search_query := "https://openlibrary.org/search.json?q=has_fulltext:true%20AND%20title:" + input + "&fields=key,author_key,author_name,availability&limit=1"
	fmt.Println(search_query)
	response, err := http.Get(search_query)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(responseData))

	var worksSearch Docs
	if err := json.Unmarshal(responseData, &worksSearch); err != nil {
		log.Fatal(err)
	}

	/*fmt.Println("WORK ID", worksSearch.Results[0].Key)
	for i, element := range worksSearch.Results[0].AuthorsKeys {
		fmt.Print("namekey ", element)
		fmt.Println(" name ", worksSearch.Results[0].AuthorsNames[i])
	}*/

	//var worksPerAuthor []string
	//worksPerAuthor = worksbyauthors(worksSearch.Results[0])
	worksbyauthors(worksSearch.Results[0], limit)
}

type Entries struct {
	Title     string `json:"title" yaml:"title"`
	Revisions int    `json:"latest_revision" yaml:"revision"`
}

/*
	API response consist of list accessed by keyword entries. Preview of reponse structure:

{ "entries":[{"title": "XYZ", "latest_revision": 123 }], [{"title": "XYZ", "latest_revision": 123 }]}
*/
type Works struct {
	Author  string    `yaml:"author"`
	Results []Entries `json:"entries" yaml:"works"`
}

/*
	Function uses "KeySearch" struct to iterate over "AuthorKeys".

For every Authorkey there is created a list of works consisting of its names and revision count.
*/
func worksbyauthors(dt KeysSearch, limit string) {
	var http_query string
	allWorks := []Works{}

	for i, element := range dt.AuthorsKeys {
		http_query = "https://openlibrary.org/authors/" + element + "/works.json?limit=" + limit
		//	fmt.Println("http_query :", http_query)
		response, err := http.Get(http_query)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		//	fmt.Println(string(responseData))

		var wrkList Works
		if err := json.Unmarshal(responseData, &wrkList); err != nil {
			log.Fatal(err)
		}

		wrkList.Author = dt.AuthorsNames[i]

		//fmt.Println("WRKL", wrkList.Results)

		sort.SliceStable(wrkList.Results, func(i, j int) bool {
			return wrkList.Results[i].Title < wrkList.Results[j].Title
		})

		sort.SliceStable(wrkList.Results, func(i, j int) bool {
			return wrkList.Results[i].Revisions > wrkList.Results[j].Revisions
		})
		//fmt.Println("WRKLS", wrkList.Author, wrkList.Results)
		allWorks = append(allWorks, wrkList)
	}

	//fmt.Println(allWorks)

	sort.SliceStable(allWorks, func(i, j int) bool {
		return allWorks[i].Author < allWorks[j].Author
	})
	//fmt.Println("SRT", allWorks)

	yamlData, err := yaml.Marshal(&allWorks)

	if err != nil {
		fmt.Printf("Error while creating YAML. %v", err)
	}
	//fmt.Println(string(yamlData))
	os.Stdout.Write(yamlData)

	//fmt.Println("allWorks: ", allWorks)
}
