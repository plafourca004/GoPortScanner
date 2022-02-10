//Package main provides the parsng functions and the scanning functions
package main

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	MAX_GO_ROUTINES = 1048475
)

var (
	ctx = context.Background() //Useful for semaphore
)

//Function maxTCPConnection.
//Function is used to determine the max number of TCP connection available to scan.
//Function is used so that we don't try to scan more TCP connection than available
func maxTCPConnection() int {
	output, err := exec.Command("/bin/sh", "-c", "ulimit -n").Output() //Equivalent to doing "ulimit -n" on a shell
	if err != nil {
		panic(err)
	}

	maxTCP, _ := strconv.Atoi(strings.TrimSpace(string(output))) //Converting it to a number

	return maxTCP
}

//Function scan.
func scan(ip string, port int, wg *sync.WaitGroup, countSemaphore *semaphore.Weighted) {
	defer countSemaphore.Release(1) //V()
	defer wg.Done()

	var result string
	fullIp := fmt.Sprintf("%s:%d", ip, port)                         //formatting "ip:host"
	connection, err := net.DialTimeout("tcp", fullIp, 3*time.Second) //Connect tcp to the ip

	if err != nil {
		result = "closed"
	} else {
		result = "open"
		defer connection.Close()
	}
	fmt.Printf("%s:%d %s\n", ip, port, result)

}

//Function main.
//Reads the different parameters from the command line.
//Can read the port as 20,21,22,23,24,25 or 20-25.
//Can read the ip adresses as 198.168.1.30 198.160.1.31 or for a whole network 198.168.1.30/30.
func main() {

	countSemaphore := semaphore.NewWeighted(int64(math.Min(MAX_GO_ROUTINES, float64(maxTCPConnection())))) //Creating weighted semaphore

	args := os.Args[1:]
	ports := args[1]
	adresses := args[2:]

	portsList := getPortsList(ports)
	ipList := getIpList(adresses)

	//fmt.Printf("args: %v\n", args)
	//fmt.Printf("ports: %v\n", portsList)
	//fmt.Printf("adresses ip: %v\n", ipList)

	//Scanning
	var wg sync.WaitGroup
	for _, ip := range ipList {
		for _, port := range portsList {
			wg.Add(1)
			countSemaphore.Acquire(ctx, 1) //P()
			go scan(ip, port, &wg, countSemaphore)
		}
	}
	//wg.Add(1)
	//go scan(ipList[0], portsList[0], 0, 0, &wg)
	wg.Wait()

}
