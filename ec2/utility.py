import argparse
import boto.ec2
import sys, os
import time
import json

your_key_path = "~/.ssh/"
your_key_name = "usenix_rsa"


if not boto.config.has_section('ec2'):
    boto.config.add_section('ec2')
    boto.config.setbool('ec2', 'use-sigv4', True)

secgroups = {
    'us-east-1': 'sg-3babf620',  # virginia--no
    'us-east-2': 'sg-cc461c85',  # Ohio--no
    'us-west-1': 'sg-2e1e2464',  # california--ok
    'us-west-2': 'sg-09473f6eb7cf6feac',  # Oregon--ok
    'eu-west-1': 'sg-be4341ed',  # Ireland--ok
    'sa-east-1': 'sg-05488e7aa31e2ef00',  # South America--
    'ap-southeast-1': 'sg-f6f794b9',  # Singapore--
    'ap-southeast-2': 'sg-7fadf33f',  # Sydney--
    'ap-northeast-1': 'sg-c7566983',  # Tokyo--
    'ca-central-1': 'sg-6e7ec313',  # Canada--
    'eu-north-1': 'sg-9fdc1af9',  #
    'eu-west-3': 'sg-8d7f1cf4',  # Paris--no
    'eu-west-2': 'sg-9fafe4e0',  # London--no
    'ap-south-1': 'sg-05f5597a',  # Asia Pacific (Mumbai)
    'ap-northeast-2': 'sg-ad5a2ad2',  # Seoul
    'ap-northeast-3': 'sg-ad5a2ad2'  # Osaka
    #    'eu-central-1':'sg-2bfe9342'  # somehow this group does not work
}
regions = sorted(secgroups.keys())[::-1]

NameFilter = 'USENIX'


def getAddrFromEC2Summary(s):
    return [
        x.split('ec2.')[-1] for x in s.replace(
            '.compute.amazonaws.com', ''
        ).replace(
            '.us-west-1', ''  # Later we need to add more such lines
        ).replace(
            '-', '.'
        ).strip().split('\n')]


def get_ec2_instances_ip(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    if ec2_conn:
        result = []
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        for reservation in reservations:
            if reservation:
                for ins in reservation.instances:
                    if ins.public_dns_name:
                        currentIP = ins.public_dns_name.split('.')[0][4:].replace('-', '.')
                        result.append(currentIP)
                        print(currentIP)
        return result
    else:
        #print('Region failed', region)
        return None


def get_ec2_instances_id(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    if ec2_conn:
        result = []
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        for reservation in reservations:
            for ins in reservation.instances:
                print
                ins.id
                result.append(ins.id)
        return result
    else:
        print
        'Region failed', region
        return None


def stop_all_instances(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    idList = []
    if ec2_conn:
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        for reservation in reservations:
            if reservation:
                for ins in reservation.instances:
                    idList.append(ins.id)
        ec2_conn.stop_instances(instance_ids=idList)


def terminate_all_instances(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    idList = []
    if ec2_conn:
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        for reservation in reservations:
            if reservation:
                for ins in reservation.instances:
                    idList.append(ins.id)
        ec2_conn.terminate_instances(instance_ids=idList)


def launch_new_instances(region, number):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    dev_sda1 = boto.ec2.blockdevicemapping.EBSBlockDeviceType(delete_on_termination=True)
    dev_sda1.size = 8  # size in Gigabytes
    dev_sda1.delete_on_termination = True
    bdm = boto.ec2.blockdevicemapping.BlockDeviceMapping()
    bdm['/dev/sda1'] = dev_sda1
    img = ec2_conn.get_all_images(filters={'name': 'ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-20220610'})[0].id
    reservation = ec2_conn.run_instances(image_id=img,  # 'ami-df6a8b9b',  # ami-9f91a5f5
                                         min_count=number,
                                         max_count=number,
                                         key_name=your_key_name,
                                         instance_type='m5.xlarge',
                                         security_group_ids=[secgroups[region], ],
                                         block_device_map=bdm)
    for instance in reservation.instances:
        instance.add_tag("Name", NameFilter)
    return reservation


def start_all_instances(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    idList = []
    if ec2_conn:
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        for reservation in reservations:
            for ins in reservation.instances:
                idList.append(ins.id)
        print(idList)
        ec2_conn.start_instances(instance_ids=idList)


def get_all_reservations(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    if ec2_conn:
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        return reservations


def reboot_all_instances(region):
    ec2_conn = boto.ec2.connect_to_region(region,
                                          aws_access_key_id=access_key,
                                          aws_secret_access_key=secret_key)
    idList = []
    if ec2_conn:
        reservations = ec2_conn.get_all_reservations(filters={'tag:Name': NameFilter})
        print(len(reservations))
        for reservation in reservations:
            for ins in reservation.instances:
                idList.append(ins.id)
        ec2_conn.reboot_instances(instance_ids=idList)


def ipAll():
    result = []
    for region in regions:
        result += get_ec2_instances_ip(region) or []
    open('hosts', 'w').write('\n'.join(result))
    callFabFromIPList(result, 'removeHosts')
    callFabFromIPList(result, 'writeHosts')
    return result


def ipRegion(region):
    result = []
    result += get_ec2_instances_ip(region) or []
    open('hosts', 'w').write('\n'.join(result))
    callFabFromIPList(result, 'removeHosts')
    callFabFromIPList(result, 'writeHosts')
    return result


def getIP():
    return [l for l in open('hosts', 'r').read().split('\n') if l]


def getClientIP():
    return [l for l in open('client_hosts', 'r').read().split('\n') if l]


def clientIPHosts(number=0):
    result = getIP()
    ip_len = len(result)
    result = result[ip_len - number:ip_len]
    open('client_hosts', 'w').write('\n'.join(result))


def idAll():
    result = []
    for region in regions:
        result += get_ec2_instances_id(region) or []
    return result


def startAll():
    for region in regions:
        start_all_instances(region)


def stopAll():
    for region in regions:
        stop_all_instances(region)


from subprocess import check_output, Popen, call, PIPE, STDOUT
import fcntl
from threading import Thread
import platform


def callFabFromIPList(l, work):
    print(l, work)
    if platform.system() == 'Darwin':
        print
        Popen(['fab', '-i', '~/.ssh/id_rsa.pem',
               '-u', 'ubuntu', '-H', ','.join(l),  # We rule out the client
               work])
    else:
        print
        'fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s %s' % (your_key_name, ','.join(l), work)
        call('fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s %s' % (your_key_name, ','.join(l), work), shell=True)


def non_block_read(output):
    ''' even in a thread, a normal read with block until the buffer is full '''
    fd = output.fileno()
    fl = fcntl.fcntl(fd, fcntl.F_GETFL)
    fcntl.fcntl(fd, fcntl.F_SETFL, fl | os.O_NONBLOCK)
    try:
        return output.readline()
    except:
        return ''


def monitor(stdout, N, t):
    starting_time = time.time()
    counter = 0
    while True:
        output = non_block_read(stdout).strip()
        print
        output
        if 'synced transactions set' in output:
            counter += 1
            if counter >= N - t:
                break
    ending_time = time.time()
    print
    'Latency from client scope:', ending_time - starting_time


def runProtocol():  # fast-path to run, assuming we already have the files ready
    callFabFromIPList(getIP(), 'runProtocol')


def runProtocolfromClient(client, key, hosts=None):
    if not hosts:
        callFabFromIPList(getIP(), 'runProtocolFromClient:%s,%s' % (client, key))
    else:
        callFabFromIPList(hosts, 'runProtocolFromClient:%s,%s' % (client, key))


def runEC2(Tx, N, t, n):  # run 4 in a row
    for i in range(1, n + 1):
        runProtocolfromClient('"%d %d %d"' % (Tx, N, t), "~/%d_%d_%d.key" % (N, t, i))


# def runServer():
#     index = 0  # server id
#     for ip in getIP():
#         l = [ip]
#         callFabFromIPList(l, 'runServer:%s' % str(index))
#         index += 1


def runServer():
    callFabFromIPList(getIP(), 'runServer')


def runClientShell():
    callFabFromIPList(getClientIP(), 'runClientShell')


def checkLog():
    callFabFromIPList(getIP(), 'checkLog')


def clearLog():
    callFabFromIPList(getIP(), 'clearLog')


def fetchLog():
    callFabFromIPList(getIP(), 'fetchLogs')

def fetchServerLog(sNum = 4):
    callFabFromIPList(getIP()[0:sNum], 'fetchLogs')


def fetchLimit():
    callFabFromIPList(getIP(), 'fetchLimit')


def fetchEvaLog():
    index = 0  # server id
    for ip in getIP():
        l = [ip]
        callFabFromIPList(l, 'fetchEvaLogs:%s' % str(index))
        index += 1


def generateConfJson(number=0):
    json_data = dict()
    # 读取json文件内容,返回字典格式
    with open('etc/conf.json', 'r', encoding='utf8')as fp:
        json_data = json.load(fp)
        json_data['replicas'] = list()
        if number == 0:
            for i, ip in enumerate(getIP()):
                replica = dict()
                replica['id'] = str(i)
                replica['host'] = ip
                replica['port'] = str(11000 + i)
                json_data['replicas'].append(replica)
            fp.close()
        else:
            for i, ip in enumerate(getIP()):
                if i == number:
                    break
                replica = dict()
                replica['id'] = str(i)
                replica['host'] = ip
                replica['port'] = str(11000 + i)
                json_data['replicas'].append(replica)
            fp.close()

    with open('etc/conf.json', 'w+', encoding='utf8')as fp:
        json.dump(json_data, fp)


def generate_server_shell(number=0):
    ip_list = getIP()
    with open('run.sh', 'w+') as fp:
        fp.write("#! /bin/bash\n")
        if number == 0:
            for index, ip in enumerate(ip_list):
                fp.write("fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s runServer:%d&\n" % (your_key_name, ip, index))
        else:
            for i in range(0, number):
                fp.write("fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s runServer:%d&\n" % (your_key_name, ip_list[i], i))
    fp.close()


def generate_client_shell(client_number=0, request_number=1):
    ip_list = getClientIP()
    with open('start_clients.sh', 'w+') as fp:
        fp.write("#! /bin/bash\n")
        fp.write("for ((i=0; i<%d; i++))\n" % client_number)
        fp.write("do\n{\n")
        fp.write("\t./excute/client $[100+$i] 0 %d\n" % request_number)
        fp.write("}&\n")
        fp.write("done\nwait\n")
    fp.close()


def generate_client_shell_local(client_number=0, request_number=1):
    ip_list = getClientIP()
    with open('run_clients.sh', 'w+') as fp:
        fp.write("#! /bin/bash\n")
        for cip in ip_list:
            for i in range(0, client_number):
                fp.write("fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s runClient:%d,0,%d&\n" % (
                your_key_name, cip, 100 + i, request_number))
    fp.close()


def generate_fetchlog_shell(number=0):
    ip_list = getIP()
    with open('fetchlog.sh', 'w+') as fp:
        fp.write("#! /bin/bash\n")
        if number == 0:
            for index, ip in enumerate(ip_list):
                fp.write("fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s fetchEvaLogs:%d&\n" % (your_key_name, ip, index))
        else:
            for i in range(0, number):
                fp.write("fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s fetchEvaLogs:%d&\n" % (your_key_name, ip_list[i], i))
    fp.close()


def check_server_state():
    callFabFromIPList(getIP(), 'checkServerState')


def stopProtocol():
    callFabFromIPList(getIP(), 'stopProtocols')


def callStartProtocolAndMonitorOutput(N, t, l, work='runProtocol'):
    if platform.system() == 'Darwin':
        popen = Popen(['fab', '-i', '~/.ssh/zbx_id_rsa.pem',
                       '-u', 'ubuntu', '-H', ','.join(l),
                       work], stdout=PIPE, stderr=STDOUT, close_fds=True, bufsize=1, universal_newlines=True)
    else:
        popen = Popen('fab -i ~/.ssh/%s.pem -u ubuntu -P -H %s %s' % (your_key_name, ','.join(l), work),
                      shell=True, stdout=PIPE, stderr=STDOUT, close_fds=True, bufsize=1, universal_newlines=True)
    thread = Thread(target=monitor, args=[popen.stdout, N, t])
    thread.daemon = True
    thread.start()

    popen.wait()
    thread.join(timeout=1)

    return  # to comment the following lines
    counter = 0
    while True:
        line = popen.stdout.readline()
        if not line: break
        if 'synced transactions set' in line:
            counter += 1
        if counter >= N - t:
            break
        print
        line  # yield line
        sys.stdout.flush()
    ending_time = time.time()
    print
    'Latency from client scope:', ending_time - starting_time


def callFab(s, work):  # Deprecated
    print
    Popen(['fab', '-i', '~/.ssh/zbx_id_rsa.pem',
           '-u', 'ubuntu', '-H', ','.join(getAddrFromEC2Summary(s)),
           work])


# short-cuts

c = callFabFromIPList


def sk():
    c(getIP(), 'syncKeys')


def sConf():
    c(getIP(), 'syncConf')


def onlySConf():
    c(getIP(), 'syncJsonOnly')


def sKeys():
    c(getIP(), 'syncKeys')

def copyKeys(startID_=0,endID_=0):
    c(getIP(), 'copyKey:%d,%d' % (startID_,endID_))


def sExcute():
    c(getIP(), 'syncExcute')


def sClientShell():
    c(getClientIP(), 'syncClientShell')


def sLimit():
    c(getIP(), 'syncLimit')


def setUlimit():
    c(getIP(), 'set_ulimit')


def runClient(id, type, number, msg, freq):
    print('runClient:%d,%d,%d,%s,%d' % (id, type, number, msg, freq))
    c(getIP(), 'runClient:%d,%d,%d,%s,%d' % (id, type, number, msg, freq))


def computeThroughput(n=4, date=''):
    throughputArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in throughputArr:
                    throughputArr[int(item[3])] = []
                historyThroughput = throughputArr[int(item[3])]
                if int(item[4]) > 0:
                    historyThroughput.append(int(item[4]))
                    throughputArr[int(item[3])] = historyThroughput
            file.close()

    throughputDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(throughputArr.keys()):
            Sum = 0
            numLen = len(throughputArr[key])
            for i in throughputArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "(" + str(key) + "," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            throughputDict[key] = avg
    dfile.close()
    return throughputDict


def computeQuorumThroughput(n=4, date='', filter=False):
    throughputArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in throughputArr:
                    throughputArr[int(item[3])] = list()
                historyThroughput = throughputArr[int(item[3])]
                if int(item[5]) > 0:
                    historyThroughput.append(int(item[5]))
                    throughputArr[int(item[3])] = historyThroughput
            file.close()

    throughputDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(throughputArr.keys()):
            Sum = 0
            numLen = len(throughputArr[key])

            if filter:
                if numLen > n:
                    throughputArr[key].sort()
                    throughputArr[key] = throughputArr[key][0: 2 * int(numLen / 3)]
                    numLen = len(throughputArr[key])
            for i in throughputArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "(" + str(key) + "," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            throughputDict[key] = avg
    dfile.close()
    return throughputDict


def computeMidThroughput(n=4, date=''):
    throughputArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in throughputArr:
                    throughputArr[int(item[3])] = list()
                historyThroughput = throughputArr[int(item[3])]
                if int(item[5]) > 0:
                    historyThroughput.append(int(item[5]))
                    throughputArr[int(item[3])] = historyThroughput
            file.close()

    throughputDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(throughputArr.keys()):
            Sum = 0
            numLen = len(throughputArr[key])

            if numLen > n:
                throughputArr[key].sort()
                throughputArr[key] = throughputArr[key][int(numLen / 4): 3 * int(numLen / 4)]
                numLen = len(throughputArr[key])

            for i in throughputArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "(" + str(key) + "," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            throughputDict[key] = avg
    dfile.close()
    return throughputDict

def latency_throughput(n=4, date='', filter=False):
    throughputArr = dict()
    latencyArr = dict()
    cresult = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                bkey = int(item[3])
                if bkey not in throughputArr:
                    throughputArr[bkey] = []
                    latencyArr[bkey] = []
                historyThroughput = throughputArr[bkey]
                historyLatency = latencyArr[bkey]
                if int(item[5]) > 0:
                    historyThroughput.append(int(item[5]))
                    throughputArr[bkey] = historyThroughput
                    historyLatency.append(int(item[6]))
                    latencyArr[bkey] = historyLatency
            file.close()
    throughputDict = dict()
    latencyDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(throughputArr.keys()):
            Sum = 0
            numLen = len(throughputArr[key])

            if filter:
                if numLen > 1:
                    throughputArr[key].sort()
                    throughputArr[key] = throughputArr[key][0: int(numLen / 2)]
                    latencyArr[key].sort(reverse=True)
                    latencyArr[key] = latencyArr[key][0: int(numLen / 2)]
                    numLen = len(throughputArr[key])

            for i in throughputArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)

            throughputDict[key] = avg

            Sum = 0
            numLen = len(latencyArr[key])
            for i in latencyArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            latencyDict[key] = avg

            p = "(" + str(throughputDict[key]) + "," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            cresult[throughputDict[key]] = float(avg) / 1000

    dfile.close()
    return cresult


def computeLatency(n=4, date='', filter=False):
    latencyArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in latencyArr:
                    latencyArr[int(item[3])] = []
                historyThroughput = latencyArr[int(item[3])]
                if int(item[6]) > 0:
                    historyThroughput.append(int(item[6]))
                    latencyArr[int(item[3])] = historyThroughput
            file.close()

    latencyDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(latencyArr.keys()):
            Sum = 0
            numLen = len(latencyArr[key])

            if filter:
                if numLen > n:
                    latencyArr[key].sort(reverse=True)
                    latencyArr[key] = latencyArr[key][0: 2 * int(numLen / 3)]
                    numLen = len(latencyArr[key])

            for i in latencyArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "($b=" + str(key) + "$," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            latencyDict[key] = avg
    dfile.close()
    return latencyDict


def computeHacssLatency(n=4, date='', filter=False):
    latencyArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in latencyArr:
                    latencyArr[int(item[3])] = []
                historyThroughput = latencyArr[int(item[3])]
                if len(item) > 9:
                    if int(item[9]) > 0:   #HacssLatency
                        historyThroughput.append(int(item[9])) #HacssLatency
                        latencyArr[int(item[3])] = historyThroughput
            file.close()

    latencyDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(latencyArr.keys()):
            Sum = 0
            numLen = len(latencyArr[key])

            if filter:
                if numLen > n:
                    latencyArr[key].sort(reverse=True)
                    latencyArr[key] = latencyArr[key][0: 2 * int(numLen / 3)]
                    numLen = len(latencyArr[key])

            for i in latencyArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "($b=" + str(key) + "$," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            latencyDict[key] = avg
    dfile.close()
    return latencyDict

def computeRBCLatency(n=4, date='', filter=False):
    latencyArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in latencyArr:
                    latencyArr[int(item[3])] = []
                historyThroughput = latencyArr[int(item[3])]
                if len(item) > 9:
                    if int(item[9]) > 0:   #HacssLatency
                        historyThroughput.append(int(item[9])) #HacssLatency
                        latencyArr[int(item[3])] = historyThroughput
            file.close()

    latencyDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(latencyArr.keys()):
            Sum = 0
            numLen = len(latencyArr[key])

            if filter:
                if numLen > n:
                    latencyArr[key].sort(reverse=True)
                    latencyArr[key] = latencyArr[key][0: 2 * int(numLen / 3)]
                    numLen = len(latencyArr[key])

            for i in latencyArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "($b=" + str(key) + "$," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            latencyDict[key] = avg
    dfile.close()
    return latencyDict


def computeABALatency(n=4, date='', filter=False):
    latencyArr = dict()
    for ID in range(0, n):
        if date == '':
            filename = "var/log/" + str(ID) + "/" + time.strftime("%Y%m%d", time.localtime()) + "_Eva.log"
        else:
            filename = "var/log/" + str(ID) + "/" + date + "_Eva.log"
        with open(filename, 'r') as file:
            while True:
                line = file.readline()
                if not line:
                    break
                item = line.strip().split(' ')
                if int(item[3]) not in latencyArr:
                    latencyArr[int(item[3])] = []
                historyThroughput = latencyArr[int(item[3])]
                if len(item) > 10:
                    if int(item[10]) > 0:   #abaLatency
                        historyThroughput.append(int(item[10])) #abaLatency
                        latencyArr[int(item[3])] = historyThroughput
            file.close()

    latencyDict = dict()
    with open("data.txt", 'w') as dfile:
        for key in sorted(latencyArr.keys()):
            Sum = 0
            numLen = len(latencyArr[key])

            if filter:
                if numLen > n:
                    latencyArr[key].sort(reverse=True)
                    latencyArr[key] = latencyArr[key][0: 2 * int(numLen / 3)]
                    numLen = len(latencyArr[key])

            for i in latencyArr[key]:
                Sum = Sum + i
            avg = int(Sum / numLen)
            p = "($b=" + str(key) + "$," + str(float(avg) / 1000) + ") "
            dfile.write(p)
            latencyDict[key] = avg
    dfile.close()
    return latencyDict


def install():
    c(getIP(), 'installPBCandCharm')


def id():
    c(getIP(), 'install_dependencies')


def gp():
    c(getIP(), 'git_pull')


def rp(srp):
    c(getIP(), 'runProtocol:%s' % srp)


def rplocal(srp):
    c(getIP(), 'runProtocol_local:%s' % srp)


def killServer():
    c(getIP(), 'kill_Server')

def killClient():
    c(getIP(), 'kill_Client')


def killAll():
    c(getIP(), 'kill_All')


if __name__ == '__main__':
    try:
        __IPYTHON__
    except NameError:
        parser = argparse.ArgumentParser()
        parser.add_argument('access_key', help='Access Key')
        parser.add_argument('secret_key', help='Secret Key')
        args = parser.parse_args()
        access_key = args.access_key
        secret_key = args.secret_key

        print(access_key)
        print(secret_key)

        import IPython

        IPython.embed()
