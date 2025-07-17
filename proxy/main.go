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

func readJSON(conn net.Conn) (map[string]interface{}, error) {
    lenBuf := make([]byte, 4)
    if _, err := io.ReadFull(conn, lenBuf); err != nil {
        return nil, err
    }
    msgLen := binary.BigEndian.Uint32(lenBuf)
    msgBuf := make([]byte, msgLen)
    if _, err := io.ReadFull(conn, msgBuf); err != nil {
        return nil, err
    }

    var msg map[string]interface{}
    err := json.Unmarshal(msgBuf, &msg)
    return msg, err
}

func writeJSON(conn net.Conn, obj interface{}) error {
    data, err := json.Marshal(obj)
    if err != nil {
        return err
    }
    length := uint32(len(data))
    lenBuf := make([]byte, 4)
    binary.BigEndian.PutUint32(lenBuf, length)
    _, err = conn.Write(append(lenBuf, data...))
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

        err = writeJSON(agentConn, msg)
        if err != nil {
            log.Println("[X] Failed to forward to agent")
            break
        }

        response, err := readJSON(agentConn)
        if err != nil {
            log.Println("[X] Failed to receive response from agent")
            break
        }

        err = writeJSON(conn, response)
        if err != nil {
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
    log.Println("[*] Waiting for controller on :9001")

    for {
        conn, err := controllerListener.Accept()
        if err != nil {
            continue
        }
        go handleController(conn, agentConn)
    }
}
