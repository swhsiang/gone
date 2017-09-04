package server

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/swhsiang/gone/log"
)

// Storage store data
type Storage struct {
	data map[string]interface{}
}

const DefaultStorageSize = 10

// Server define server's property
type Server struct {
	Host    string
	Port    string
	Storage *Storage
}

// NewServer return an instance
func NewServer(host, port string) Server {
	return Server{Host: host, Port: port, Storage: &Storage{
		data: make(map[string]interface{}, DefaultStorageSize),
	}}
}

// Message define Message's property
type Message struct {
	Command string
	// FIXME ignore STATS command first
	Key       string
	Value     string
	ValueType string
}

var re = regexp.MustCompile(`^(PUT|GET|DELETE)\;(\w+)\;(.*)\;(string|int|array)?$`)

func parseMessage(message string) (*Message, error) {
	_res := re.FindAllStringSubmatch(message, -1)
	if len(_res) == 0 {
		return &Message{}, fmt.Errorf("Command '%s' is invalid", message)

	}
	res := _res[0]
	if len(res) != 5 {
		return nil, errors.New("The message is invalid")
	}
	m := &Message{
		Command:   res[1],
		Key:       res[2],
		Value:     res[3],
		ValueType: res[4],
	}
	return m, nil
}

// Run start server
func (s *Server) Run() {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", s.Host+":"+s.Port)
	if err != nil {
		log.Errorf("Unable to resolve tcp address: %s", err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Errorf("Unable to listen host %s port %s : %s", s.Host, s.Port, err.Error())
		return
	}

	log.Infof("Listening on %s:%s ...", s.Host, s.Port)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Warnf("Unable to receive request: %s", err.Error())
			continue
		}
		_ = connection.SetReadDeadline(time.Now().Add(time.Second * 15))

		buf := make([]byte, 4*1024)
		for {

			n, err := connection.Read(buf)
			if err != nil || n == 0 {
				break
			}
			query := strings.TrimSpace(string(buf[:n]))
			res, err := s.Write(query)
			if err != nil {
				connection.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
				continue
			}

			connection.Write([]byte(fmt.Sprintf("Response: %s\n", res.Res)))
		}

		connection.Close()
	}
}

// Response object
type Response struct {
	Res string
	Err error
}

func (s *Server) Write(buf string) (Response, error) {
	message, err := parseMessage(buf)
	if err != nil {
		return Response{}, fmt.Errorf("Unable to parse message: %s", err.Error())
	}

	switch message.Command {
	case "PUT":
		log.Debugf("message: %+v", message)
		s.Storage.data[message.Key] = message.Value
		return Response{Res: fmt.Sprintf("Insert %s=%s", message.Key, message.Value)}, nil
	case "GET":
		if val, ok := s.Storage.data[message.Key]; ok {
			return Response{Res: fmt.Sprintf("%v", val)}, nil
		}
		return Response{Res: ""}, nil

	case "DELETE":
		if _, ok := s.Storage.data[message.Key]; ok {
			delete(s.Storage.data, message.Key)
			return Response{Res: fmt.Sprintf("DELETED key=%s", message.Key)}, nil
		}
		return Response{Res: ""}, nil
	}

	return Response{}, errors.New("Invalid operation")
}
