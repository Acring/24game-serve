package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"net/http"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"strings"
	"html/template"
	"dot24_server/client"
	"dot24_server/frame"
	"os/exec"
	"os"
	"path/filepath"
)



var manager = client.GetInstance()

func wsPage(res http.ResponseWriter, req *http.Request) { // 登录页面, 注册一个连接
	fmt.Print("请求登录\n")
	// 将http连接升级为socket连接
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)

	if err != nil {
		fmt.Print(err.Error())
		http.NotFound(res, req)
		return
	}
	// 分配id
	uid, _ := uuid.NewV4()
	strid := fmt.Sprintf("%s", uid)
	newclient := client.Client{Id: strid, Socket: conn, Send: make(chan *frame.Frame), Status: "waiting"}
	// 向管理类注册新的客户端连接
	fmt.Print("注册完成\n")
	go newclient.Read()
	go newclient.Write()
	manager.Register <- &newclient
}

type View struct{
	indexTemplate *template.Template
	templates *template.Template
}
var view = View{}

func (view *View)initView(){
	proDir,err := GetProDir()
	var allfile []string
	files, err := ioutil.ReadDir(proDir + "/dist/")
	if err != nil {
		fmt.Print("获取静态资源失败")
		return
	}

	for _, file := range files {
		fileName := file.Name()
		//log.Println(fileName)
		if strings.HasSuffix(fileName, ".html") {

			if err != nil{
				fmt.Println("获取本地路径失败")
			}
			allfile = append(allfile, proDir + "/dist/"+fileName)
		}
	}
	view.templates = template.Must(template.ParseFiles(allfile...))
	view.indexTemplate = view.templates.Lookup("index.html")

}

func (view *View)loadTemplate(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if r.Method == "GET" {
		view.indexTemplate.Execute(w, nil)
		return
	}
}
// GetProDir 用于获取项目根目录
func GetProDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)

	if err != nil {
		return "", err
	}

	end := strings.LastIndex(path, string(os.PathSeparator))
	proPath := path[:end]

	return proPath, nil
}

func main() {
	view.initView()  // 初始化静态前端文件
	fmt.Println("Starting application...")

	manager.Init()
	go manager.Matching()
	go manager.Start()

	http.HandleFunc("/ws", wsPage)
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("dist"))))
	http.HandleFunc("/", view.loadTemplate)
	http.ListenAndServe(":6606", nil)

	//index.html bash href="/src/"
}