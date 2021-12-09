# cs655-project

### Setup
This system can have up to 100 servers, and an rspec is provided to create the network topology.

In each server, run

```bash
sudo -s
git clone https://github.com/cshi0/cs655-project.git
cd cs655-project
sh install.sh
```

### Run
To run the program, change directory to where cs655-project is, and run

```bash
./cracker
```

The http server will be on, and requests can be sent to any server in the network

### Script
In order to get the data used in the report, the python script should be ran.

```bash
python3 hash_gen.py
```

It will generate 10 random 5-character strings, and send 10 crackTask requests to the server running on localhost. After seeing no tasks are performed in the log printed, the metrics can be retrieved by using an API.

```bash
curl -X GET http://localhost:8080/metrics
```

### Publicly router
It is recommended to run requests on localhost which makes it easier to use the script and makes it more stable, but all servers are publicly routable. For example, node-1 in my slice has an IP of 143.215.216.204, which can be called.