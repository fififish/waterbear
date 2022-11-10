package cryptolib

import (
	"log"
	"utils"
	prf "cryptolib/threshprf"
)

var k int64
//var sharemap utils.IntBytesMap 

func GenPRFShare(input []byte) []byte{
	share := prf.Compute_share_own(input)
	return share
}

//TODO: verify this is correct
func CombineShares(idarr []int64, shares [][]byte) []byte{
	prfval := prf.Compute_prf_from_shares(idarr,k,shares)
	return prfval
}

func VerifyShare(C []byte, nodeid int64, share []byte) bool{
	return prf.Verify_share_node(C,nodeid,share)
}



func StartThSig(id int64, threshold int){
	nid = id 
	k = utils.IntToInt64(threshold)
	prf.SetHomeDir()
	prf.LoadkeyFromFiles(id)
	log.Printf("Starting threshold signature setup")
}