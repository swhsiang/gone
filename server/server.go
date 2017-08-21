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

// Server define server's property
type Server struct {
	Host string
	Port string
}

// NewServer return an instance
func NewServer(host, port string) Server {
	return Server{Host: host, Port: port}
}

// Message define Message's property
type Message struct {
	Command string
	// FIXME ignore STATS command first
	Key string
	// FIXME Support string type first
	Value     string
	ValueType string
}

var re = regexp.MustCompile(`(PUT|GET|DELETE)\;(\w+)\;(.*)\;(string|int|array)?$`)

func parseMessage(message string) (*Message, error) {

	res := re.FindAllStringSubmatch(message, -1)[0]
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
		connection.SetReadDeadline(time.Now().Add(time.Second * 15))

		buf := make([]byte, 4*1024)
		for {

			n, err := connection.Read(buf)
			if err != nil || n == 0 {
				break
			}
			query := strings.TrimSpace(string(buf[:n]))
			log.WithFields(log.Fields{
				"address": fmt.Sprintf("%+v", connection.RemoteAddr()),
			}).Infof("Data: %s", query)

			connection.Write([]byte(fmt.Sprintf("Response: %s\n", query)))

		}

		connection.Close()
	}
}
