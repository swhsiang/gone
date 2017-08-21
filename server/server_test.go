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
func TestCommandParser(t *testing.T) {
	Convey("Parse message", t, func() {
		type testData struct {
			Input  string
			Expect Message
		}
		testDataSet := []testData{
			testData{
				Input: "PUT;dream;I will be a millionare;string",
				Expect: Message{
					Command:   "PUT",
					Key:       "dream",
					Value:     "I will be a millionare",
					ValueType: "string",
				},
			},
			testData{
				Input: "GET;dream;;",
				Expect: Message{
					Command:   "GET",
					Key:       "dream",
					Value:     "",
					ValueType: "",
				},
			},
		}

		for _, test := range testDataSet {
			m, err := parseMessage(test.Input)
			ShouldBeNil(err)
			ShouldEqual(*m, test.Expect)
		}
	})
}
