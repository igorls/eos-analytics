package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"os"
	"io/ioutil"
)

type GetInfoResponse struct {
	ServerVersion            string `json:"server_version"`
	HeadBlockNum             int    `json:"head_block_num"`
	HeadBlockProducer        string `json:"head_block_producer"`
	HeadBlockTime            string `json:"head_block_time"`
	HeadBlockID              string `json:"head_block_id"`
	LastIrreversibleBlockNum int    `json:"last_irreversible_block_num"`
}

type NodeList struct {
	Nodes []Node `json:"blockProducerList"`
}

type Node struct {
	Name        string `json:"np_name"`
	Org         string `json:"organisation"`
	Location    string `json:"location"`
	NodeAddress string `json:"node_addr"`
	PortHTTP    string `json:"port_http"`
	PortSSL     string `json:"port_ssl"`
	PortP2P     string `json:"port_p2p"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func getJson(url string, target interface{}) (int64, error) {
	t := int64(0)
	start := time.Now()
	fmt.Println("Fetching data from", url)
	r, err := httpClient.Get(url)
	if err != nil {
		return t, err
	}
	defer r.Body.Close()
	elapsed := time.Since(start).Nanoseconds()
	t = elapsed
	return t, json.NewDecoder(r.Body).Decode(target)
}

func findPublicIP(server string) string {
	resp, err := httpClient.Get(server)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return bodyString
	} else {
		return "none"
	}
}

func main() {
	fmt.Println("EOS-Analytics")
	filepath := "testnets/jungle3.json"
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + filepath)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var nodeList NodeList
	json.Unmarshal(byteValue, &nodeList)

	// check node's external ip address
	externalIP := findPublicIP("http://myexternalip.com/raw")
	fmt.Println("Firing requests from:", externalIP)
	externalIP2 := findPublicIP("https://ipv4.icanhazip.com/")
	fmt.Println("Firing requests from:", externalIP2)

	fmt.Println(len(nodeList.Nodes), "nodes on list")

	for _, element := range nodeList.Nodes {
		data := new(GetInfoResponse)
		nodeURL := "http://" + element.NodeAddress + ":" + element.PortHTTP + "/v1/chain/get_info"
		responseTime, err := getJson(nodeURL, &data)
		if err != nil {
			fmt.Println("Server is down")
		} else {
			fmt.Println("Block:", data.HeadBlockNum)
			fmt.Println("Producer:", data.HeadBlockProducer)
			fmt.Println("Latency:", float64(responseTime)/1000000, "ms")
		}
		fmt.Println("-------")
	}
}
