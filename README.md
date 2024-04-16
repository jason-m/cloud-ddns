# cloud-ddns
A simple DDNS Server to update DNS entries using various web APIs (ie aws route53) 
allows you to use a standard DDNS client to send updates 

 - zero configuration required no secrets stored or saved
 - running the cloud-ddns binary with no args binds to localhost:8080 combining this with a reverse proxy for SSL is encouraged
 - if you want to bind to an ip address other than localhost just run cloud-ddns IPADDRESS PORT 

Client configuration: 
 - for AWS Route53
   - point your ddns client to http://ipaddress:port/aws/ZONEIDSTRING/
   - use your AWS client id for the username and client secret for the password
  
 - for Cloudflare
   - point your ddns client to http://ipaddress:port/cloudflare/
   - set username to your zone name (ie domain.com) and password to your API Token
   - pairs nicely with cloudlfares tunneld to provide https termination and public accessibility 
