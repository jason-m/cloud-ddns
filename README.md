# cloud-ddns
A simple DDNS Server to update DNS entries using various web APIs (ie aws route53) 
allows you to use a standard DDNS client to send updates 

zero configuration required no secrets stored or saved

running the cloud-ddns binary with no args binds to localhost:8080 combining this with a reverse proxy for SSL is encouraged

if you want to bind to an ip address other than localhost just run cloud-ddns IPADDRESS PORT 

configure your ddns client to use http://ipaddress:port/aws/ZONEIDSTRING/ 
