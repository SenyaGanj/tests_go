package main
import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func parsFileName (path string) ([]string, string)  {

	//массив с результатами
	var arResult []string
	//вытаскиваем файлы из каталога
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return arResult, err.Error()
	}

	//пропускаем через регулярное выражение название каждого файла
	r := regexp.MustCompile(`image-(\S+)-\d{8}T\d{6}.tar.gz`)
	for _, file := range files {
		str := r.FindStringSubmatch(file.Name())
		if (len(str) > 1) {
			arResult = append(arResult, str[1])
		}
	}

	return arResult, ""
}

func main() {
	arr, err := parsFileName("/root/projects/src/catalog/")
	if(len(err) > 0) {

		fmt.Println(err)

	}else{
		for _, v := range arr {
			fmt.Println(v)
		}
	}
}