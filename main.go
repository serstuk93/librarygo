package main

/*
1. Create simple cli application for finding all works from authors of specific book
2. Application has to find all authors for book and it will print list of all their works
3. Create list of works for each author (name, revision)
4. Print result to stdout in yaml format sorted by author name, count of revision (asc, desc as argument).
Names of authors have to be part of output.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"encoding/json"
)

/*
Type of struct needed for parsing useful data. "Key" value is ID of work(book).
AuthorsKeys gets slice of ID values for later searching of those authors again.
AuthorsNames are literal names in string format.
Index [0] in AuthorsKeys and AuthorsNames shows data forthe same person.
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
		"To search authors of the book type '-search <name of the book>'"}
	for _, element := range introCli {
		fmt.Printf("%s  \n", element)
	}
	response, err := http.Get("https://openlibrary.org/search.json?q=has_fulltext:true%20AND%20title:the%20lord%20of%20the%20rings&fields=key,author_key,author_name,availability&limit=1")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

	var worksSearch Docs
	if err := json.Unmarshal(responseData, &worksSearch); err != nil {
		log.Fatal(err)
	}

	fmt.Println("WORK ID", worksSearch.Results[0].Key)
	for i, element := range worksSearch.Results[0].AuthorsKeys {
		fmt.Print("namekey ", element)
		fmt.Println(" name ", worksSearch.Results[0].AuthorsNames[i])
	}

}
