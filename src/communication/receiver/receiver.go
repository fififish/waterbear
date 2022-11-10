/*
Receiver functions.
It implements all the gRPC services defined in communication.proto file.
*/

package receiver

import (
	wb "aba/waterbear"
	"communication"
	"config"
	"fmt"
	"google.golang.org/grpc"
	"log"
	logging "logging"
	"net"
	"os"
	pb "proto/proto/communication"
	"sync"
	"context"
	"consensus"
	"utils"
	"cryptolib"
	"broadcast/rbc"
	"broadcast/ecrbc"

	bit "aba/biasedit"
	cobalt "aba/cobalt"
	coin "aba/coin"
)

var id string
var wg sync.WaitGroup
var sleepTimerValue int
var con int 

type server struct {
	pb.UnimplementedSendServer
}

type reserver struct {
	pb.UnimplementedSendServer
}

/*
Handle replica messages (consensus normal operations)
*/
func (s *server) SendMsg(ctx context.Context, in *pb.RawMessage) (*pb.Empty, error) {
	//go handler.HandleByteMsg(in.GetMsg())
	return &pb.Empty{}, nil
}


func (s *server) SendRequest(ctx context.Context, in *pb.Request) (*pb.RawMessage, error) {
	return HandleRequest(in)
}

func (s *reserver) SendRequest(ctx context.Context, in *pb.Request) (*pb.RawMessage, error) {
	return HandleRequest(in)
}

func HandleRequest(in *pb.Request) (*pb.RawMessage, error) {
	/*h := cryptolib.GenHash(in.GetRequest())
	rtype := in.GetType()*/


	/*go handler.HandleRequest(in.GetRequest(), utils.BytesToString(h))


	replies := make(chan []byte)
	go handler.GetResponseViaChan(utils.BytesToString(h), replies)
	reply := <-replies*/
	wtype := in.GetType()
	switch wtype{
	case pb.MessageType_WRITE_BATCH:
		consensus.HandleBatchRequest(in.GetRequest())
		reply := []byte("batch rep")

		return &pb.RawMessage{Msg: reply}, nil
	default:
		h := cryptolib.GenHash(in.GetRequest())
		go consensus.HandleRequest(in.GetRequest(), utils.BytesToString(h))

		reply := []byte("rep")

		return &pb.RawMessage{Msg: reply}, nil
	}

}

func (s *server) RBCSendByteMsg(ctx context.Context, in *pb.RawMessage) (*pb.Empty, error) {
	go rbc.HandleRBCMsg(in.GetMsg())
	return &pb.Empty{}, nil
}

func (s *server) ECRBCSendByteMsg(ctx context.Context, in *pb.RawMessage) (*pb.Empty, error) {
	go ecrbc.HandleECRBCMsg(in.GetMsg())
	return &pb.Empty{}, nil
}


func (s *server) ABASendByteMsg(ctx context.Context, in *pb.RawMessage) (*pb.Empty, error) {
	switch consensus.ConsensusType(con){
	case consensus.ITBFT:
		go bit.HandleABAMsg(in.GetMsg())
	case consensus.BEATCobalt:
		go cobalt.HandleABAMsg(in.GetMsg())
	case consensus.WaterBearBiased:
		go wb.HandleABAMsg(in.GetMsg())
	default:
		log.Fatalf("consensus type %v not supported in ABASendByteMsg function", con)
	}
	
	return &pb.Empty{}, nil
}

func (s *server) PRFSendByteMsg(ctx context.Context, in *pb.RawMessage) (*pb.Empty, error) {
	go coin.HandleCoinMsg(in.GetMsg())
	return &pb.Empty{}, nil
}

/*
Handle join requests for both static membership (initialization) and dynamic membership.
Each replica gets a conformation for a membership request.
*/
func (s *server) Join(ctx context.Context, in *pb.RawMessage) (*pb.RawMessage, error) {
	
	reply := []byte("hi")//handler.HandleJoinRequest(in.GetMsg())
	result := true

	return &pb.RawMessage{Msg: reply, Result: result}, nil
}

/*
Register rpc socket via port number and ip address
*/
func register(port string, splitPort bool) {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		p := fmt.Sprintf("[Communication Receiver Error] failed to listen %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		os.Exit(1)
	}
	if config.FetchVerbose() {
		p := fmt.Sprintf("[Communication Receiver] listening to port %v", port)
		logging.PrintLog(config.FetchVerbose(), logging.NormalLog, p)
	}

	log.Printf("ready to listen to port %v", port)
	go serveGRPC(lis, splitPort)

}

/*
Have serve grpc as a function (could be used together with goroutine)
*/
func serveGRPC(lis net.Listener, splitPort bool) {
	defer wg.Done()


	if splitPort {

		s1 := grpc.NewServer(grpc.MaxRecvMsgSize(52428800), grpc.MaxSendMsgSize(52428800))

		pb.RegisterSendServer(s1, &reserver{})
		log.Printf("listening to split port")
		if err := s1.Serve(lis); err != nil {
			p := fmt.Sprintf("[Communication Receiver Error] failed to serve: %v", err)
			logging.PrintLog(true, logging.ErrorLog, p)
			os.Exit(1)
		}

		return
	}

	s := grpc.NewServer(grpc.MaxRecvMsgSize(52428800), grpc.MaxSendMsgSize(52428800))

	pb.RegisterSendServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		p := fmt.Sprintf("[Communication Receiver Error] failed to serve: %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		os.Exit(1)
	}


}

/*
Start receiver parameters initialization
*/
func StartReceiver(rid string, cons bool) {
	id = rid
	logging.SetID(rid)
	

	config.LoadConfig()
	logging.SetLogOpt(config.FetchLogOpt())
	con = config.Consensus()

	sleepTimerValue = config.FetchSleepTimer()
	if cons{
		consensus.StartHandler(rid)
	}
	
	
	if config.SplitPorts() {
		//wg.Add(1)
		go register(communication.GetPortNumber(config.FetchPort(rid)), true)
	}
	wg.Add(1)
	register(config.FetchPort(rid), false)
	wg.Wait()


}
