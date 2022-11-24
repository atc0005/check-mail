#!/bin/bash

user="me@there.com"
token="adfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfa"
format="user=${user}\001auth=Bearer ${token}\001\001"

echo -en "user=${user}\001auth=Bearer ${token}\001\001" | base64 -w 0 > shell-output.txt

# Not sufficient to produce correct output (-e is needed):
# echo -n "user=${user}\001auth=Bearer ${token}\001\001" | base64 -w 0 > shell-output.txt


# go run . --token "$token" --username "$user" > go-output.txt
./xoauth2 --token "$token" --username "$user" > go-output.txt
md5sum *-output.txt
