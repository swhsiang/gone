package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewServer(t *testing.T) {
	Convey("Start server", t, func() {
		Convey("Test connection", func() {
			host := "localhost"
			port := "7777"

			go func() {
				s := NewServer(host, port)
				s.Run()
			}()

			tcpAddress, err := net.ResolveTCPAddr("tcp4", host+":"+port)
			if err != nil {
				t.Errorf("err should be nil: %s", err.Error())
				t.Fail()
			}
			connection, err := net.DialTCP("tcp", nil, tcpAddress)
			ShouldBeNil(nil)
			connection.Close()

		})
		Convey("Should be able to sent message and receive a response from server", func() {
			host := "localhost"
			port := "7777"
			testString := "I am a program monkey"
			expectString := fmt.Sprintf("Response: %s\n", testString)

			go func() {
				s := NewServer(host, port)
				s.Run()
			}()

			tcpAddress, err := net.ResolveTCPAddr("tcp4", host+":"+port)
			checkErr(err, t)

			connection, err := net.DialTCP("tcp", nil, tcpAddress)
			checkErr(err, t)

			_, err = connection.Write([]byte(testString))
			checkErr(err, t)

			res, err := ioutil.ReadAll(connection)
			checkErr(err, t)
			fmt.Println(string(res[:]))

			connection.Close()
			ShouldEqual(string(res[:]), expectString)
		})
	})
}

func checkErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf("err should be nil: %s", err.Error())
		t.Fail()
	}
}
