# WaterBear Repository

Asynchronous fault-tolerant protocols for the following paper: 

Haibin Zhang, Sisi Duan, Boxin Zhao, and Liehuang Zhu. "WaterBear: Practical Asynchronous BFT Matching Security Guarantees of Partially Synchronous BFT." Usenix Security 2023. 

Eprint version: https://eprint.iacr.org/2022/021

This repository implements five BFT protocols:

+ BEAT-Cobalt (BEAT0)
+ WaterBear-C (using Bracha's ABA and Cubic-RABA)
+ WaterBear-Q (using Bracha's ABA and Quadratic-RABA)
+ WaterBear-QS-C (using CT RBC and  Cubic-RABA).
+ WaterBear-QS-Q (using CT RBC and  Quadratic-RABA).

Different RBC modules are implemented under src/broadcast, and different ABA modules are implemented under src/aba. See the README files under each folder for details.

### Configuration

The codebase can only be compiled for x86 machines for now. The src/cryptolib/word/word_amd64.s needs to be changed to adapt to arm machines. 

Configuration is under etc/conf.json

Change "consensus" to switch between the protocols. See note.txt for details. 

### Experimental environment && Dependency

+ system: ubuntu 20.04
+ golang: go1.15.14 linux/amd64
+ google.golang.org/grpc: 1.35.0
+ github.com/cbergoon/merkletree: 0.1.0
+ golang.org/x/net: 
+ golang.org/x/text: 
+ golang.org/x/crypto: 
+ golang.org/x/sys: 
+ google.golang.org/genproto: 
+ github.com/klauspost/reedsolomon: 
+ github.com/klauspost/cpuid:
+ github.com/golang/protobuf: 


### Installation && How to run the code

#### Install dependencies 

+ enter the directory and run the following commands
```
export GOPATH=$PWD
export GOBIN=$PWD/bin
export GO111MODULE=off
```

+ If you need to update grpc, net, genproto and so on, run
```
make go
```

+ If you only need to update remote github entries, run 
```
make install 
```

+ If you only need to build again, run 
```
make build
```

##### Launch the code

+ Modify the configuration file  "etc/conf.json" to choose which protocol to execute, and modify the IP addresses 
and port numbers of all servers. Details about the protocols are included in "$etc/node.txt$". The “id” of each server should be unique. 
By default, we use monotonically increasing ids, 0, 1, 2, ....
+ To run BEAT-Cobalt, generate keys for threshold PRF first by running the following command:  
```
        keygen [n] [k]
```
Here, "n" is the number of servers, and "k" is the threshold to generate the common coin. We set up "n=3f+1" and "k" to "f+1" for 
most of our experiments.
+ For all the servers, run the command below to start the servers:
```
         server [id]
```
Here, "id" is configured in conf.json and is different at each server.
+ Start a client to send transaction to start the protocol by running the following command:
```
        client [id] 1 [batch] [msg]
```
Here, "id" is the identifier of the client. We do not require the client to be registered. One can use any id that is unique,
e.g., 1000. "b" is the batch size. "msg" can be any message. One can ignore the "msg" field and a default message is included in the codebase.

##### Quick experiment locally
+ Open four terminals and start server in each terminal (start one server in one terminal).
```
        server 0
        server 1
        server 2
        server 3
```
+ Open one terminal to start client and send transactions.
```
        client 100 1 10000
```
All server terminal will print text like " *****epoch ends", which represents the success of the epoch.
One can repeat the operation of client after the epoch end.

### How to deploy the Amazon EC2 experiment

+ Scripts are included in ec2 folder. 
