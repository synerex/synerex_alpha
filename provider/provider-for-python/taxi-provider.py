
import time
import sys
from sxutil import SynerexClient

# connect synerex and nodeid 
sclient = SynerexClient()

def demand_callback(demand):
    print("Get demand!\n", demand)
    
    # send confirm request
    supply_arg = {
        "depart_point": "名古屋",
        "depart_time": "2019/09/23 10:26:00",
        "arrive_point": "熱海",
        "arrive_time": "2019/09/23 12:16:00",
        "available_sheets": 4,
        "price": 28000,
        "traffic_type": "Taxi",
    }
    sclient.register_supply(supply_arg)


if __name__ == "__main__":
    try:
        # subscribe demand
        sclient.subscribe_demand(demand_callback)

    except KeyboardInterrupt:
        print("interrupted!\n")
        sys.exit(0)
