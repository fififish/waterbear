package cryptolib


import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"os"
)


type CryptoLibrary int
var nid int64

const (
	NoCrypto CryptoLibrary = 0 // Return true for all operations. Used for testing and when there is an authenticated channel. 
)

var TypeOfCrypto_name = map[int]CryptoLibrary{
	0: NoCrypto,
}

var cryptoOption CryptoLibrary



func GenMAC(id int64, msg []byte) []byte {
	result := []byte("")

	mac := hmac.New(sha256.New,[]byte("123456789"))
	mac.Write(msg)
	result = mac.Sum(nil)

	return result
}

func VerifyMAC(id int64, msg []byte, sig []byte) bool {
	result := false

	
	mac := hmac.New(sha256.New,[]byte("123456789"))
	mac.Write(msg)
	expectedMac := mac.Sum(nil)
	result = hmac.Equal(expectedMac,sig)

	return result
}

func StartCrypto(id int64, cryptoOpt int) {
	var exist bool
	nid = id
	cryptoOption,exist = TypeOfCrypto_name[cryptoOpt]
	if !exist{
		log.Fatalf("Crypto option is not supported")
		os.Exit(1)
	}

	switch cryptoOption{
	case NoCrypto:
		log.Printf("Alert. Not using any crypto library.")
	default:
		log.Fatalf("The crypto library is not supported by the system")
	}
}
