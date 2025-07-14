# MK-TAN_INFL

└── README.md                  # Maybe read it
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


Default builds operate in isolated simulation mode.
 
