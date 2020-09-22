package main

import (
	"flag"
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"github.com/tarm/serial"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)
const (
	gwState = rpio.Pin(17)
	lan9514 = rpio.Pin(20)
	lan9512 = rpio.Pin(21)
	CTL_1 = rpio.Pin(36)
	CTL_2 = rpio.Pin(37)
	rs485A = rpio.Pin(22)
	wireless = rpio.Pin(25)
)
func LedTest(wg *sync.WaitGroup){
	if err := rpio.Open(); err!= nil{
		fmt.Println(err)
		os.Exit(1)
	}
	gwState.Output()
	lan9514.Output()
	lan9512.Output()
	rs485A.Output()
	CTL_1.Output()
	CTL_2.Output()
	wireless.Output()

	for{    //GPIO LED Toggle
		gwState.Toggle()
		lan9514.Toggle()
		lan9512.Toggle()
		rs485A.Toggle()
		CTL_1.Toggle()
		CTL_2.Toggle()
		wireless.Toggle()
		time.Sleep(time.Second)
	}
}
func SerialTest(wg *sync.WaitGroup){
	rs485_A:= &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, StopBits: 1, Parity: 'N'}
	rs485_IO := &serial.Config{Name: "/dev/ttyUSB1", Baud: 115200, StopBits: 1, Parity: 'N'}
	ble := &serial.Config{Name: "/dev/ttyUSB2", Baud: 115200, StopBits: 1, Parity: 'N'}

	a, err := serial.OpenPort(rs485_A)
	if err != nil{
		log.Fatal(err)
	}
	b, err := serial.OpenPort(rs485_IO)
	if err != nil{
		log.Fatal(err)
	}
	c, err := serial.OpenPort(ble)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Serial port Open")
	for{
		n, err := a.Write([]byte("test\n"))
		if err != nil {
			log.Fatal(n)
		}
		m, err := b.Write([]byte("test\n"))
		if err != nil {
			log.Fatal(m)
		}
		l, err := c.Write([]byte("test\n"))
		if err != nil {
			log.Fatal(l)
		}
	}
}
func EthernetTest(wg *sync.WaitGroup){
	port := flag.Int("port", 3334, "Port to accept connections on.")
	flag.Parse()

	l, err := net.Listen("tcp",":"+strconv.Itoa(*port))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to connections at on port", strconv.Itoa(*port))
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	log.Println("Accepted new connection.")

	for{
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]
		log.Println("Read new data from connection", data)
		conn.Write(data)

	}
}
func main() {
	//gpio pin setting
	//Ethernet TCP setting
	//Board Test
	//Led start
	//Rs485 start
	//Ethernet
	var wg sync.WaitGroup

	log.Println("start led toggle")
	wg.Add(3)
	go LedTest(&wg)

	log.Println("start Serial server")
	//wg.Add(2)
	go SerialTest(&wg)

	log.Println("start tcp server")
	//wg.Add(3)
	go EthernetTest(&wg)


	wg.Wait()
}
