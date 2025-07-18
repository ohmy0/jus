# JUS — Just UID/GID Set

A minimal `sudo`/`doas` alternative written in **Go** (only **252 lines!**).  
Uses **PAM** for auth.

---  

## 🔧 Build

```sh
# Clone repository
git clone https://github.com/ohmy0/jus
cd ./jus
# Please install golang before building
go mod tidy
go build jus
# Move bin file ( run as root )
mv ./jus /usr/local/bin
# Set privileges ( run as root )
chown root:root /usr/local/bin/jus 
chmod u+s /usr/local/bin/jus
# Create config file ( run as root )
touch /etc/jus.toml
chmod 644 /etc/jus.toml
# Finish.

```

---

## ⚙️ Config
```toml  
[[permit]] # There can be many such constructs

user="youruser" # The user under which this configuration will be applied

as="root" # Under whose identity the command will be executed

paths=["/usr/bin","/bin"] # Paths in which commands will be searched ( Optional, std paths = /bin /sbin /usr/bin /usr/sbin /usr/local/bin /usr/local/sbin ) 

keepenv=true # saves variables PATH, TERM, LANG
```

---

### 🚀 Usage
```sh
jus command args
```
