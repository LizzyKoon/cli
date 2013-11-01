package apps

import (
	"appsdeck/cli/api"
	"appsdeck/cli/config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"syscall"
)

func Run(app string, command []string) error {
	res, err := api.Run(app, command)
	if err != nil {
		return err
	}
	runStruct := make(map[string]interface{})
	ReadJson(res.Body, &runStruct)

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Not authorized")
	}

	if _, ok := runStruct["attach"]; !ok {
		return fmt.Errorf("Unexpected answer from server")
	}

	res, socket, err := connectToRunServer(runStruct["attach"].(string))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Fail to attach: %s", res.Status)
	}

	sttyArgs := []string{"stty", "-echo", "raw"}
	fd := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	_, err = syscall.ForkExec("/bin/stty", sttyArgs, &syscall.ProcAttr{Dir: "", Files: fd})
	if err != nil {
		return err
	}

	go io.Copy(socket, os.Stdout)
	io.Copy(os.Stdin, socket)
	socket.Close()

	sttyArgs = []string{"stty", "echo", "cooked"}
	fd = []uintptr{os.Stdout.Fd(), os.Stderr.Fd()}
	syscall.ForkExec("/bin/stty", sttyArgs, &syscall.ProcAttr{Dir: "", Files: fd})
	return nil
}

func connectToRunServer(rawUrl string) (*http.Response, net.Conn, error) {
	params := map[string]string{
		"user_email": api.AuthEmail,
		"user_token": api.AuthToken,
	}
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return nil, nil, err
	}
	paramsReader := bytes.NewReader(paramsJson)

	req, err := http.NewRequest("POST", rawUrl, paramsReader)
	if err != nil {
		return nil, nil, err
	}

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, nil, err
	}

	dial, err := net.Dial("tcp", url.Host)
	if err != nil {
		return nil, nil, err
	}

	tls_conn := tls.Client(dial, config.TlsConfig)
	conn := httputil.NewClientConn(tls_conn, nil)

	res, err := conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		return nil, nil, err
	}

	connection, _ := conn.Hijack()
	return res, connection, nil
}
