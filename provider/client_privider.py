
import sys
sys.path.append('..')

import grpc
import json, os
import threading
import time
import sys

from api import synerex_pb2, synerex_pb2_grpc
from nodeapi import nodeid_pb2, nodeid_pb2_grpc

NODE_URL = os.getenv('NODE_URL', 'onemile.synergic.mobi')
NODE_PORT = os.getenv('NODE_PORT', 9990)

SYNEREX_URL = os.getenv('SYNEREX_URL', 'onemile.synergic.mobi')
SYNEREX_PORT = os.getenv('SYNEREX_PORT', 10000)


class NodeClient():
    node_id = None
    secret = None
    keepalive_duration = None
    update_count = 0
    node_name = None

    def __init__(self, url=NODE_URL, port=NODE_PORT):
        self.channel = grpc.insecure_channel('%s:%d' % (url, port))
        self.stub = nodeid_pb2_grpc.NodeStub(self.channel)
        self.register()
        print("Connected to Node Server, with nodeid:", self.node_id, "nodename:", self.node_name)
        threading.Thread(target=self.__keep_alive).start()

    def register(self):
        nodeinfo = nodeid_pb2.NodeInfo()

        if 'runserver' in sys.argv:
            self.node_name = "kotacivic-provider-test"
        else:
            self.node_name = "kotacivic-provider-test"

        nodeinfo.node_name = self.node_name

        nodeinfo.is_server = False
        nodeid = self.stub.RegisterNode(nodeinfo)
        self.node_id = nodeid.node_id
        self.secret = nodeid.secret
        self.keepalive_duration = nodeid.keepalive_duration

    def __keep_alive(self):
        while True:
            try:
                self.update_count += 1
                nodeupdate = nodeid_pb2.NodeUpdate()
                nodeupdate.node_id = self.node_id
                nodeupdate.secret = self.secret
                nodeupdate.update_count = self.update_count
                nodeupdate.node_status = 0
                nodeupdate.node_arg = ""
                self.stub.KeepAlive(nodeupdate)
            except:
                print("Node Keep Alive Error")
            finally:
                time.sleep(5)

    def close(self):
        try:
            nodeid = nodeid_pb2.NodeID()
            nodeid.node_id = self.node_id
            nodeid.secret = self.secret
            nodeid.keepalive_duration = self.keepalive_duration
            self.stub.UnRegisterNode(nodeid)

            print("NODE server closed;")
        except Exception as err:
            pass

    def __del__(self):
        self.close()


class SynerexClient():

    nodeclient = None
    client_id = None

    def __connect__(self, url=SYNEREX_URL, port=SYNEREX_PORT):

        if self.nodeclient is not None:
            self.nodeclient.close()

        self.nodeclient = NodeClient()

        self.channel = grpc.insecure_channel('%s:%d' % (url, port))
        self.stub = synerex_pb2_grpc.SynerexStub(self.channel)
        threading.Thread(target=self.__subscribe_supply).start()
        print("Connected to Synerex Server.")
        self._machine_id = self.nodeclient.node_id
        self._epoch = 0
        self._serial_no = 0

    def __init__(self):
        self.__connect__()
        self.client_id = self.gen_snowflack()

    def gen_snowflack(self):

        nodeid = self._machine_id
        i = self._serial_no

        ts = bin(int(time.time()))[2:]
        zero = "0" * (41 - len(ts))
        ts = ts + zero

        nodeid = bin(nodeid)[2:]
        zero = "0" * (10 - len(nodeid))
        nodeid = zero + nodeid

        i = bin(i)[2:]
        zero = "0" * (12 - len(i))
        i = zero + i

        self._serial_no += 1

        return int(ts + nodeid + i, 2)


    def register_demand(self, user, *args):
        import time

        synerex_demand = synerex_pb2.Demand()
        #synerex_demand.id = int(demand.snowflake)
        synerex_demand.id = 0
        synerex_demand.sender_id = self.client_id
        synerex_demand.target_id = int(0)
        synerex_demand.type = synerex_pb2.ChannelType.Value("RIDE_SHARE")
        synerex_demand.arg_json = json.dumps("")

        response = self.stub.RegisterDemand(synerex_demand)

        result = {
            #"demand_id": int(demand.snowflake),
            "demand_id": int(0),
            "err": response.err,
            "ok": response.ok
        }

        print("register_demand to sx:", smarket_demand)

        return result

    def __del__(self):
        self.is_subscribe_supply = False


if __name__ == "__main__":
	sclient = SynerexClient()
	# Regist To Node

	# Regist To Synerex

	# SendDemand
	sclient.register_demand()

	time.Sleep(5)