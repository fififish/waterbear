# WaterBear Repository

Asynchronous fault-tolerant protocols for the following paper: 

<<<<<<< HEAD
Haibin Zhang, Sisi Duan, Boxin Zhao, and Liehuang Zhu. "WaterBear: Practical Asynchronous BFT Matching Security Guarantees of Partially Synchronous BFT." Usenix Security 2023, to appear. 

Epring version: https://eprint.iacr.org/2022/021
=======
Haibin Zhang, Sisi Duan, Boxin Zhao, and Liehuang Zhu. "WaterBear: Practical Asynchronous BFT Matching Security Guarantees of Partially Synchronous BFT." Cryptology ePrint Archive (2022).
https://eprint.iacr.org/2022/021
>>>>>>> 4b4f769195a72c9d6e59047cd5ce8f37c7dff923

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

<<<<<<< HEAD
+ Scripts are included in ec2 folder. 
=======
+ Scripts are not included in this repository. 
>>>>>>> 4b4f769195a72c9d6e59047cd5ce8f37c7dff923
