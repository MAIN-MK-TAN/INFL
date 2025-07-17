use std::io::{self, Read, Write};
use std::net::TcpStream;
use serde_json::{json, Value};

const CONTROLLER_ID: u32 = 1337;
const PROXY_HOST: &str = "127.0.0.1";
const PROXY_PORT: u16 = 9001;

fn send_json(mut stream: &TcpStream, obj: &Value) -> io::Result<()> {
    let data = serde_json::to_vec(obj)?;
    let len = (data.len() as u32).to_be_bytes();
    stream.write_all(&len)?;
    stream.write_all(&data)?;
    Ok(())
}

fn recv_json(mut stream: &TcpStream) -> io::Result<Option<Value>> {
    let mut len_buf = [0u8; 4];
    if stream.read_exact(&mut len_buf).is_err() {
        return Ok(None);
    }

    let len = u32::from_be_bytes(len_buf) as usize;
    let mut data_buf = vec![0u8; len];
    stream.read_exact(&mut data_buf)?;

    let obj: Value = serde_json::from_slice(&data_buf)?;
    Ok(Some(obj))
}

fn main() -> io::Result<()> {
    let addr = format!("{}:{}", PROXY_HOST, PROXY_PORT);
    let stream = TcpStream::connect(addr)?;

    send_json(&stream, &json!({ "controller_id": CONTROLLER_ID }))?;

    let stdin = io::stdin();
    let mut input = String::new();

    loop {
        print!("cmd > ");
        io::stdout().flush()?;
        input.clear();
        if stdin.read_line(&mut input)? == 0 {
            break;
        }

        let cmd = input.trim();
        if cmd.is_empty() {
            continue;
        }

        send_json(&stream, &json!({ "cmd": cmd }))?;

        match recv_json(&stream)? {
            Some(resp) => {
                let output = resp.get("output").and_then(|v| v.as_str()).unwrap_or("[X] No output");
                println!("{}", output);
            }
            None => {
                println!("[X] No response from proxy.");
                break;
            }
        }
    }

    Ok(())
}
