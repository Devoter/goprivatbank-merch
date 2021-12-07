package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	XMLName  xml.Name `xml:"request"`
	Version  string   `xml:"version,attr"`
	Merchant Merchant
	Data     Data
}

type Merchant struct {
	XMLName   xml.Name `xml:"merchant"`
	ID        string   `xml:"id"`
	Signature string   `xml:"signature"`
}

type Data struct {
	XMLName xml.Name `xml:"data"`
	Oper    string   `xml:"oper"`
	Wait    int      `xml:"wait"`
	Test    int      `xml:"test"`
	Payment Payment  `xml:"payment"`
}

type Payment struct {
	XMLName xml.Name `xml:"payment"`
	ID      string   `xml:"id,attr"`
	Props   []Prop
}

type Prop struct {
	XMLName xml.Name `xml:"prop"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

var Version = "undefined"

func main() {
	url := flag.String("url", "", "target URL")
	passphrase := flag.String("passphrase", "", "private key")
	version := flag.Bool("version", false, "print the application version")
	dryRun := flag.Bool("dry-run", false, "just prepare a body")
	help := flag.Bool("help", false, "show this text")

	flag.Parse()

	if *version {
		fmt.Printf("version=%s\n", Version)
		return
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	data := Data{
		Oper: "cnt",
		Wait: 0,
		Test: 0,
		Payment: Payment{
			Props: []Prop{
				{Name: "sd", Value: "11.08.2013"},
				{Name: "ed", Value: "11.09.2013"},
				{Name: "card", Value: "5168742060221193"},
			},
		},
	}

	encodedData, err := xml.Marshal(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not marshal a data error=%v\n", err)
		os.Exit(1)
	}

	encodedData = append(encodedData[6:len(encodedData)-7], []byte(*passphrase)...)

	fmt.Printf("Data inner XML:\n%s\n", encodedData)

	md5sum := md5.Sum(encodedData)
	signature := sha1.Sum([]byte(hex.EncodeToString(md5sum[:])))

	reqContent := Request{
		Version: "1.0",
		Merchant: Merchant{
			ID:        "75482",
			Signature: hex.EncodeToString(signature[:]),
		},
		Data: data,
	}

	encoded, err := xml.Marshal(&reqContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not marshal a request body error=%v\n", err)
		os.Exit(1)
	}

	encoded = append([]byte(xml.Header), encoded...)

	fmt.Printf("\nRequest body:\n%s\n", encoded)

	if *dryRun {
		return
	}

	buf := bytes.NewBuffer(encoded)

	_, err = http.Post(*url, "application/xml", buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP request failed error=%v\n", err)
		os.Exit(3)
	}
}
