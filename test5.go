package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"encoding/json"
	"strconv"
)

//структура url
type StructUrl struct {
	info map[string]string
	url string
}
//принимает: путь до правил, и структуру URL. возвращает: балы, ошибку
func priorityUrl(regulationsFile string, structUrl StructUrl) (float64, string){
	//забираем из файла правила и обрабатываем JSON
	file, err := os.Open(regulationsFile)
	if err != nil{
		os.Exit(1)
		return 0, err.Error()
	}
	defer file.Close()

	data := make([]byte, 64)

	var regulations string
	for{
		n, err := file.Read(data)
		if err == io.EOF{   // если конец файла
			break           // выходим из цикла
		}
		regulations += strings.TrimSpace(string(data[:n]))

	}


	byt := []byte(regulations)
	var datRegulations map[string]interface{}
	if err := json.Unmarshal(byt, &datRegulations); err != nil {
		return 0, err.Error()
	}

	//задаём начальный приоритет
	priority := datRegulations["begin"].(float64)

	//отделяем параметры и сам адрес
	url := strings.Split(structUrl.url, "?")
	//разбиваем адрес
	lastSec := strings.Split(url[0], "/")
	pointDelim := strings.Split(string(lastSec[len(lastSec)-1]), ".")
	//выискиваем окончание адреса и смотрим
	var endElem float64
	//if(len(pointDelim) > 1) {
		endMap := datRegulations["url"].(map[string]interface{})["end"].(map[string]interface{})
		endLElem := strings.ToLower("." + pointDelim[len(pointDelim) - 1])
		if _, ok := endMap[endLElem]; ok {
			endElem = endMap[endLElem].(float64)

		}
	//}
	//выискиваем окончания всех гет параметров
	var getElem float64
	if(len(url) > 1){
		allParams := strings.Split(url[1], "&")
		for _, par := range allParams {
			param := strings.Split(par, "=")
			if(len(param) > 1){
				valParam := strings.Split(param[1], ".")

					getMap := datRegulations["url"].(map[string]interface{})["get"].(map[string]interface{})
					endGet := strings.ToLower("."+valParam[len(valParam) - 1])
					if _, ok := getMap[endGet]; ok {
						getElem = getMap[endGet].(float64)
						break
					}


			}

		}
	}

		//если условия с окончанием адреса и окончанием параметра выполняется, добавляем
		// фиксированное значение из правил в иных случаях, либо одно, либо другуое
	if(endElem > 0 && getElem > 0){
		priority += datRegulations["url"].(map[string]interface{})["get_and_end"].(float64)
	} else if(endElem > 0){
		priority += endElem
	} else if(getElem > 0){
		priority += getElem
	}

		//перебираем уровни доменов и смотрим, есть ли такие в правилах, добаляем соответствующие балы
	domain := strings.Split(string(lastSec[2]), ".")
	lenDom := len(domain)

	for i, levelVal := range domain {
		real_i := int(lenDom) - int(i)
		str_i := strconv.Itoa(real_i)
		domMap := datRegulations["url"].(map[string]interface{})["domain"].(map[string]interface{})
		if _, oklevel := domMap[str_i]; oklevel {
			domMapVal := domMap[str_i].(map[string]interface{})
			if _, oklevelval := domMapVal[levelVal]; oklevelval {
				priority += domMapVal[levelVal].(float64)
			}
		}
	}

	//проибегаемся по свойствам info, если такие в правилах есть, добавляем балы
	infoMap := datRegulations["info"].(map[string]interface{})
	for key, valinfo := range structUrl.info {
		if _, okinf := infoMap[key]; okinf {
			//structUrl
			valInfoMap := infoMap[key].(map[string]interface{})
			if _, okinfval := valInfoMap[valinfo]; okinfval {
				priority += valInfoMap[valinfo].(float64)
			}
		}
	}

	return priority, ""
}


func main() {

	var input1 StructUrl

	inf1 := make(map[string]string)
	inf1["as"] = "19574"
	inf1["as_org"] = "Corporation Service Company"
	inf1["city"] = "Wilmington"
	inf1["country"] = "United States"
	inf1["iso"] = "US"
	inf1["isp"] = "Corporation Service Company"
	inf1["org"] = "Corporation Service Company"
	input1.info = inf1
	input1.url = "http://degeuzen.nl/jeygtgv.exe"


	a, err := priorityUrl("/root/projects/src/regulations.json", input1)

	if(len(err) > 0){
		fmt.Println(err)
	}else{
		fmt.Println(a)
	}



	var input2 StructUrl

	inf2 := make(map[string]string)
	inf2["as"] = "197068"
	inf2["as_org"] = "HLL LLC"
	inf2["city"] = ""
	inf2["country"] = "Russia"
	inf2["iso"] = "RU"
	inf2["isp"] = "HLL LLC"
	inf2["org"] = "HLL LLC"
	input2.info = inf2
	input2.url = "https://download.geo.drweb.com/pub/drweb/windows/katana/1.0/drweb-1.0-katana.exe?download=MSXML3.DLL"

	b, err2 := priorityUrl("/root/projects/src/regulations.json", input2)

	if(len(err2) > 0){
		fmt.Println(err2)
	}else{
		fmt.Println(b)
	}


}