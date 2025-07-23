###/ MODULES \###
import json,socket,struct
###/ VARS \###
id="1337"
ip="127.0.0.1"
pt=9001
###/ CODE \###
def sndjsn(conc,dta):
 encoded=json.dumps(dta).encode()
 conc.sendall(struct.pack(">I",len(encoded))+encoded)
def rcvjsn(conc):
 l=struct.unpack(">I",conc.recv(4))[0]
 d=b""
 while len(d)<l:
  c=conc.recv(l-len(d))
  if not c:raise ConnectionError("Disconnected during recv")
  d+=c
 return json.loads(d.decode())
def identify(conc,id):
 sndjsn(conc,{"controller_id":int(id)})
 print(f"Identified to {ip}:{pt} with id {id}")
def conn2srvr(ip,pt):
 conc=socket.create_connection((ip,pt))
 identify(conc,id)
 while 1:
  x=input(">>> ")
  if not x:continue
  sndjsn(conc,{"cmd":x})
  print(rcvjsn(conc).get("output",""))
conn2srvr(ip,pt)
