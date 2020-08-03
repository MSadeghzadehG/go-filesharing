package file

import (
	"net"
	"fmt"
	"bufio"
	"time"
	"bytes"
	"io"
	"strings"
	"os"
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func checkFileServer(port int, ip string, nodes map[string]int, directory string) {
	p := make([]byte, 1024)
	protocol := "udp"
	addr := net.UDPAddr{
        Port: port,
        IP: net.ParseIP(ip),
    }
    //Create the connection
	udpConn, err := net.ListenUDP(protocol, &addr)
	if err != nil {
		fmt.Println(err)
	}
	for {
        _, remoteaddr, err := udpConn.ReadFromUDP(p)
		fileName := string(bytes.Trim(p, "\x00")) // to remove null char at the end of bytes
		// fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		var message string
		if fileExists(directory + fileName) {
			message = "I HAVE"
	 	} else {
			message = "I DONT HAVE"
		}
		if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
		_,err = udpConn.WriteToUDP([]byte(message), remoteaddr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
    }
}


func GetFileServer(port int, ip string, directory string) {
	p := make([]byte, 1024)
	protocol := "tcp"
	addr := net.TCPAddr{
        Port: port,
        IP: net.ParseIP(ip),
    }
    //Create the connection
	listener, err := net.ListenTCP(protocol, &addr)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		listener.Close()
		fmt.Println("Listener closed")
	}()
	for {
		conn, err := listener.Accept()
		_, err = conn.Read(p)
		// fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			println("get Filename failed:", err.Error())
			continue
		}
		fileName := string(bytes.Trim(p, "\x00")) // to remove null char at the end of bytes
		file, err := os.Open(strings.TrimSpace(fileName)) // For read access.
		if err != nil {
			println("load File failed:", err.Error())
			continue
		}
		defer file.Close() // make sure to close the file even if we panic.
		n, err := io.Copy(conn, file)
		if err != nil {
			println("send File failed:", err.Error())
			continue
		}
		fmt.Println(n, "bytes sent")
    }
}

func GetFileClient(port int, node_ip string, fileName string, directory string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d",node_ip, port))
    if err != nil {
		fmt.Printf("Some error %v", err)
        return
    }
	defer conn.Close()
	fmt.Fprintf(conn, fileName)
	destination, err := os.Create(directory + fileName)
	if err != nil {
		fmt.Printf("error on create file %v", err)
		return
	}
	defer destination.Close()
	_, err = io.Copy(destination, conn)
	if err != nil {
		fmt.Printf("error on save file %v", err)
		return
	}
}

func checkFileFromNode(port int, ip string, fileName string) (hasFile bool, responseTime int) {
	p :=  make([]byte, 1024)
	
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d",ip, port))
    if err != nil {
		fmt.Printf("Some error %v", err)
        return
    }
	defer conn.Close()
	fmt.Fprintf(conn, fileName)
	start := time.Now()
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Printf("Discovery client error %v\n", err)
	}
	if fmt.Sprintf("%s", p) == "I HAVE" {
		t := time.Now()
		return true, int(t.Sub(start).Nanoseconds())
	} else {
		return false, -1
	}
}

func checkFile(port int, ip string, fileName string, nodes map[string]int) (bestNode string ) {
	minRpIP, minRp := "", -1
	for node_ip, _ := range(nodes) {
		if hasFile, rp := checkFileFromNode(port, ip, fileName); hasFile {
			nodes[node_ip] = rp
			if (minRp == -1 || rp < minRp) {
				minRp = rp
				minRpIP = node_ip
			}
		}
	}
	return minRpIP
}

func StartService(getFileStruct GetFile, nodes map[string]int) {
	go checkFileServer(getFileStruct.CheckPort, getFileStruct.Ip, nodes, getFileStruct.Directory)
	go GetFileServer(getFileStruct.GetPort, getFileStruct.Ip, getFileStruct.Directory)
}

func GetFileByName(getFileStruct GetFile, fileName string, nodes map[string]int) {
	node_ip := checkFile(getFileStruct.CheckPort, getFileStruct.Ip, fileName, nodes)
	GetFileClient(getFileStruct.GetPort, node_ip, fileName, getFileStruct.Directory)
}