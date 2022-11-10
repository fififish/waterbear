/*
System logs 
*/

package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type TypeOfLog int
var id string
var logOpt int
var homepath string
var err error

var Normal_LogFile *os.File
var Error_LogFile *os.File
var Evalue_LogFile *os.File
var Info_LogFile *os.File

var NormalLogger *log.Logger 
var ErrorLogger *log.Logger 
var EvalLogger *log.Logger 
var InfoLogger *log.Logger 

const(
	NormalLog     TypeOfLog = 0
	ErrorLog      TypeOfLog = 1
	EvaluationLog TypeOfLog = 2
	InfoLog       TypeOfLog = 3
)

/*
Used as the directory for logs
*/
func SetID(rid string){
	id = rid
	exepath, err := os.Executable()
	if err != nil {
		log.Fatalln("[WriteLog Error] cannot get absolute path of executable.",err)
	}
	
	p1 := path.Dir(exepath)
	homepath =  path.Dir(p1)
}

func SetLogOpt(LOpt int)  {
	logOpt = LOpt
}


func PrintLog(verbose bool, logType TypeOfLog, msg string){
	if !verbose {
		return
	}
	//log.Printf(msg)
	//if logType != ErrorLog && logType != EvaluationLog{
	//	return
	//}

	var fpath string
	var fileName string
	//fmt.Printf("   %d   %s\n",logOpt,msg)
	if logOpt == 0{
		switch logType{
		case NormalLog:
			if NormalLogger == nil{
				fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
				fileName = fmt.Sprintf("%s_Normal.log",time.Now().Format("20060102"))
				if !IsExist(fpath) {
					err := CreateDir(fpath)
					if err != nil{
						log.Printf("[WriteLog Error] create log file error!",err)
						return 
					}
					fmt.Println("create "+fpath)
				}
				
				fmt.Println("open ",fpath+fileName)
				Normal_LogFile, err = os.OpenFile(fpath+fileName,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Printf("[WriteLog Error] open log file error!", err)
					return
				}
				NormalLogger = log.New(Normal_LogFile,"",log.LstdFlags)
			}

			//debugLog := log.New(Normal_LogFile,"",log.LstdFlags)
			NormalLogger.Println(msg)
		case ErrorLog:
			if ErrorLogger == nil{
				fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
				fileName = fmt.Sprintf("%s_Error.log",time.Now().Format("20060102"))
				if !IsExist(fpath) {
					err := CreateDir(fpath)
					if err != nil{
						log.Printf("[WriteLog Error] create log file error!",err)
						return
					}
					fmt.Println("create "+fpath)

				}
				fmt.Println("open ",fpath+fileName)
				Error_LogFile, err = os.OpenFile(fpath+fileName,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Printf("[WriteLog Error] open log file error!", err)
					return
				}
				ErrorLogger = log.New(Error_LogFile,"",log.LstdFlags)
			}

			//debugLog := log.New(Error_LogFile,"",log.LstdFlags)
			ErrorLogger.Println(msg)
		case EvaluationLog:
			if EvalLogger == nil{
				fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
				fileName = fmt.Sprintf("%s_Eva.log",time.Now().Format("20060102"))
				if !IsExist(fpath) {
					err := CreateDir(fpath)
					if err != nil{
						log.Printf("[WriteLog Error] create log file error!",err)
						return
					}
					fmt.Println("create "+fpath)
				}
				fmt.Println("open ",fpath+fileName)
				Evalue_LogFile, err = os.OpenFile(fpath+fileName,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Printf("[WriteLog Error] open log file error!", err)
					return
				}
				EvalLogger = log.New(Evalue_LogFile,"",log.LstdFlags)
			}

			//debugLog := log.New(Evalue_LogFile,"",log.LstdFlags)
			EvalLogger.Println(msg)
		case InfoLog:
			if Info_LogFile == nil{
				fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
				fileName = fmt.Sprintf("%s_Info.log",time.Now().Format("20060102"))
				if !IsExist(fpath) {
					err := CreateDir(fpath)
					if err != nil{
						log.Printf("[WriteLog Error] creat log file error!",err)
						return
					}
					fmt.Println("create "+fpath)
				}
				fmt.Println("open ",fpath+fileName)
				Info_LogFile, err = os.OpenFile(fpath+fileName,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Printf("[WriteLog Error] open log file error!", err)
					return
				}
				InfoLogger = log.New(Info_LogFile,"",log.LstdFlags)
			}

			//debugLog := log.New(Info_LogFile,"",log.LstdFlags)
			InfoLogger.Println(msg)
		}
	}else if logOpt == 1{
		switch logType{
		case NormalLog:
			fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
			fileName = fmt.Sprintf("%s_Normal.log",time.Now().Format("20060102"))
		case ErrorLog:
			fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
			fileName = fmt.Sprintf("%s_Error.log",time.Now().Format("20060102"))
		case EvaluationLog:
			fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
			fileName = fmt.Sprintf("%s_Eva.log",time.Now().Format("20060102"))
		case InfoLog:
			fpath = fmt.Sprintf(homepath+"/var/log/%s/", id)
			fileName = fmt.Sprintf("%s_Info.log",time.Now().Format("20060102"))
		}

		if !IsExist(fpath) {
			err := CreateDir(fpath)
			if err != nil{
				log.Printf("[WriteLog Error] create log file error!",err)
			}
			fmt.Println("create"+fpath)
		}
		logFile,err  := os.OpenFile(fpath+fileName,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("[WriteLog Error] open log file error!",err)
			return
		}
		defer logFile.Close()
		debugLog := log.New(logFile,"",log.LstdFlags)
		debugLog.Println(msg)
	}

}

//Create Dir
func CreateDir(fpath string) error {
	err := os.MkdirAll(fpath, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(fpath, os.ModePerm)
	return nil
}

//Determines whether folders exist,=
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

