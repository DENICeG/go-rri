package client

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/satori/go.uuid"
	"github.com/sebidude/go-rri/messagetypes"
	"gopkg.in/yaml.v2"
)

type RriClient struct {
	Username   string
	Password   string
	Address    string
	Connection *tls.Conn
	TlsConfig  *tls.Config
}

func NewRriClient(user string, pass string, address string) (*RriClient, error) {
	client := &RriClient{
		Username: user,
		Password: pass,
		Address:  address,
	}
	client.TlsConfig = &tls.Config{
		MinVersion:         tls.VersionSSL30,
		CipherSuites:       []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA},
		InsecureSkipVerify: true,
	}
	var err error
	client.Connection, err = tls.Dial("tcp", client.Address, client.TlsConfig)

	if err != nil {
		fmt.Println(err)
		return &RriClient{}, err
	}

	return client, nil
}

func (client *RriClient) Close() error {
	return client.Connection.Close()
}

func (client *RriClient) Login() error {
	loginMsg := &messagetype.LoginMessage{
		Action:   "LOGIN",
		Version:  "2.0",
		Username: client.Username,
		Password: client.Password,
	}
	err := client.send(loginMsg)
	if err != nil {
		return err
	}

	return nil
}

func (client *RriClient) Logout() error {
	logoutMsg := &messagetype.LogoutMessage{
		Action:  "LOGOUT",
		Version: "2.0",
	}
	err := client.send(logoutMsg)
	if err != nil {
		return err
	}
	return nil
}

func (client *RriClient) Read() error {
	for {

		recvBuf := make([]byte, 1024)

		n, err := client.Connection.Read(recvBuf[:]) // recv data
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("read timeout:", err)
				// time out
			}
		}

		if n > 0 {

			fmt.Println(string(recvBuf))
		}
	}

}

func (client *RriClient) SendOrder(order interface{}) error {

	order_data, _ := order.(map[string]interface{})
	switch order_data["ACTION"] {
	case "CREATE":
		fallthrough
	case "UPDATE":
		fallthrough
	case "DELETE":
		order_data["CTID"] = uuid.NewV4().String()
	}

	return client.send(order_data)
}

func (client *RriClient) send(msg interface{}) error {
	yamlmsg, _ := yaml.Marshal(msg)

	// we need to strip away the quotes around the version string.
	formated_yamlmsg := strings.Replace(string(yamlmsg), "\"", "", -1)
	len := len(formated_yamlmsg)

	header := make([]byte, 4)
	databuf := new(bytes.Buffer)

	binary.BigEndian.PutUint32(header, uint32(len))
	var data = []interface{}{
		header,
		[]byte(formated_yamlmsg),
	}
	for _, v := range data {
		err := binary.Write(databuf, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}

	_, err := client.Connection.Write(databuf.Bytes())
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
