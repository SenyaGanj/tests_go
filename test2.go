package main
import (
	"fmt"
	"os/exec"
	"strings"
	"regexp"
	"time"
)

//возвращает 4 значения: время создания, список серверов, имя организации, текст ошибки
func whoisFunc(name string) (string, []string, string, string){
	var arServers []string
	name = strings.TrimSpace(name)
	//вызываем whois
	cmd := exec.Command("whois", name)
	stdout, err := cmd.Output()
	//обрабатываем ошибку
	if err != nil {
		return "", arServers, "", err.Error()
	}

	//разбиваем пришедшую строку на массив по отступам
	output := strings.Split(string(stdout), "\n")

	//переменные для возврата
	var nameOrg string
	var createTime string

	for _, val := range output{

		val := strings.ToLower(val)
		//пропускаем все строки с "%", т к это комментарии (информация сверху)
		v, _ := regexp.MatchString("%", val)
		if(v){
			continue
		}
		//выходим из массива, если натыкаемся на "last update", т к это последняя значимая строка
		l_update, _ := regexp.MatchString("last update", val)
		if(l_update){
			break
		}

		//разделяем строки чере ": ", для получения ключа и его значения
		elem := strings.Split(val, ": ")
		if(len(elem) <= 1){
			continue
		}
		//ищем время создания
		flagCreateTime, _ := regexp.MatchString("creat", elem[0])
		if(flagCreateTime){
			//не нашёл вывод в структуре, как на Питоне((( только парсинг в другой формат
			createTime = strings.TrimSpace(elem[1])
		}

		//ищем список серверов, и скрладываем в массив arServers
		flagServers, _ := regexp.MatchString("server", elem[0])
		flagWhois, _ := regexp.MatchString("whois", elem[0])
		if(flagServers && !flagWhois){
			arServers = append(arServers, strings.TrimSpace(elem[1]))
		}

		//ищем имя организации
		flagNameOrg, _ := regexp.MatchString("org", elem[0])
		if(flagNameOrg){
			nameOrg = strings.TrimSpace(elem[1])
		}

	}

	return createTime, arServers, nameOrg, ""
}

func main() {
	//t, _ := time.Parse("2018-09-05T18:16:32Z", "2018-09-05T18:16:32Z")
	//fmt.Println(t)
	a, b, c, err := whoisFunc("drweb.ru")
	//fmt.Printf("", )
	if(len(err) > 0){
		fmt.Println(err)
	}else {
		fmt.Println(a)
		fmt.Println(b)
		fmt.Println(c)
	}
}