package temp

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type Transaction struct {
	From  string
	To    string
	Value int
}

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}

func DecodeToPerson(s []byte) Transaction {

	p := Transaction{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func Test() {
	fmt.Println("test222222")
}
func WriteCsv() {

	record := []string{"first_name"}

	// x := [5]int{10, 20, 30, 40, 50}

	f, err := os.Create("users.csv")
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}

func UpdateCsv(filename string, data string) {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open example.txt")
		return
	}

	defer file.Close()

	_, err2 := file.WriteString(data)

	if err2 != nil {
		fmt.Println("Could not write text to example.txt")

	} else {
		fmt.Println("Operation successful! Text has been appended")
	}
}

func ReadCsv(filename string) [][]string {

	f, err := os.Open(filename)

	if err != nil {
		fmt.Println("eeeeeerrrrrrrrrrrrrrrrrrrooooooor")
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line

	records, err := r.ReadAll()

	if err != nil {
		fmt.Println("eeeeeerrrrrrrrrrrrrrrrrrrooooooor")
	}

	return records

}
