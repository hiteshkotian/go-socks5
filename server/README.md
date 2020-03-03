# Server code

The server will act as a proxy to tunnel TCP Request.

Flow :

CLIENT -> <CONNECT_PKT> -> SERVER
                             |
                             V
                           Validates the connection
                         request and also authenicates the user
                             |
                             V
CLIENT <- <CONNECT_ACPT> <- Establishes a tunnel with the 
                            server
                    

CONNECT_PKT:

| <ID> | 0x00 | <SERVER> | <PORT> | <NAMESERVER> | <AUTH_TYPE> | <Authentication_TOKEN> |


CONNECT_ACPT:


|  8   |  
+------+---------+----------+--------+------------+----------+
| <ID> | 0x01 | <SERVER> | <PORT> | <PROXY_ID> | <STATUS> |

CONNECT -> 

STATUS => OK -> 0x00, 
          NOAUTH -> 0x01,
          INVALID_MSG -> 0x02
          SERVER_UNREACHABLE -> 0x03

========

State machine :


Client -> UNAUTH -> Authenticate with the server -> awaiting connection -> proxying -> end connection -> Terminated


Server -> Awaiting connection -> connecting (auth + validation) -> connected -> proxying -> terminating


Auth ->

Client -> <AUTH_INITIATE> -> SERVER

SERVER -> <AUTH_MAN> -> CLIENT

CLIENT -> <AUTH_REQ> -> SERVER

SERVER -> <AUTH_RESP> -> CLIENT


AUTH_INITIATE:

| ID | AUTH_INIT