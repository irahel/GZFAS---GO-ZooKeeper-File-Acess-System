package main

import(
	"fmt"
	"strconv"
	"time"
	"net"
	"bufio"
	"strings"
	"../messages"
	"../zk"	
)

type Client struct {
	ID 		string
	STATUS 	string	
	CONN 	net.Conn
}

var(
	//Connection to clients
	SERVER_MAIN_IP = "127.0.0.1"
	SERVER_MAIN_PORT = "1234"

	//Connection to ZooKeeper
	SERVER_ZK_IP string = "127.0.0.1"
	SERVER_ZK_PORT = "2181"

	SERVER_ID int64 = 0	
	DATA_PATH string = "/GZFAS/data"
	CONF_PATH string = "/GZFAS/conf"
	MAX_CONNECTED_CLIENTS int = 10
	MAX_READ_CLIENTS int = 10

	svListener *net.TCPListener
	zkServerConn chan *zk.Conn
	ConnectedClients chan []Client


)

					
func Connect_zk(){

	ConnectedClients = make(chan []Client, 1)
	zkServerConn = make(chan *zk.Conn, 1)
	ConnectedClients <- []Client {}
	//Connect to zk server
	var err error

	

	zkServerConnChan, _, err := zk.Connect([]string{SERVER_ZK_IP}, time.Second)
	//err check
	if err != nil{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)		
	}


	//Architecture check
	check, _ , err := zkServerConnChan.Exists("/GZFAS/server")
	if !check{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_ARCHITECTURE_ERROR)
		panic(err)		
	}else if err != nil {
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)		
	}
	check, _ , err = zkServerConnChan.Exists("/GZFAS/clients")
	if !check{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_ARCHITECTURE_ERROR)
		panic(err)		
	}else if err != nil {
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)		
	}
	check, _ , err = zkServerConnChan.Exists("/GZFAS/data")
	if !check{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_ARCHITECTURE_ERROR)
		panic(err)		
	}else if err != nil {
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)		
	}
	check, _ , err = zkServerConnChan.Exists("/GZFAS/conf")
	if !check{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_ARCHITECTURE_ERROR)
		panic(err)		
	}else if err != nil {
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)		
	}

	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_ARCHITECTURE_OK)

	//Other server check
	serverConnectionStatus, _, err := zkServerConnChan.Get("/GZFAS/server")
	if string(serverConnectionStatus) == "ON"{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_SERVER_ALREADY)
		panic(err)		
	}else if err != nil{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_SERVER_ERROR)
		panic(err)
	}
	
	//Setting default values

	//Online
	_, statt, _ := zkServerConnChan.Get("/GZFAS/server")
	_, err = zkServerConnChan.Set("/GZFAS/server", []byte("ON"), statt.Version)	
	if err != nil{
		fmt.Printf("error on 1\n")
		fmt.Printf("%t\n", err)
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//Getting ID
	SERVER_ID = zkServerConnChan.SessionID()

	//ID
	_, statt, _ = zkServerConnChan.Get("/GZFAS/server/ID")
	_, err = zkServerConnChan.Set("/GZFAS/server/ID", []byte(strconv.FormatInt(SERVER_ID, 10)), statt.Version)
	if err != nil{
		fmt.Printf("error on 2")
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//IP
	_, statt, _ = zkServerConnChan.Get("/GZFAS/server/IP")
	_, err = zkServerConnChan.Set("/GZFAS/server/IP", []byte(SERVER_ZK_IP), statt.Version)		
	if err != nil{
		fmt.Printf("error on 2")
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//PORT
	_, statt, _ = zkServerConnChan.Get("/GZFAS/server/PORT")
	_, err = zkServerConnChan.Set("/GZFAS/server/PORT", []byte(SERVER_ZK_PORT), statt.Version)	
	if err != nil{
		fmt.Printf("error on 3")
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}
	
	//Data path
	_, statt, _ = zkServerConnChan.Get("/GZFAS/server/DATA")
	_, err = zkServerConnChan.Set("/GZFAS/server/DATA", []byte(DATA_PATH), statt.Version)	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//Conf path
	_, statt, _ = zkServerConnChan.Get("/GZFAS/server/CONF")
	_, err = zkServerConnChan.Set("/GZFAS/server/CONF", []byte(CONF_PATH), statt.Version)	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//Configuration maximum number of connections
	_, statt, _ = zkServerConnChan.Get(CONF_PATH+"/MAX_CONN")
	_, err = zkServerConnChan.Set(CONF_PATH+"/MAX_CONN", []byte(strconv.Itoa(MAX_CONNECTED_CLIENTS)), statt.Version)
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//Configuration maximum number of reading clients
	_, statt, _ = zkServerConnChan.Get(CONF_PATH+"/MAX_READ")
	_, err = zkServerConnChan.Set(CONF_PATH+"/MAX_READ", []byte(strconv.Itoa(MAX_READ_CLIENTS)), statt.Version)
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	//Clients init
	_, statt, _ = zkServerConnChan.Get("/GZFAS/clients")
	_, err = zkServerConnChan.Set("/GZFAS/clients", []byte("0"), statt.Version)
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	zkServerConn <- zkServerConnChan

	//Finally
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CONNECTED)
}

func Disconnect_zk(){

	zkServerConnChan := <- zkServerConn

	_, statt, _ := zkServerConnChan.Get("/GZFAS/server")
	_, err := zkServerConnChan.Set("/GZFAS/server", []byte("OFF"), statt.Version)	
	if err != nil{		
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_SERVER_DEFAULT_VALUES)
		return
	}

	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DISCONNECTED)
	zkServerConnChan.Close()
}

func Connect_cli(ID_SEND string, CONN_SEND net.Conn){	
	
	zkServerConnChan := <- zkServerConn

	//DEBUG
	fmt.Printf("Connect specific client %s %s\n", ID_SEND, CONN_SEND)

	
	cliChan  := <- ConnectedClients
	//DEBUG
	fmt.Printf("len before append %s \n", string(len(cliChan)))
	cliChan = append(cliChan, Client{ID: ID_SEND, STATUS: "FREE", CONN: CONN_SEND})
	//DEBUG
	fmt.Printf("len after append %s \n", string(len(cliChan)))
	ConnectedClients <- cliChan

	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CONNECTED_CLI)

	data, statt, err := zkServerConnChan.Get("/GZFAS/clients")
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
		return
	}	
	atualConnected, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_STR_I_ERROR)
		return
	}	
	//DEBUG
	fmt.Printf("Connected atual %s \n", strconv.FormatInt(atualConnected, 10))
	_, err = zkServerConnChan.Set("/GZFAS/clients", []byte(strconv.FormatInt(atualConnected +1, 10)), statt.Version)	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_UPDATE_FAILED)
		return
	}	
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_UPDATE_SUCESS)

	ID_SEND = strings.Replace(ID_SEND, "\n", "", -1)
	ID_SEND = strings.Replace(ID_SEND, " ", "", -1)
	path := "/GZFAS/clients/cli_" + ID_SEND
	
	//DEBUG
	fmt.Printf("PATH = %s\n", path)
	_, err = zkServerConnChan.Create(path , []byte(" "), 0 , zk.WorldACL(zk.PermAll))		
	if err != nil{
		//DEBUG
		fmt.Printf("error on 1\n")
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_FAILED)
		panic(err)
	}
	//DEBUG
	fmt.Printf("Create client/x \n")	

	path = path + "/ID"
	//DEBUG
	fmt.Printf("PATH = %s\n", path)
	_, err = zkServerConnChan.Create(path , []byte(ID_SEND), 0 , zk.WorldACL(zk.PermAll))	
	if err != nil{
		//DEBUG
		fmt.Printf("error on 2\n")
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_FAILED)
		panic(err)
	}
	//DEBUG
	fmt.Printf("Create client/x/ID \n")	

	path = strings.Replace(path, "/ID", "/STATUS", -1)
	//DEBUG
	fmt.Printf("PATH = %s\n", path)
	_, err = zkServerConnChan.Create(path , []byte("FREE"), 0 , zk.WorldACL(zk.PermAll))	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_FAILED)
		panic(err)
	}
	//DEBUG
	fmt.Printf("Create client/x/STATUS \n")	

	path = strings.Replace(path, "/STATUS", "/ACTIONS", -1)
	//DEBUG
	fmt.Printf("PATH = %s\n", path)
	_, err = zkServerConnChan.Create(path , []byte("0"), 0 , zk.WorldACL(zk.PermAll))	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_FAILED)
		panic(err)
	}
		
	//DEBUG
	fmt.Printf("Create client/x/ACTIONS \n")	

	zkServerConn <- zkServerConnChan

	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_SUCESS)	
}

func Disconnect_cli(ID_SEND string){
	
	//DEBUG
	fmt.Printf("QUIT INIT \n")
	
	var founded = false	
	
	cliChan  := <- ConnectedClients
	
	cliNew := []Client{}

	//DEBUG
	fmt.Printf("FIND item \n")

	for _, item := range cliChan {
        if item.ID != ID_SEND{
			//DEBUG
			fmt.Printf("NOT item = %s\n", item.ID)
			cliNew = append(cliNew, item)
		}else{
			//DEBUG
			fmt.Printf("Item founded \n")
			fmt.Printf("NOT item = %s\n", item.ID)
			founded = true
		}
    }

	
	if founded{
		
		zkServerConnChan := <- zkServerConn

		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DISCONNECTED_CLI)		


		data, statt, err := zkServerConnChan.Get("/GZFAS/clients")
		if err != nil{
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
			panic(err)			
		}	
		atualConnected, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil{
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_STR_I_ERROR)
			panic(err)
		}	
		_, err = zkServerConnChan.Set("/GZFAS/clients", []byte(strconv.FormatInt(atualConnected - 1, 10)), statt.Version)	
		if err != nil{
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_UPDATE_FAILED)
			panic(err)
		}	
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_UPDATE_SUCESS)
				
		//Delete Znodes

		ID_SEND = strings.Replace(ID_SEND, "\n", "", -1)
		ID_SEND = strings.Replace(ID_SEND, " ", "", -1)
		path := "/GZFAS/clients/cli_" + ID_SEND
		path = path + "/ID"

		err = zkServerConnChan.Delete(path, -1)	
		if err != nil{			
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_FAILED)
			panic(err)			
		}
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_SUCESS)


		
		//DEBUG
		fmt.Printf("Delete client/x/ID \n")	
		
		path = strings.Replace(path, "/ID", "/STATUS", -1)
		//DEBUG
		err = zkServerConnChan.Delete(path, -1)	
		if err != nil{			
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_FAILED)
			panic(err)			
		}
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_SUCESS)
		//DEBUG
		fmt.Printf("Delete client/x/STATUS \n")	
		
		
		path = strings.Replace(path, "/STATUS", "/ACTIONS", -1)
		childrens, _ , err := zkServerConnChan.Children(path)	
		if err != nil{
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_GET_CHILDRENS_ERROR)
			panic(err)			
		}
		
		if len(childrens) != 0{
			//DEBUG
			fmt.Printf("FILHOOOOOOOOOOOOOS = %s\n", strconv.Itoa(len(childrens)))
			for _, item := range childrens {	
				err = zkServerConnChan.Delete(path +"/"+item, -1)	
				if err != nil{
					fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_GET_CHILDRENS_ERROR)
					panic(err)			
				}
				fmt.Printf("FILHO DELETED\n")

			}

		}
		
		//DEBUG
		err = zkServerConnChan.Delete(path, -1)
		if err != nil{			
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_FAILED)
			panic(err)			
		}
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_SUCESS)
		//DEBUG
		fmt.Printf("Delete client/x/ACTIONS \n")	
																

		path = "/GZFAS/clients/cli_" + ID_SEND
		
		err = zkServerConnChan.Delete(path, -1)	
		if err != nil{			
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_FAILED)
			panic(err)			
		}
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_SUCESS)


		ConnectedClients <- cliNew
		zkServerConn <-  zkServerConnChan
	}else{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CLIENT_NOT_EXISTS)		
		ConnectedClients <- cliChan
	}

	
	
}

func StartLocal(){

	//TCP connection
	//Address
	tcpAddr, err := net.ResolveTCPAddr("tcp", ""+ SERVER_MAIN_IP + ":" + SERVER_MAIN_PORT)
	if err != nil {
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_TCP_A_ERROR)
		panic(err)
	}

	//Listener
	svListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_TCP_C_ERROR)
		panic(err)
	}

	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_TCP_C_SUCCESS)

}

func CloseLocal(){
	svListener.Close()
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_TCP_DC_SUCCESS)
}

func HandleCliConnection(conn net.Conn){
	
	

	//DEBUG
	fmt.Printf("Enter In Handle COnnection\n")

	//Possible commands
	//CONNECT ID
	//READ ID ArchivePath
	//WRITE ID ArchivePath
	//QUIT ID

	for {
		
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil{
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_TCP_B_ERROR)
			panic(err)
		}
		
		//DEBUG
		fmt.Printf("Received message = %s\n", message)
		
		commands := strings.Split(message, " ")

		if strings.HasPrefix(message, "CONNECT"){		
			//DEBUG				
			fmt.Printf("connecting %s \n", commands[1])
			Connect_cli(commands[1], conn)
		}else if strings.HasPrefix(message, "READ"){
			//Check if the archive exists
			zkServerConnChan := <- zkServerConn
			check, _ , err := zkServerConnChan.Exists(DATA_PATH+"/"+commands[2])	
			if err != nil{				
				fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
				continue			
			}
			if check{
				fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_TRUE)
			}else{
				fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
				continue
			}

			//Get Childrens
			childrens, _ , err := zkServerConnChan.Children(DATA_PATH)
			if err != nil{				
				fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
				panic(err)
			}
			for _, element := range childrens {		
				if element == commands[2]{
					//Check status
					data, _ , err := zkServerConnChan.Get(DATA_PATH+"/"+commands[2]+"/STATUS")	
					if err != nil{						
						fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
						panic(err)
					}
					
					if string(data) == "FREE"{
						//Dar aceso ao cliente
					}
				}
			}
			zkServerConn <- zkServerConnChan
		}else if strings.HasPrefix(message, "WRITE"){
		}else if strings.HasPrefix(message, "QUIT"){
			Disconnect_cli(commands[1])
			break
		}

		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
	
}

func main(){		
	Connect_zk()
	defer Disconnect_zk()

	StartLocal()
	defer CloseLocal()
	
	//Accept connections	
	
	for {

        conn, err := svListener.Accept()
        if err != nil {
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONN_A_ERROR)
			panic(err)
        }
		
		go HandleCliConnection (conn)
	}
	
	
}

//TO DO ANNOTATION
