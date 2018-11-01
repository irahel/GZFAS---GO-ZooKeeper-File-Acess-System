package main

import(
	"fmt"
	"time"
	"bufio"
	"os"
	"strings"
	"./messages"
	"./zk"
)

var SERVER_IP = "127.0.0.1"
var reserved_words_commands = [...]string{"CREATE", "READ", "EXISTS", "LIST", "UPDATE", "DELETE", "DELETEC", "QUIT"}
var zkCON *zk.Conn

//Check errors on enter command
func Command_treatment(enter string) bool{		
	if strings.HasPrefix(enter, reserved_words_commands[0]) || strings.HasPrefix(enter, reserved_words_commands[4]) {
		var separateString = strings.Split(enter, " ")
		if len(separateString) > 3{
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_COMMAND_MANY_ARGUMMENTS)
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_USAGE_CREATE_COMMAND)
			return false
		}else if len(separateString) < 3{
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_COMMAND_FEW_ARGUMMENTS)
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_USAGE_CREATE_COMMAND)
			return false
		}
	}else if strings.HasPrefix(enter, reserved_words_commands[1]) || strings.HasPrefix(enter, reserved_words_commands[5]) || strings.HasPrefix(enter, reserved_words_commands[3]) || strings.HasPrefix(enter, reserved_words_commands[2]) || strings.HasPrefix(enter, reserved_words_commands[6]){
		var separateString = strings.Split(enter, " ")
		if len(separateString) > 2{
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_COMMAND_MANY_ARGUMMENTS)
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_USAGE_READ_COMMAND)
			return false
		}else if len(separateString) < 2{
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_COMMAND_FEW_ARGUMMENTS)
			fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_USAGE_READ_COMMAND)
			return false
		}
	}else if strings.HasPrefix(enter, reserved_words_commands[7]) == false{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_COMMAND_EMPTY)
		return false
	}	
	return true
}

//CRUD
//C
func Create(path string, value string){		
	_, err := zkCON.Create(path, []byte(value), 0, zk.WorldACL(zk.PermAll))	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_FAILED)
		return
	}	
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CREATE_SUCESS)	
}

//R
func Read(path string){	
	data, _ , err := zkCON.Get(path)	
	if err != nil{
		fmt.Printf("%+v", err)
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
		return
	}
	fmt.Printf("%s\n", string(data))		
}
func Exists(path string){
	check, _ , err := zkCON.Exists(path)	
	if err != nil{
		fmt.Printf("%+v", err)
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
		return
	}
	if check{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_TRUE)
	}else{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
	}
	
}
func List(path string){
	childrens, _ , err := zkCON.Children(path)
	if err != nil{
		fmt.Printf("%+v", err)
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_EXISTS_FALSE)
	}
	for _, element := range childrens {		
		fmt.Printf("\t %s\n", element)
	}
}

//U
func Update(path string, value string){
	_, err := zkCON.Set(path, []byte(value), 0)	
	if err != nil{
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_UPDATE_FAILED)
		return
	}	
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_UPDATE_SUCESS)
}

//D
func Delete(path string){
	err := zkCON.Delete(path, -1)	
	if err != nil{
		fmt.Printf("%+v", err)
		fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_FAILED)
		return
	}
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DELETE_SUCESS)
}
func Delete_c(path string){	
}

func main(){
	var err error
	zkCON, _, err = zk.Connect([]string{SERVER_IP}, time.Second)
	
	//err check
	if err != nil{
		fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_CONNECT_ERROR)
		panic(err)		
	}
	fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_CONNECTED)
	fmt.Printf("%s\n", messages.LOG_PREFIX)	
	defer fmt.Printf("%s %s\n", messages.SERVER_PREFIX, messages.MESSAGE_DISCONNECTED)
	defer zkCON.Close()
	
	for {		
		fmt.Printf("%s", messages.CLI_PREFIX)		
		var sentence string		
		reader := bufio.NewReader(os.Stdin)		
		sentence, _ = reader.ReadString('\n')
		sentence = strings.Replace(sentence, "\n", "", -1)
		if Command_treatment(sentence){
			if strings.HasPrefix(sentence, reserved_words_commands[0]){
				var separateString = strings.Split(sentence, " ")
				Create(separateString[1], separateString[2])
			}else if strings.HasPrefix(sentence, reserved_words_commands[1]){
				var separateString = strings.Split(sentence, " ")
				Read(separateString[1])
			}else if strings.HasPrefix(sentence, reserved_words_commands[2]){
				var separateString = strings.Split(sentence, " ")
				Exists(separateString[1])
			}else if strings.HasPrefix(sentence, reserved_words_commands[3]){
				var separateString = strings.Split(sentence, " ")
				List(separateString[1])
			}else if strings.HasPrefix(sentence, reserved_words_commands[4]){
				var separateString = strings.Split(sentence, " ")
				Update(separateString[1], separateString[2])
			}else if strings.HasPrefix(sentence, reserved_words_commands[5]){
				var separateString = strings.Split(sentence, " ")
				Delete(separateString[1])
			}else{
				break
			}
		}else{
			fmt.Printf("%s %s\n", messages.ERROR_PREFIX, messages.MESSAGE_INCORRECT_USAGE_COMMAND)
		}
		sentence = ""
	}					
}

