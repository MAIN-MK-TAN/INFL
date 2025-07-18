import socket
import json
import struct

CONTROLLER_ID = 1337
PROXY_HOST = "127.0.0.1"
PROXY_PORT = 9001

def send_json(sock, obj):
    data = json.dumps(obj).encode()
    sock.sendall(struct.pack(">I", len(data)) + data)

def recv_json(sock):
    length_bytes = sock.recv(4)
    if not length_bytes:
        return None
    length = struct.unpack(">I", length_bytes)[0]
    data = sock.recv(length)
    return json.loads(data.decode())

def main():
    s = socket.create_connection((PROXY_HOST, PROXY_PORT))
    send_json(s, { "controller_id": CONTROLLER_ID })

    try:
        while True:
            cmd = input("cmd > ").strip()
            if not cmd:
                continue
            send_json(s, { "cmd": cmd })
            response = recv_json(s)
            if not response:
                print("[X] No response from proxy.")
                break
            print(response.get("output", "[X] No output"))
    except KeyboardInterrupt:
        s.close()

if __name__ == "__main__":
    main()
