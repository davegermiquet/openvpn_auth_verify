Example on how to install:


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







