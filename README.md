# NIXDB - Unix /etc Database

Main goal of this project is to provide AD domain like functionality for unix environments  
LDAP and Kerberos are an overkill for network of < 20 hosts and come with burden of maintenance
NIS is almost impossible to use in multi-distro environment.

### Features

This is under heavy development but as of now we have following 

- `/login` Basic auth against PAM.d, if member of authorised groups receives JWT token to use other endpoints
- `/v1/api/passwd` Read only access to passwd file (JSON Format for transport) if requestor is a member of authorized group
- `/v1/api/group` Same as above for /etc/group


### TODO

- Support BSD Style passwd and groups
- CLI Client that provides synchronsiation on client hosts
