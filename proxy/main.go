package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	controllerConn net.Conn
	controllerLock sync.Mutex
)

func readJSON(r *bufio.Reader) (map[string]interface{}, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
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

// isConnDead checks if conn is dead by peeking 1 byte without consuming it.
func isConnDead(conn net.Conn) bool {
	if conn == nil {
		return true
	}
	br := bufio.NewReader(conn)
	_ = conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
	_, err := br.Peek(1)
	_ = conn.SetReadDeadline(time.Time{})
	if err == io.EOF {
		return true
	}
	return err != nil
}

func replaceController(conn net.Conn) bool {
	// Check liveness before locking
	if controllerConn != nil && !isConnDead(controllerConn) {
		log.Println("[!] Existing controller still alive — rejecting new connection")
		return false
	}

	controllerLock.Lock()
	defer controllerLock.Unlock()

	if controllerConn != nil {
		controllerConn.Close()
		log.Println("[*] Previous controller connection closed")
	}
	controllerConn = conn
	return true
}

func handleController(conn net.Conn, agentConn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	agentReader := bufio.NewReader(agentConn)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	msg, err := readJSON(br)
	conn.SetReadDeadline(time.Time{}) // clear timeout
	if err != nil {
		log.Println("[X] Handshake read failed:", err)
		return
	}

	rawID, ok := msg["controller_id"]
	idFloat, valid := rawID.(float64)
	if !ok || !valid || int(idFloat) != 1337 {
		log.Println("[!] Invalid or missing controller_id — dropping")
		return
	}

	if !replaceController(conn) {
		return
	}

	log.Printf("[+] Controller 1337 connected from %s", conn.RemoteAddr())

	for {
		msg, err := readJSON(br)
		if err != nil {
			log.Println("[X] Controller read error:", err)
			break
		}

		if isConnDead(agentConn) {
			log.Println("[X] Agent connection is dead — aborting")
			break
		}

		if err := writeJSON(agentConn, msg); err != nil {
			log.Println("[X] Write to agent failed:", err)
			break
		}

		resp, err := readJSON(agentReader)
		if err != nil {
			log.Println("[X] Read from agent failed:", err)
			break
		}

		if err := writeJSON(conn, resp); err != nil {
			log.Println("[X] Return to controller failed:", err)
			break
		}
	}

	controllerLock.Lock()
	if controllerConn == conn {
		controllerConn = nil
		log.Printf("[*] Controller %s disconnected", conn.RemoteAddr())
	}
	controllerLock.Unlock()
}

func main() {
	agentListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("[FATAL] Failed to bind agent listener:", err)
	}

	controllerListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("[FATAL] Failed to bind controller listener:", err)
	}

	log.Println("[*] Awaiting agent on :9000...")
	agentConn, err := agentListener.Accept()
	if err != nil {
		log.Fatal("[FATAL] Agent accept failed:", err)
	}
	log.Printf("[+] Agent connected from %s", agentConn.RemoteAddr())

	log.Println("[*] Awaiting controller connections on :9001...")
	for {
		conn, err := controllerListener.Accept()
		if err != nil {
			log.Println("[!] Controller accept error:", err)
			continue
		}
		go handleController(conn, agentConn)
	}
}
