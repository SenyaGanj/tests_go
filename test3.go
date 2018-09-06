package main

import (
    "fmt"
    "golang.org/x/net/html"
    "net/http"
)

//возвращает массив данных и ошибку
func GetScriptTag(url string) ([]string, string){

    //забираем содержимое страницы
    resp, err := http.Get(url)
    var arResult []string
    if err != nil {
        return arResult, err.Error()
    }

    defer resp.Body.Close()
    doc, err := html.Parse(resp.Body)
    if err != nil {
        return arResult, err.Error()
    }

    //разбираем разметку с помощью библиотеки golang.org/x/net/html
    var f func(*html.Node)
    f = func(n *html.Node) {
        //вытаскиваем все теги script
        if n.Type == html.ElementNode && n.Data == "script" {
            flagSrc := false
            //пробегаем по каждому, и если есть аттрибут src, забираем его,
            // flagSrc = true, если нет, забираем содержимое flagSrc = false
            for _, val := range(n.Attr) {
                if(val.Key == "src"){
                    arResult = append(arResult, val.Val)
                    flagSrc = true
                    break
                }
            }
            // в зависимости от флага flagSrc заполняем массив результата
            if(!flagSrc){
                arResult = append(arResult, n.FirstChild.Data)
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    return arResult, ""
}

func main() {
    a, err := GetScriptTag("http://sen.mcart.ru/test.php")

    if(len(err) > 0){
        fmt.Println(err)
    }else{
        fmt.Println(a)
    }

}