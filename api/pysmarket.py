import grpc
import os
from abc import abstractmethod, ABCMeta
import time
import smarket_pb2
import smarket_pb2_grpc


SMARKET_URL = os.getenv('SMARKET_URL', '127.0.0.1')
SMARKET_PORT = os.getenv('SMARKET_PORT', 10000)

class SmarketClient(metaclass=ABCMeta):
    def __init__(self, url=SMARKET_URL, port=SMARKET_PORT):
        self.channel = grpc.insecure_channel('%s:%d' % (url, port))
        self.stub = smarket_pb2_grpc.SMarketStub(self.channel)

        self.__subscribe_vehicle_status()

    @abstractmethod
    def on_vehicle_status_handler(self, data):
        pass

    def __subscribe_vehicle_status(self):
        client_id = int(time.time())

        stream = self.stub.SubscribeSupply(
        smarket_pb2.Channel(client_id=1, type=smarket_pb2.MarketType.Value("RIDE_SHARE"), arg_json=str({})))
        for data in stream:
            self.on_vehicle_status_handler(data=data)




if __name__ == "__main__":
    SmarketClient().run()
