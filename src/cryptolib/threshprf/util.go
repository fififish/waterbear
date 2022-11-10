package threshprf

import (
	"fmt"
	"io/ioutil"
	"log"
	logging "logging"
	"os"
	"path"

	word "cryptolib/word"
)

var homepath string
var nsk []byte 
var nvkx []byte
var nvky []byte

func SetHomeDir() {
	exepath, err := os.Executable()
	if err != nil {
		log.Fatalf("[Sdibc.cn] Cannot find absolute path of executable")
	}

	p1 := path.Dir(exepath)
	homepath = path.Dir(p1)
}

func GenPath(id int64) string {
	// if id == -1 {
	// 	return fmt.Sprintf(homepath + "/etc/newclient/")
	// } else if id == -2 {
	// 	return fmt.Sprintf(homepath + "/etc/newserver/")
	// }
	return fmt.Sprintf(homepath+"/etc/thresprf_key/%d/", id)
}

//Create Dir. Copying from storage to make sdcsm2 library independent
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

//Determines whether folders exist. Copying from storage to make smcgo library independent
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

/*存储n个用户的vk//*/
func Store_share(share []byte, id int64) {
	if len(share) != 128 {
		return
	}

	path := GenPath(int64(id))
	var sharefileName = fmt.Sprintf("share")

	if !IsExist(path) {
		err := CreateDir(path)
		if err != nil {
			log.Fatalf("[Store Share Error] create file error!", err)
		}
		fmt.Println("creating " + path)
	}

	err := ioutil.WriteFile(path+sharefileName, share, 0644)
	if err != nil {
		log.Fatalf("[Store Share Error] write share error!", err)
		return
	}

}

//*下载用户id的vk //*/
func LoadshareFromFiles(id int64) (share []byte) {

	path := GenPath(id)
	var sharefileName = fmt.Sprintf("share")

	//p := fmt.Sprintf("loading from %s", path)
	//logging.PrintLog(true, logging.NormalLog, p)
	var err error

	share, err = ioutil.ReadFile(path + sharefileName)
	if err != nil {
		//p := fmt.Sprintf("[Loadvk Error] open vk file error!", err)
		//logging.PrintLog(true, logging.ErrorLog, p)
		return nil
	}
	return
}

/*存储n个用户的vk,sk，此函数只针对dealer//*/
func Store_key_dealer(vk, sk []uint64, n int64) {
	if len(vk) != 8*int(n) || len(sk) != 4*int(n) {
		return
	}

	var vkb, skb []byte
	vkb = make([]byte, 64)
	skb = make([]byte, 32)
	var sk64, vk_x, vk_y [4]uint64

	for id := 0; id < int(n); id++ {
		vk_x[0] = vk[8*id]
		vk_x[1] = vk[8*id+1]
		vk_x[2] = vk[8*id+2]
		vk_x[3] = vk[8*id+3]
		vk_y[0] = vk[8*id+4]
		vk_y[1] = vk[8*id+5]
		vk_y[2] = vk[8*id+6]
		vk_y[3] = vk[8*id+7]
		vkb = word.U64toByte_256(vk_x)
		vkb = append(vkb, word.U64toByte_256(vk_y)...)

		sk64[0] = sk[4*int(id)]
		sk64[1] = sk[4*int(id)+1]
		sk64[2] = sk[4*int(id)+2]
		sk64[3] = sk[4*int(id)+3]
		skb = word.U64toByte_256(sk64)

		path := GenPath(int64(id))
		var vkfileName = fmt.Sprintf("ver.key")
		var skfileName = fmt.Sprintf("sec.key")

		if !IsExist(path) {
			err := CreateDir(path)
			if err != nil {
				log.Fatalf("[KeyGen Error] create file error!", err)
			}
			fmt.Println("creating " + path)
		}

		err := ioutil.WriteFile(path+vkfileName, vkb, 0644)
		if err != nil {
			log.Fatalf("[Store_key_dealer Error] write vk error!", err)
			return
		}
		err = ioutil.WriteFile(path+skfileName, skb, 0644)
		if err != nil {
			log.Fatalf("[Store_key_dealer Error] write sk error!", err)
			return
		}
	}
}

/*存储n个用户的vk//*/
func Store_vk_user(vk []byte, n int64) {
	if len(vk) != int(n)*64 {
		return
	}

	var vkb []byte
	for id := 0; id < int(n); id++ {
		vkb = vk[id*int(n) : id*int(n)+64]

		path := GenPath(int64(id))
		var vkfileName = fmt.Sprintf("ver.key")

		if !IsExist(path) {
			err := CreateDir(path)
			if err != nil {
				log.Fatalf("[KeyGen Error] create file error!", err)
			}
			fmt.Println("creating " + path)
		}

		err := ioutil.WriteFile(path+vkfileName, vkb, 0644)
		if err != nil {
			log.Fatalf("[Store_vk_user Error] write vk error!", err)
			return
		}
	}
}

//*下载用户id的vk //*/
func LoadvkFromFiles(id int64) (vk_x, vk_y []byte) {

	path := GenPath(id)
	var vkfileName = fmt.Sprintf("ver.key")

	//p := fmt.Sprintf("loading from %s", path)
	//logging.PrintLog(true, logging.NormalLog, p)
	var err error

	vkb, err := ioutil.ReadFile(path + vkfileName)
	if err != nil {
		p := fmt.Sprintf("[Loadvk Error] open vk file error!", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return nil, nil
	}

	vk_x = vkb[0:32]
	vk_y = vkb[32:64]

	return vk_x, vk_y
}

//*下载用户id的sk, vk //*/
func LoadkeyFromFiles(id int64) (sk, vk_x, vk_y []byte) {

	path := GenPath(id)
	var vkfileName = fmt.Sprintf("ver.key")
	var skfileName = fmt.Sprintf("sec.key")

	//p := fmt.Sprintf("loading from %s", path)
	//logging.PrintLog(true, logging.NormalLog, p)
	var err error

	sk, err = ioutil.ReadFile(path + skfileName)
	if err != nil {
		p := fmt.Sprintf("[Loadvk Error] open sk file error! %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return nil, nil, nil
	}

	vkb, err := ioutil.ReadFile(path + vkfileName)
	if err != nil {
		p := fmt.Sprintf("[Loadvk Error] open vk file error!", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return nil, nil, nil
	}

	vk_x = vkb[0:32]
	vk_y = vkb[32:64]
	nsk = sk 
	nvkx = vk_x 
	nvky = vk_y
	
	return sk, vk_x, vk_y
}
