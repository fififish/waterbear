# WaterBear Reposity

Asynchronous fault-tolerant protocols for the following paper: 

Zhang, Haibin, et al. "WaterBear: Asynchronous BFT with Information-Theoretic Security and Quantum Security." Cryptology ePrint Archive (2022).

This repository implements five BFT protocols:

+ BEAT-Cobalt (BEAT0)
+ WaterBear-C (using Bracha's ABA and Cubic-RABA)
+ WaterBear-Q (using Bracha's ABA and Quadratic-RABA)
+ WaterBear-QS-C (using AVID and  Cubic-RABA).
+ WaterBear-QC-Q (using AVID and  Quadratic-RABA).

Different RBC modules are implemented under src/broadcast, and differen ABA modules are impleented under src/aba. See the README files under each folder for details.

### Configuration

Configuration is under etc/conf.json

Change "consensus" to switch between the protocols. See note.txt for details. 

### Installation && How to run the code

#### Install dependencies 

+ enter the directory and run the following commands
```
make all 
```

+ If you only need to update reomte github entries, run 
```
make install 
```

+ If you only need to build again, run 
```
make build
```

##### Launch the code

+ For all the servers, run the command below to start the servers. The [id] is configured in conf.json
```
        server [id]
```

+ Start a client
```
        client [id] 1 [batch]
```

- [id] could be anything. [batch] is the batch size. 

+ The default message is "abcdefg", one could change it by either editing src/main/client.go and compiling again, or run the following command 
```
        client [id] [batch] [msg]
```

+ To run BEAT-Cobalt, generate keys for threshold PRF first
```
        keygen [n] [k]
```
- [k] is set to f+1

### How to deploy the Amazon EC2 experiment

+ Scripts are not included in this repository. 