This is a custom script for openvpn authentication using environment as password

To create the database



python /usr/sbin/open_vpn_auth.py  -c

To Create A User:
 python /usr/sbin/open_vpn_auth.py  -a -u david -p pass

To Delete A User:


 python /usr/sbin/open_vpn_auth.py  -d -u david

 requires python-bcrypt requiremnet.
