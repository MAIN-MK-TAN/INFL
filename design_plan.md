# MK-TAN_INFL  
  
  
> Languages: C++, Go, Rust, and Python    
  
---  
  
## 00. Purpose  
  
This project outlines the structure of a modular, multi-platform control interface for agent-based communication. It defines a cleanroom, non-operational representation of task dispatch, agent tracking, and transport layering strategies common in distributed C2 infrastructure.  

---  
  
## 01. Top-Level Structure  


```
MK-TAN_INFL/  
├── controller/                # Operator-side core (Rust)  
│   ├── interface.rs           # CLI and control logic  
│   ├── agents.rs              # Agent memory registry  
│   ├── tasking.rs             # Task engine  
│   └── transport.rs           # Transport abstraction stubs  
│  
├── agent/                     # Cross-platform implant stubs (C++ / Rust)  
│   ├── beacon.cpp             # Beacon logic  
│   ├── exec.cpp               # Task execution (windows-specific. Helps with anti-detection)  
│   ├── shutdown.cpp           # Kill switch, shutdown triggers  
│   └── transport.cpp          # Pluggable transport layer  
│  
├── protocol/                  # Message schema definitions  
│   ├── agent.proto            # Protobuf schema (primary)  
│   └── fallback.json          # Legacy/static format fallback  
│  
├── proxy/                     # Optional redirector node (Go)  
│   └── main.go  
│  
├── interface/                 # Experimental operator TUI (Python)  
│   └── tui.py  
│  
├── tools/                     # Mutator and entropy seeding tools  
│   ├── mutator.py  
│   └── seeder.py  
│  
├── persistence/               # Optional persistence modules (manual opt-in)  
│   ├── WIN32_regkey.cpp  
│   ├── ELF32_service_systemctl.cpp  
│   ├── WINx8664_regkey.cpp
│   ├── ELFx8664_service_systemctl
│   └── ELFUNIV_pers.bash
│
├── encryption                 # Encryption logic
│   ├── server.rs
│   ├── WIN32.cpp
│   ├── ELF32.cpp
│   └── EXPIRIM_UNIV.py
│
│
├── LICENSE.md                 # Ethical/legal control license  
├── DESIGN_PLAN.md             # Project structural blueprint  
└── README.md                  # Maybe read it
```
---  


## 02. Architecture Philosophy


| Attribute       | Design Rule |
|----------------|-------------|
| Storage        | Memory-only by default. No local writes. |
| Transport      | Abstracted. Implementations are runtime-swappable. |
| Agent Identity | Ephemeral. UUID regenerated per deployment. |
| Result Handling| Padded and transient. No on-disk logs. |
| Code Split     | Controller/implant separation enforced across languages. |
| Deployment     | Build-time flags control behavioral profile. |


---


## 03. Agent Composition


- Agents do not persist unless toggled manually
- All instructions are short-lived and time-expiring
- Memory-resident execution is assumed; disk is treated as hostile
- Beacon intervals, if defined, are jittered and seeded
- Header structure is mutable and randomized
- Execution artifacts are not retained. Shell memory cleared post-task

---


## 04. Transport Model


Transport logic is decoupled from core operation logic.


Planned modules include:
- HTTP(S) over randomized User-Agent and Accept headers
- DNS (TXT polling and domain sharding)
- QUIC (stateless key exchange mode)
- TLS-wrapped POST chains
- Passive mode via hosted redirector
- Transport behaviors defined per-build. No hardcoded protocols.

---


## 05. Operational Profiles


Build flags define behavior classes. Examples include:


| Profile        | Effects |
|----------------|---------|
| `offline`      | Disables all network interfaces |
| `burnable`     | Agent self-destructs after N tasks |
| `mutate`       | Randomizes binary identifiers, task seeds |
| `persist`      | Enables platform-native persistence modules |
 
