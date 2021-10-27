package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"web/genWords"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	booksAPI := app.Party("/books")
	rootPath := app.Party("/")
	{
		booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/books
		booksAPI.Get("/", list)
		// POST: http://localhost:8080/books
		booksAPI.Post("/", create)
		rootPath.Get("healthz", healthz)
		rootPath.Get("ip", fetchIp)
		rootPath.Get("words", getWords)

	}
	// runForever()
	run()
	app.Listen(":8090")

}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func list(ctx iris.Context) {
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}

	ctx.JSON(books)
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

func healthz(ctx iris.Context) {
	ctx.Text("ok!")
}

func fetchIp(ctx iris.Context) {
	ip1, err := GetOutBoundIP()
	if err != nil {
		ctx.StatusCode(400)
		ctx.Text(err.Error())
		return
	}
	fmt.Println("the inner ip is:", ip1)

	ip2, err := getIp()
	if err != nil {
		ctx.StatusCode(400)
		ctx.Text(err.Error())
		return
	}
	fmt.Println("the outer ip is:", ip2)
	ctx.StatusCode(200)
	ctx.Text("the inner ip is: %s  \nthe outer ip is: %s", ip1, ip2)
}
func create(ctx iris.Context) {
	var b Book
	err := ctx.ReadJSON(&b)
	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Book creation failure").DetailErr(err))
		// TIP: use ctx.StopWithError(code, err) when only
		// plain text responses are expected on errors.
		return
	}

	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}

// 获取IP用于测试!

func getIp() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func getWords(ctx iris.Context) {
	w, err := genWords.GetWords()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.JSON(w)
}

func runForever() {
	// while true
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Second * 3)
		w, err := genWords.GetWords()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(w)
	}
}

func run() {
	go runForever()
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		runForever()
		os.Exit(0)
	}()
}
