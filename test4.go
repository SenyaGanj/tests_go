package main

import (
	"fmt"
	"os"
	"io"
	"regexp"
	"strings"
)

//возвращаемая структура
type Event struct {
	PID, SF_AT, SF_TEXT string
}
//принимает путь к логу и событие
func parsLog(log, event string) ([]Event, string){

	var arResult []Event
	//считываем лог
	file, err := os.Open(log)
	if err != nil{
		os.Exit(1)
		return arResult, err.Error()
	}
	defer file.Close()

	data := make([]byte, 64)

	var logTextArr []string
	for{
		n, err := file.Read(data)
		if err == io.EOF{   // если конец файла
			break           // выходим из цикла
		}
		logTextArr = append(logTextArr, strings.TrimSpace(string(data[:n])))

	}
	//объединяем весь файл в одну строку
	LogText := strings.Join(logTextArr, " ")
	//разбиваем через "/n"
	logTextArr = strings.Split(LogText, "\n")
	var arSF_TEXT []string
	var OneStruct Event
	for _, str := range logTextArr {
		//если в структуре нет PID добавляем его
		if(len(OneStruct.PID) == 0){
			//если есть событие в строке, добавляем PID, SF_AT в структуру
			v, _ := regexp.MatchString(event, str)
			if(v){
				//ищем PID
				r := regexp.MustCompile("F-(\\d+):")
				PIDreg := r.FindStringSubmatch(str)
				if(len(PIDreg) > 1) {
					OneStruct.PID = strings.TrimSpace(PIDreg[1])
				}

				//ищем SF_AT
				SF_ATregT := event + "at (\\d+)"
				r = regexp.MustCompile(SF_ATregT)
				SF_ATreg := r.FindStringSubmatch(str)
				if(len(SF_ATreg) > 1) {
					OneStruct.SF_AT = strings.TrimSpace(SF_ATreg[1])

				}
			}
		} else {
			//если в структуре есть PID, ищем поле SF_TEXT из каждой строки, пока событие не закончится
			pid_ex, _ := regexp.MatchString("F-" + OneStruct.PID, str)
			if(pid_ex){
				r := regexp.MustCompile("Dump: (.+)")
				SF_TEXTreg := r.FindStringSubmatch(str)
				if(len(SF_TEXTreg) > 1) {
					arSF_TEXT = append(arSF_TEXT, strings.TrimSpace(SF_TEXTreg[1]))
					arSF_TEXT = append(arSF_TEXT, "\n")
				}
			} else {
				//видим, что в строке нет такого PID, обнуляем структуру, для поиска нового события
				OneStruct.SF_TEXT = strings.Join(arSF_TEXT, " ")
				arResult = append(arResult, OneStruct)
				OneStruct = Event{"","",""}
			}
		}

	}


	return arResult, ""
}


func main() {
	a, err := parsLog("/root/projects/src/log.txt", "Segmentation fault")

	if(len(err) > 0){
		fmt.Println(err)
	}else{
		fmt.Println(a)
	}

}