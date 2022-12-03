### Librarygo

The purpose of this CLI software is to search for authors of any book. Additionaly all works of authors are displayed in YAML.
Software is written in Go language. Data are gathered from openlibrary.org API.

### Prerequisites

1. go 1.19
2. gopkg.in/yaml.v3

* to run code in your IDE type 
  ```
  go run main.go
  ```
* or build executable file

  ```
  go build
  ```

### Usage example

1. Run code ( see prequisities above)
2. Type ```-help``` to show additional info
3. Type ```-limit 10``` to limit results for every author to 10. Any positive number is accepted in range of int datatype.
4. Type for example ```-search the lord of the rings``` to display additional works for every author contributing to searched book. 


---
<p  align="center" >Welcome screen</p>
<p align="center">
  <img src="https://github.com/serstuk93/librarygo/blob/master/preview/scr1.png" alt="screenshot" />
</p>


---
<p  align="center">Results</p>
<p align="center">
  <img src="https://github.com/serstuk93/librarygo/blob/master/preview/scr2.png" alt="screenshot" />
</p>

