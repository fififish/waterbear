/*
generating threshold prf keys
*/

package main

import (
	"os"
	"flag"
	"log"
	"utils"
	prf "cryptolib/threshprf"
)

const (
	helpText_keygen = `
Generating keys for threshold prf
keygen [n] [k]
n: number of replicas, k: threshold
`
)

func main() {
	
	helpPtr := flag.Bool("help", false, helpText_keygen)
	flag.Parse()
	
	if *helpPtr || len(os.Args) < 2{
		log.Printf(helpText_keygen)
		return
	}
	
	n,_ := utils.StringToInt64(os.Args[1])
	k,_ := utils.StringToInt64(os.Args[2])

	log.Printf("[Keygen] Generating keys for n=%v, k=%v.", n, k)
	prf.SetHomeDir()
	prf.Init_key_dealer(n,k)
	
}