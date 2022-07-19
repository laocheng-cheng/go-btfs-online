# Status-server

These are steps to manually build a local status-server.

## Fetch source code
We need to switch Makefile after fetch source code.
```bash
git clone https://github.com/TRON-US/status-server.git
cd status-server
mv Makefile Makefile.cloud
mv Makefile.local Makefile
```

## Install prerequisites  
Install GoLang on the machine.  
Install PostgreSQL on Mac and create the status-server db run:
```bash
make install
```

## Config ENV for local instance of status-server    
Replace username, password, host, db name for PostgreSQL database setting.
```bash
export DB_USER=<username>:<password>
export DB_HOSTNAME=localhost
export DB_NAME=db_status
```

## Build
```bash
make build
```
status-server binary will be built into ./bin directory.

## Start the status server at the backend
Issue the following command:
```bash
cd bin
nohup status-server </dev/null >/dev/null 2>&1 &
```
If copy and run status-server in another directory, make sure to copy the IP2LOCATION DB along with status-server binary.

## Updating the IP2LOCATION DB
The binary is located along with the status-server binary. When replacing/updating it, make sure to update it there as well. Run:
```bash
which status-server
```
to get the location of the binary.

# Migration Tools
If run status-server in directory other than default ./bin, make sure to put status-server, migration binary and migrations dir in same directory.
```bash
$ ls
IP2LOCATION-LITE-DB5.BIN	migration	migrations	status-server
```

## Deploy to dev via kubernetes

```make -f Makefile.local deploy```

It will complie a version tag with latest git hash and timestamp appended, and use it to build and push to dev env.(eg. master-c61c91ca106dc1ebf69af04af0c33697d58f9149-1578002075-) The version tag will also be clean up after the deploy.

## Auto migrations

### Build Tool
```bash
make tools
```
### Init (first time run)
```bash
make db_init
```
### Up
```bash
make db_up
```
### Down
```bash
make db_down
```

## To developer - Format and clean up code
Please run this or make sure your code is formated before submiting your coding.
```bash
make fmt
```
