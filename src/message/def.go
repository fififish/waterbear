
package message

type TypeOfMessage int

const (
	RBC_ALL       	TypeOfMessage = 1
	RBC_SEND      	TypeOfMessage = 2
	RBC_ECHO      	TypeOfMessage = 3
	RBC_READY     	TypeOfMessage = 4
	ABA_ALL       	TypeOfMessage = 5
	ABA_BVAL      	TypeOfMessage = 6
	ABA_AUX       	TypeOfMessage = 7
	ABA_CONF      	TypeOfMessage = 8
	ABA_FINAL       TypeOfMessage = 9
	PRF           	TypeOfMessage = 10
	ECRBC_ALL     	TypeOfMessage = 11
	CBC_ALL       	TypeOfMessage = 12
	CBC_SEND      	TypeOfMessage = 13
	CBC_REPLY     	TypeOfMessage = 14
	CBC_ECHO	  	TypeOfMessage = 15
	CBC_EREPLY    	TypeOfMessage = 16
	CBC_READY     	TypeOfMessage = 17
	MVBA_DISTRIBUTE TypeOfMessage = 18
	EVCBC_ALL       TypeOfMessage = 19
	RETRIEVE    	TypeOfMessage = 20
)


type ProtocolType int 

const (
	RBC       	ProtocolType = 1
	ABA  	  	ProtocolType = 2
	ECRBC     	ProtocolType = 3
	CBC       	ProtocolType = 4
	EVCBC     	ProtocolType = 5
	MVBA      	ProtocolType = 6
)

type VCBCType int

const (
	DEFAULT_HASH			VCBCType = 0
	MERKLE					VCBCType = 1
)
