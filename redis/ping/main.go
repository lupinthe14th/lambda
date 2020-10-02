package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Println(ipnet.IP.String())
			}
		}
	}
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	log.Printf("addr: %v port: %v", addr, port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}
	defer conn.Close()

	out := make([]string, 0)
	for _, s := range []string{"ping", "quit"} {
		_, err = conn.Write([]byte(s + "\n"))
		if err != nil {
			log.Fatal(err)
			return events.APIGatewayProxyResponse{}, err
		}
		log.Print(s)
		reply, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return events.APIGatewayProxyResponse{}, err
		}
		out = append(out, reply)
	}
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       strings.Join(out, ","),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
