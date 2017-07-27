[![Build Status](https://travis-ci.org/n0r1sk/smail.svg?branch=master)](https://travis-ci.org/n0r1sk/smail)

# smail
Simple/Stupid/Short command line mailer (smail) binary for Ops usage written in golang with MX support

# Usage
./smail -s "test email" -t a@example.com -t a@example.com -f f@example.com -m mx.example.com -a /some/text/file

```
# ./smail -h
Usage of ./smail:
  -a string
    	Attachment [optional]
  -d	Debug [default=false]
  -f string
    	From address
  -m string
    	MX DNS record
  -s string
    	Subject text
  -t value
    	To address(es)
```
