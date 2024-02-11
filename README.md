For Python:

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




Example on how to install for go application:



install golang on platform
edit conf_location to point to /etc/openvpn/server/user.db

```
make build
```

copy open_vpn_auth to the openvpn server /usr/sbin

On OpenVPN Server:

make sure openvpn has access to open_vpn_auth

```
chgrp openvpn /usr/sbin/open_vpn_auth
chmod 750 /usr/sbin/open_vpn_auth
```

Edit the /etc/openvpn/server/server.conf:

```
script-security 3
auth-user-pass-verify /usr/sbin/open_vpn_auth via-env
```

Create the database:
```
/usr/sbin/open_vpn_auth -c
```
To Create A User:
```
/usr/sbin/open_vpn_auth  -a -u david -p pass

chgrp openvpn /etc/openvpn/server/user.db
chmod 740 /etc/openvpn/server/user.db
```
It should be good to go.
As long as the user.db has 740 permissions and is grouped by openvpn







