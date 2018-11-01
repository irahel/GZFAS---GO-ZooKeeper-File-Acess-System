package messages

const(
	//Prefix
	LOG_PREFIX string = "[LOG]"
	CLI_PREFIX string = "[Client] "
	ERROR_PREFIX string = "[ERROR ] "
	SERVER_PREFIX string = "[Server] "

	//Message
		//Connection
	MESSAGE_CONNECTED string = "Connection successfully ESTABLISHED "
	MESSAGE_CONNECTED_CLI string = "Client Successfully connected "
	MESSAGE_DISCONNECTED_CLI string = "Client Successfully disconnected "
	MESSAGE_CLIENT_NOT_EXISTS string = "Client not exists "
	MESSAGE_ARCHITECTURE_OK string = "Architecture OK "
	MESSAGE_DISCONNECTED string = "Connection successfully STOPPED "
	MESSAGE_CONNECT_ERROR string = "Connection error "
	MESSAGE_CONNECT_SERVER_ERROR string = "Server Start error"
	MESSAGE_ARCHITECTURE_ERROR string = "Architecture ZooKeeper not defined"
	MESSAGE_SERVER_ALREADY string = "An online server already exists"
	MESSAGE_SERVER_DEFAULT_VALUES string = "Error on setting default server values"

	MESSAGE_STR_I_ERROR string = "String to Float conversion error"
	
	MESSAGE_TCP_A_ERROR string = "Error on solve TCP address."
	MESSAGE_TCP_B_ERROR string = "Error on creating TCP buffer"
	MESSAGE_TCP_C_ERROR string = "Error on solve TCP connection."
	MESSAGE_CONN_A_ERROR string = "Error on accept connection"

	MESSAGE_TCP_C_SUCCESS string = "TCP connection successfully connected."
	MESSAGE_TCP_DC_SUCCESS string = "TCP connection successfully disconnected."

		//Commands
	MESSAGE_COMMAND_EMPTY string = "Empty command"
	MESSAGE_COMMAND_MANY_ARGUMMENTS string = "Command have many argumments"
	MESSAGE_COMMAND_FEW_ARGUMMENTS string = "Command have few argumments"
			//Create    
	MESSAGE_CREATE_SUCESS string = "Node created successfully"
	MESSAGE_CREATE_FAILED string = "Node created failed"
			//Read
	MESSAGE_EXISTS_TRUE string = "Node exists"
	MESSAGE_EXISTS_FALSE string = "Node NOT exists"
			//List
	MESSAGE_LIST string = " childrens path: "
			//Update
	MESSAGE_UPDATE_SUCESS string = "Node updated successfully"
	MESSAGE_UPDATE_FAILED string = "Node NOT updated"
			//Delete
	MESSAGE_DELETE_SUCESS string = "Node deleted successfully"
	MESSAGE_DELETE_FAILED string = "Node NOT deleted"
		//Usage
			//Incorrect
	MESSAGE_INCORRECT_USAGE_LINE_COMMAND string = "Incorrect calling line command"
	MESSAGE_INCORRECT_USAGE_COMMAND string = "Incorrect calling command"
			//Best Usage
	MESSAGE_USAGE_LINE_COMMAND string = "Usage:'''go run Cli.go <Server IP> <Server Port>'''"
	MESSAGE_USAGE_CREATE_COMMAND string = "Usage:'''CREATE/UPDATE <path> <value>'''"
	MESSAGE_USAGE_READ_COMMAND string = "Usage:'''READ/EXISTS/LIST <path> <value>'''"
		//Erros
	MESSAGE_BAD_VERSION string = "Version doesn’t match"
)