package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"os"
	"io/ioutil"
)

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"strconv"
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
	Name        string `json:"bp_name"`
	Org         string `json:"organisation"`
	Location    string `json:"location"`
	NodeAddress string `json:"node_addr"`
	PortHTTP    string `json:"port_http"`
	PortSSL     string `json:"port_ssl"`
	PortP2P     string `json:"port_p2p"`
	Coordinates string
	Responses   []float64
	AVR         float64
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

func avg(input []float64) float64 {
	var total float64 = 0
	for _, val := range input {
		total += val
	}
	return total / float64(len(input))
}

func printUsage(u *disk.UsageStat) {
	fmt.Println(u.Path + "\t" + strconv.FormatFloat(u.UsedPercent, 'f', 2, 64) + "% full.")
	fmt.Println("Total: "  + strconv.FormatUint(u.Total/1024/1024/1024, 10) + " GiB")
	fmt.Println("Free:  "  + strconv.FormatUint(u.Free /1024/1024/1024, 10) + " GiB")
	fmt.Println("Used:  "  + strconv.FormatUint(u.Used /1024/1024/1024, 10) + " GiB")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}


func main() {
	fmt.Println("EOS-Analytics")

	v, _ := mem.VirtualMemory()
	ram := float64(v.Total)/1000000000
	fmt.Printf("Total System RAM: %.2f GB\n", ram)
	fmt.Printf("UsedPercent: %.2f%%\n", v.UsedPercent)
	fmt.Printf("------- CPU ------ \n")
	c,_ := cpu.Info()
	fmt.Printf("CPU Model: %v\n",c[0].ModelName)
	fmt.Printf("CPU Clock: %v\n",c[0].Mhz)
	fmt.Printf("CPU Cores: %v\n",c[0].Cores)
	cpuUsage,_ := cpu.Percent(1,true)
	fmt.Println(cpuUsage)
	fmt.Printf("------- Storage ------ \n")

	parts, err := disk.Partitions(false)
	check(err)

	var usage []*disk.UsageStat

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		check(err)
		usage = append(usage, u)
		printUsage(u)
	}

	//filepath := "testnets/jungle3.json"
	//jsonFile, err := os.Open(filepath)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Successfully Opened " + filepath)
	//defer jsonFile.Close()
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//var nodeList NodeList
	//json.Unmarshal(byteValue, &nodeList)
	//
	//// check node's external ip address
	//externalIP := findPublicIP("http://myexternalip.com/raw")
	//fmt.Println("Firing requests from:", externalIP)
	//externalIP2 := findPublicIP("https://ipv4.icanhazip.com/")
	//fmt.Println("Firing requests from:", externalIP2)
	//
	//fmt.Println(len(nodeList.Nodes), "nodes on list")
	//
	//for index, _ := range nodeList.Nodes {
	//	nodeList.Nodes[index].Responses = []float64{}
	//}
	//
	//for i := 0; i < 5; i++ {
	//	for index, element := range nodeList.Nodes {
	//		data := new(GetInfoResponse)
	//		nodeURL := "http://" + element.NodeAddress + ":" + element.PortHTTP + "/v1/chain/get_info"
	//		responseTime, err := getJson(nodeURL, &data)
	//		if err != nil {
	//			fmt.Println("Server is down")
	//		} else {
	//			//fmt.Println("Block:", data.HeadBlockNum)
	//			//fmt.Println("Producer:", data.HeadBlockProducer)
	//			respTime := float64(responseTime) / 1000000
	//			fmt.Println("Latency:", respTime, "ms")
	//			nodeList.Nodes[index].Responses = append(element.Responses, respTime)
	//		}
	//		fmt.Println("-------")
	//	}
	//}
	//
	//for index, element := range nodeList.Nodes {
	//	fmt.Println(element.Name + " | " + element.Org)
	//	fmt.Printf("%v\n", element.Responses)
	//	if len(element.Responses) > 0 {
	//		fmt.Println("Average Response Time:", avg(element.Responses))
	//		nodeList.Nodes[index].AVR = avg(element.Responses)
	//	}
	//}
	//
	//sort.Slice(nodeList.Nodes, func(i, j int) bool {
	//	return nodeList.Nodes[i].AVR < nodeList.Nodes[j].AVR
	//})
	//
	//var output []Node
	//for _, element := range nodeList.Nodes {
	//	if len(element.Responses) > 0 {
	//		output = append(output, element)
	//	}
	//}
	//
	//fmt.Println("----------------------------")
	//fmt.Println("Fastest nodes for config.ini")
	//fmt.Println("----------------------------")
	//
	//output = output[0:6]
	//for _, element := range output {
	//	fmt.Println("p2p-peer-address = " + element.NodeAddress + ":" + element.PortP2P)
	//}

}
