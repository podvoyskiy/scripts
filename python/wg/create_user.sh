#!/bin/bash

if [ $# -eq 0 ]; then
    echo "username is required argument"
    exit 1
fi
username=$1

if [[ ! $username =~ ^[a-zA-Z]+$ ]]; then
    echo "only latin letters"
    exit 1
fi

dir="/etc/wireguard/clients"

if [ ! -d "$dir" ]; then
    echo "directory $dir doesn't exist"
    exit 1
fi

cd $dir

server_public_key=$(cat "../server_public_key")
server_ip=$(hostname -I | awk '{print $1}')

#? calc user id
folders=$(ls -l | grep '^d' | awk '{print $9}' | grep '^[0-9]')
last_number=$(echo "$folders" | sed 's/[^0-9]*//g' | sort -nr | head -n 1)
current_number=$((last_number + 1))
ip="10.7.0.$current_number/32"

#? create user dir & keys
mkdir "${current_number}_$username"
cd "./${current_number}_$username"
wg genkey | tee private_key | wg pubkey > public_key
private_key=$(<private_key)
public_key=$(<public_key)

#? create user config file
cat << EOF > wg_"$username".conf
[Interface]
PrivateKey = $private_key
Address = $ip
DNS = 8.8.8.8

[Peer]
PublicKey = $server_public_key
AllowedIPs = 0.0.0.0/0
Endpoint = $server_ip:51820
PersistentKeepalive = 20
EOF

#? append new peer to config wireguard
cat << EOF >> /etc/wireguard/wg0.conf

[Peer]
PublicKey = $public_key
AllowedIPs = $ip
EOF

echo "user $username created successfully"
echo "wg config file:"
cat "$dir/${current_number}_$username/wg_$username.conf"


