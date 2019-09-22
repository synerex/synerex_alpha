
import time
import sys
from sxutil import SynerexClient

# connect synerex and nodeid 
sclient = SynerexClient()
demand_arg = {
    "depart_point": "名古屋",
    "depart_time": "2019/09/23 10:36:00",
    "arrive_point": "熱海",
    "arrive_time": "2019/09/23 12:26:00",
    "num_adult": 4,
}

def supply_callback(supply):
    print("Get supply!", supply)
    print("Confirm the reservation?[Y/n]\n")
    answer = input('>> ')
    if answer == "n" or answer == "N":
        sclient.register_demand(demand_arg)
    else:
        print("Good Job! Finish your reservation. \n")


if __name__ == "__main__":
    try:
        # subscribe supply
        sclient.subscribe_supply(supply_callback)

        # send demand
        sclient.register_demand(demand_arg)

    except KeyboardInterrupt:
        print("interrupted!")
        sys.exit(0)