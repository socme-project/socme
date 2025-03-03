# Infrastructure - DEV

Réseau: 10.8.176.0/20
DNS: 10.8.17.5
Passerelle: 10.8.176.1

**PROXMOX**:

- usr: root
- psswd: 9#@5H%v2B$a3
- url: <https://10.8.176.109:8006>

**PFSENSE**:

- wan: 10.8.178.20/20 (VLAN 9, vmbr0 -> école)
- lan: 192.168.1.1/24 (VLAN 13)
- usr: admin
- psswd: admin

**KALI**:

- lan: 192.168.1.100 (VLAN 13)
- usr: kali
- pass: kali

**WINDOWS**:
- usr: windows
- pass: windows
