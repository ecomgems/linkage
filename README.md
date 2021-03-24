# linkage

Linkage is an SSH TCP Port Forwarding app to open multiple tunnels simultaneously written totally in GoLang.

## Run using Config file

```bash
linkage -c example.yml
```

The config file example:

```yaml
servers:
  - host: remote.server1.tld
    port: 22
    user: dev
    key_file: ~/.ssh/id_rsa.pub
    tunnels:
      - remote_port: 80
        remote_host: web_app
        local_port: 80 
      - remote_port: 443
        remote_host: web_app
        local_port: 443  
          
  - host: remote.server2.tld
    port: 22
    user: dev
    password: <password>
    tunnels:
      - remote_port: 9200
        remote_host: elastic
        local_port: 9200
```
