import json
import os
import re
import subprocess
import sys
from enum import Enum
sys.path.append(os.path.abspath(__file__).replace(os.path.basename(__file__), '') + '..')
from Helper import Helper

class Actions(str, Enum):
    CHECK = 'check'
    CLEAR = 'clear'

class Wireguard:
    __LOG_FILE = os.path.abspath(__file__).replace(os.path.basename(__file__), '') + 'log_traffic_gb.json'
            
    @classmethod
    def get_current_used_traffic(cls) -> dict:
        output_lines = Helper.exec_cli_cmd('wg | grep -E "ips|sent"').splitlines()

        response = {}

        for line in output_lines:
            row = line.lstrip()
            if row.startswith("allowed ips"): 
                ip = row[13:]
                response[ip] = 1
            if row.startswith("transfer"): 
                try:
                    traffic_used = re.search(r', (\d+\.\d+\s[a-zA-Z]{3}) sent', row).group(1)
                    traffic_gb = float(traffic_used.split()[0])
                    if traffic_used.endswith('KiB') or traffic_used.endswith('MiB'): traffic_gb = 1 #round up to 1 gig
                except Exception as e:
                    traffic_gb = 1

                response[ip] = traffic_gb
        return response

    @classmethod
    def get_previous_used_traffic(cls) -> dict:
        if not os.path.exists(cls.__LOG_FILE): cls.write_to_log_file()
        with open(cls.__LOG_FILE, "r") as log_file: return json.load(log_file)

    @classmethod
    def write_to_log_file(cls) -> None:
        current_traffic = cls.get_current_used_traffic()
        log_file = open(cls.__LOG_FILE, "w+")
        log_file.write(json.dumps(current_traffic))
        log_file.close()

    @staticmethod
    def user_is_blocked(ip: str) -> bool:
        result = subprocess.run(f"iptables -C FORWARD -s {ip} -j DROP", shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if result.returncode == 0: return True
        return False

    @staticmethod
    def block_user(ip: str) -> None:
        #! add the rule second on the list (arg 2 in cmd), first: FORWARD 0.0.0.0/0 state RELATED,ESTABLISHED
        subprocess.run(f"iptables -I FORWARD 2 -s {ip} -j DROP", shell=True, stdout=subprocess.PIPE)

    @classmethod
    def remove_all_from_blocking(cls) -> None:
        data = cls.get_current_used_traffic()
        for ip in data:
            if cls.user_is_blocked(ip):
                subprocess.run(f"iptables -D FORWARD -s {ip} -j DROP", shell=True, stdout=subprocess.PIPE)
