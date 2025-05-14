import os
from Wireguard import Actions, Wireguard
import sys
sys.path.append(os.path.abspath(__file__).replace(os.path.basename(__file__), '') + '..')
from Telegram import Telegram
from Helper import Helper
from dotenv import find_dotenv, load_dotenv

if (__name__ != '__main__'): exit()
if (len(sys.argv) != 2): exit(1)
if sys.argv[1] not in Actions: exit(1)

load_dotenv(find_dotenv())

daily_traffic_limit_gb = os.getenv('DAILY_TRAFFIC_LIMIT_GB')
daily_traffic_limit_gb = 2 if daily_traffic_limit_gb is None or daily_traffic_limit_gb == '' else float(daily_traffic_limit_gb)


#! exclude all users from block, write current traffic to log and exit
if sys.argv[1] == Actions.CLEAR:
    Wireguard.remove_all_from_blocking()
    Wireguard.write_to_log_file()
    exit()

previous_traffic = Wireguard.get_previous_used_traffic()
current_traffic = Wireguard.get_current_used_traffic()
print('prev: ', previous_traffic)
print('current: ', current_traffic)

for ip, gb in current_traffic.items():
    if ip not in previous_traffic: continue
    diff = float(gb) - float(previous_traffic[ip])
    if (diff > daily_traffic_limit_gb and Wireguard.user_is_blocked(ip) == False):
        user_dir = Helper.exec_cli_cmd(f"grep -r {ip} /etc/wireguard/clients/")
        latest_handshake = Helper.exec_cli_cmd(f"wg | grep '{ip}' -A 1 | grep 'latest handshake'")
        Telegram.send(f"Wg block user withip {ip}. Spent today {round(diff, 2)} gb\n{user_dir}\n{latest_handshake}")
        Wireguard.block_user(ip)