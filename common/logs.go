package common

import (
	"os"
	"log"
	"io"
)

var (
	Debug   *log.Logger
	Info    *log.Logger

)
func CheckError(e error) {
	Init(os.Stdout, os.Stdout)
	if e != nil {
	Info.Println(e)
	}
}




func Init(
	infoHandle io.Writer, debugHandle io.Writer) {

    Debug = log.New(debugHandle,
        "DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
		
	Info = log.New(infoHandle,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)

}