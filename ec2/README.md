# Experimental Script Tutorial

The tutorial describes how to deploy our experiment on EC2 of AWS. 
Part of the codes is borrowed from the codebase of HoneyBadgerBFT: https://github.com/amiller/HoneyBadgerBFT

## Prepare
1. Install python 3.8.10 and ipython 8.14.0
2. Pip install "boto" and "fabric" to use EC2. We use the following fabric versions:
Fabric3 1.14; Paramiko 2.7.2
3. Amazon account --> My Security Credentials --> Access keys --> Create New Access Key --> Show Access Key --> store the keys (Access key ID and Secret access key).
4. Generate key pair for SSH and connect regions.
+ ssh-keygen -m PEM, an rsa key pair is generated in the "~/.ssh/" directory. The default names are id_rsa and id_rsa.pub. You can change the names as required.
+ Change the private key file to a "pem" format, use "mv id_rsa id_rsa.pem".
+ Change the permission for the key pair by "chmod 777 id_rsa.pem"
+ Enter "EC2" of AWS, and import the pub key for the region you will access. 
EC2 --> network and security --> key pair --> import key pair.
5. Compile the code to generate "server", "client" and "keygen".
6. Change the config file to choose protocol, IP of server can be changed by script later.
7. For BEAT, generate prf key by "keygen [n] [k]", keys are generated in the "etc/thresprf_key" directory.
8. Enter the directory of the script
9. Copy "server", "client" and "keygen" to the directory, and copy "conf.json" and "thresprf_key" to "etc/".

## Experiment
1. start the script
```
python utility.py xxxx(Access key ID) xxxxx(Secret access key)
```
2. Start EC2 instances in different regions, where [region] is the name of region and [number] represent the
   amount of instance
```
launch_new_instances(region,number)
```
3. Download all instance ip addresses and save them to the local file "hosts", and upload the hosts file to all instances
```
ipAll()
```
4. If there too many servers, for example 40, set the "limits.conf" for each instace. Choose arbitrary instace
```
fab -i ~/.ssh/id_rsa.pem -u ubuntu -P -H [IP] fetchLimit
```
download it to your local etc/, add “* - nofile 65535” in "limits.conf". Then updown the "limits.conf"
to each instance by
```
sLimit()
```
reboot all instance for each [region]
```
reboot_all_instances(region)
```
5. Change the information of server according to the IP in "hosts", where [number] is the amount of total server
```
generateConfJson(number)
```
6. Upload the config file and keys to each instance 
```
sConf()
```
7. Upload the executable file "server" and "client" to each instance
```
sExcute()
```
8. Generate bash file to start all server, where [number] is the amount of total server, "run.sh" will be created
```
generate_server_shell(number)
```
9. Start all server
```
bash run.sh
```
10. SSH to an instance to start the client to start the protocol
```
sudo ssh -i ~/.ssh/xxx.pem ubuntu@ip
cd excute/
./client 100 1 [batchsize]
```

## Compute result
1. All the results are store the log files of each server, down them to local "var/log/"
```
fetchLog()
```
2. Compute throughput, compute the average of all results, where [n] is the amount of servers, [date] is the date of the log. Since there are some
results too big or two small, we will choose [filter=True] to sort the results and discard the first fourth and last fourth to compute average.
```
computeQuorumThroughput(n=, date='', filter=)
```
3. Compute latency
```
computeLatency(n=, date='', filter=)
```
4. Compute latency vs throughput
```
latency_throughput(n=, date='', filter=)
```

5. Clear the log files if you want to deploy next experiment
```
clearLog()
```
6. Kill all servers if you finish an experiment
```
killServer()
```
7. Terminate all EC2 instances of all regions when finish all experiments to avoid cost.
```
terminate_all_instances(region)
```
