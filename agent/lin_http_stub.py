import requests
import json
import subprocess

URL = "http://127.0.0.1:8080/agent"

def main():
    while True:
        try:
            cmd_obj = requests.post(URL, json={"init": True}).json()
            cmd = cmd_obj.get("cmd", "")
            if not cmd:
                continue
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
            out = result.stdout + result.stderr
            requests.post(URL, json={"output": out})
        except Exception:
            pass

if __name__ == "__main__":
    main()
