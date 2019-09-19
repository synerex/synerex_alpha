import grpc
import os

NODE_PORT = os.getenv('NODE_PORT', 9990)
SYNEREX_PORT = os.getenv('NODE_PORT', 10000)

def sendDemand():

def connect():
	channel = grpc.insecure_channel()
	stub = synerex_pb2_grpc.SynerexStub(channel)


if __name__ == "__main__":

	# Regist To Node

	# Regist To Synerex

	# SendDemand
	for():
		sendDemand()
		time.Sleep(5)