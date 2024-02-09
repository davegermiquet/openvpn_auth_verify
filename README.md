This is a custom script for openvpn authentication using environment as password

Python Application:

To create the database

python /usr/sbin/open_vpn_auth.py  -c

To Create A User:
 python /usr/sbin/open_vpn_auth.py  -a -u david -p pass

To Delete A User:

 python /usr/sbin/open_vpn_auth.py  -d -u david

 requires python-bcrypt requiremnet.

Add this to your openvpn server config:
```
script-security 3
auth-user-pass-verify /usr/sbin/open_vpn_auth.py via-env
```
Make sure your openvpn server can read the script and file that you create.


GoLang Application Instructions

Change the default location of conf_location to your openvpn server config example:

/etc/openvpn/server/user.db

go mod tidy
go build 
create the database

./cmd -c 
./cmd -a -u user -p pass
./cmd -l
(list users)

copy over the file to /usr/sbin/open_vpn_auth

Edit the /etc/openvpn/server config


