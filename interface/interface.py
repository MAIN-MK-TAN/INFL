import json,socket,struct
id="1337";ip="127.0.0.1";pt=9001
def sndjsn(c,d):
 try:e=json.dumps(d).encode();c.sendall(struct.pack(">I",len(e))+e)
 except Exception as ex:print(f"[sndjsn err]{ex}");raise
def rcvjsn(c):
 try:
  h=c.recv(4)
  if len(h)<4:raise ConnectionError("No header/conn lost")
  l=struct.unpack(">I",h)[0];b=b""
  while len(b)<l:
   p=c.recv(l-len(b))
   if not p:raise ConnectionError("Mid-transfer drop")
   b+=p
  return json.loads(b.decode())
 except Exception as ex:print(f"[rcvjsn err]{ex}");raise
def identify(c):
 try:sndjsn(c,{"controller_id":int(id)});print(f"ID {id}→{ip}:{pt}")
 except Exception as ex:print(f"[identify err]{ex}");raise
def conn():
 try:
  c=socket.create_connection((ip,pt));identify(c)
  while 1:
   try:
    x=input(">>> ").strip()
    if not x:continue
    sndjsn(c,{"cmd":x})
    r=rcvjsn(c)
    print(r.get("output","[no output]"))
   except (EOFError,KeyboardInterrupt):print("\n[exit]");break
   except Exception as ex:print(f"[loop err]{ex}");break
 except Exception as ex:print(f"[conn err]{ex}")
if __name__=="__main__":conn()
