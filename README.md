sqslite
=======

Amazon SQS Lite/Simple CLI for getting messages in and out of SQS

= Send a message

```
echo 'message' | sqslite send
```

= Receive a message

```
export AWS_ACCESS_KEY_ID=whatever
export AWS_SECRET_ACCESS_KEY=whatever
sqslite receive -q queue-name -r us-east-1
```