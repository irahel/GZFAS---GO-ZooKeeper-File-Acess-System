package main

import(
	"fmt"	
	"time"
	"net"
	"bufio"
	"os"
	"../messages"
	"../zk"		
)

var(
	CON_IP_ZK string = "127.0.0.1"
	CON_PORT_ZK string = "2181"

	CON_IP_SERVER string = "127.0.0.1"
	CON_PORT_SERVER string = "1234"

	CON_TYPE = "tcp"
	zkCON *zk.Conn
)

func Start() *zk.Conn{
	var err error
	zkCON, _, err = zk.Connect([]string{CON_IP_ZK}, time.Second)
	//err check
	if err != nil{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)		
	}
	fmt.Printf("%s %s\n", messages.CLI_PREFIX, messages.MESSAGE_CONNECTED)
		

	return zkCON
}

func main(){
	fmt.Printf("CLient main  \n")	
		
	conn, err := net.Dial(CON_TYPE, CON_IP_SERVER+":"+CON_PORT_SERVER)
	if err != nil{
		fmt.Printf("%s %s\n", messages.CLI_PREFIX, messages.MESSAGE_TCP_C_ERROR)
		panic(err)
	}

	for { 
					
		reader := bufio.NewReader(os.Stdin)
	   	fmt.Print("Text to send: ")
	   	text, _ := reader.ReadString('\n')
	   	// send to socket
	   	fmt.Fprintf(conn, text + "\n")
	   	// listen for reply
	   	message, _ := bufio.NewReader(conn).ReadString('\n')
	   	fmt.Print("Message from server: "+message)
	}	
}