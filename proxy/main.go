package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
)

type ControllerConn struct {
	Conn net.Conn
	ID   int
}

var (
	controllerSessions = make(map[int]net.Conn)
	controllerLock     sync.Mutex
)

// Reads a JSON message prefixed with 4-byte big-endian length
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

// Writes a JSON message with a 4-byte big-endian length prefix
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
	if !ok {
		log.Println("[X] Invalid controller_id")
		return
	}
	controllerID := int(idFloat)

	controllerLock.Lock()
	if _, exists := controllerSessions[controllerID]; exists {
		controllerLock.Unlock()
		log.Printf("[!] Duplicate controller ID %d rejected\n", controllerID)
		return
	}
	controllerSessions[controllerID] = conn
	controllerLock.Unlock()

	log.Printf("[+] Controller %d connected\n", controllerID)

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
	delete(controllerSessions, controllerID)
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
