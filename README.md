# cloud-ddns
A simple DDNS to AWS API Client

this will listen for /aws/ZONEIDESTRING/?ip=ip.ad.dr.es&domain=domain.to.update
the username will be your aws_access_key_id 
password will be your aws_secret_access_key

so no configuration is required appart from the cli args of what ip / port to bind to (defaults to 127.0.0.1 port 8080)

