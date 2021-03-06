package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewServer(t *testing.T) {
	Convey("Should get an instance of Server struct", t, func() {
		s := NewServer("google.com", "80")
		So(s, ShouldResemble, Server{
			Host: "google.com",
			Port: "80",
			Storage: &Storage{
				data: make(map[string]interface{}, DefaultStorageSize),
			},
		})
	})
}

func TestRun(t *testing.T) {
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
			So(err, ShouldBeNil)
			connection.Close()

		})
		Convey("Should be able to sent message and receive a response from server", func() {
			host := "localhost"
			port := "7777"
			testString := "I am a program monkey"
			expectString := fmt.Sprintf("Error: Unable to parse message: Command '%s' is invalid\n", testString)

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

			So(string(res[:]), ShouldEqual, expectString)
		})
		Convey("Should be able to insert new data", func() {
			host := "localhost"
			port := "7777"
			type testData struct {
				Input  string
				Expect string
			}
			testDataSet := []testData{
				testData{
					Input:  "PUT;dream;I will be a millionare;string",
					Expect: "Response: Insert dream=I will be a millionare\n",
				},
				testData{
					Input:  "GET;dream;;",
					Expect: "Response: I will be a millionare\n",
				},
				testData{
					Input:  "DELETE;dream;;",
					Expect: "Response: DELETED key=dream\n",
				},
			}

			go func() {
				s := NewServer(host, port)
				s.Run()
			}()

			for _, testCase := range testDataSet {

				tcpAddress, err := net.ResolveTCPAddr("tcp4", host+":"+port)
				checkErr(err, t)

				connection, err := net.DialTCP("tcp", nil, tcpAddress)
				checkErr(err, t)

				_, err = connection.Write([]byte(testCase.Input))
				checkErr(err, t)

				res, err := ioutil.ReadAll(connection)
				checkErr(err, t)

				connection.Close()

				So(string(res[:]), ShouldEqual, testCase.Expect)

			}
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
			So(err, ShouldBeNil)
			So(*m, ShouldResemble, test.Expect)
		}
	})
}
