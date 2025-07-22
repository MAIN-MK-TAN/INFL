package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
)

var (
	controllerConn net.Conn
	controllerLock sync.Mutex
)

func readJSON(conn net.Conn) (map[string]interface{}, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	payload := make([]byte, length)

	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}

	var msg map[string]interface{}
	err := json.Unmarshal(payload, &msg)
	return msg, err
}

func writeJSON(conn net.Conn, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	length := uint32(len(data))
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, length)

	_, err = conn.Write(append(lengthBuf, data...))
	return err
}

func handleController(conn net.Conn, agentConn net.Conn) {
	defer conn.Close()

	msg, err := readJSON(conn)
	if err != nil {
		log.Println("[X] Controller handshake failed:", err)
		return
	}

	idFloat, ok := msg["controller_id"].(float64)
	if !ok || int(idFloat) != 13 {
		log.Println("[!] Invalid controller_id, rejected")
		return
	}

	controllerLock.Lock()
	if controllerConn != nil {
		// Check if existing conn is still alive
		one := make([]byte, 1)
		connCheck := make(chan bool, 1)
		go func(c net.Conn) {
			_, err := c.Read(one)
			connCheck <- err != nil
		}(controllerConn)

		if <-connCheck == false {
			log.Println("[!] Controller 1337 already connected and alive, rejecting new")
			controllerLock.Unlock()
			return
		}

		controllerConn.Close()
		controllerConn = nil
		log.Println("[*] Previous controller appears dead, accepting new")
	}
	controllerConn = conn
	controllerLock.Unlock()

	log.Println("[+] Controller 1337 connected")

	for {
		msg, err := readJSON(conn)
		if err != nil {
			log.Println("[X] Controller connection closed")
			break
		}

		if err := writeJSON(agentConn, msg); err != nil {
			log.Println("[X] Failed to forward to agent")
			break
		}

		response, err := readJSON(agentConn)
		if err != nil {
			log.Println("[X] Failed to receive response from agent")
			break
		}

		if err := writeJSON(conn, response); err != nil {
			log.Println("[X] Failed to return response to controller")
			break
		}
	}

	controllerLock.Lock()
	if controllerConn == conn {
		controllerConn = nil
	}
	controllerLock.Unlock()
}

func main() {
	agentListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	controllerListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[*] Waiting for agent on :9000")
	agentConn, err := agentListener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[+] Agent connected")
	log.Println("[*] Waiting for controllers on :9001")

	for {
		conn, err := controllerListener.Accept()
		if err != nil {
			continue
		}
		go handleController(conn, agentConn)
	}
}
