alerts1='{
   "receiver": "default-receiver",
   "status": "firing",
   "alerts": [
     {
       "status": "firing",
       "labels": {
         "alertname": "DiskRunningFull",
         "dev": "sda1",
         "instance": "example3",
         "severity": "critical"
       },
       "annotations": {
         "info": "The disk sda1 is running full",
         "summary": "please check the instance example1"
       },
       "startsAt": "2018-08-17T10:19:09.269354561+01:00",
       "endsAt": "0001-01-01T00:00:00Z",
       "generatorURL": ""
     }
   ],
   "groupLabels": {},
   "commonLabels": {
     "alertname": "DiskRunningFull",
     "dev": "sda1",
     "instance": "example3",
     "severity": "critical"
   },
   "commonAnnotations": {
     "info": "The disk sda1 is running full",
     "summary": "please check the instance example1"
   },
   "externalURL": "http://max-Inspiron-7577:9093",
   "version": "4",
   "groupKey": "{}:{}"
}'
curl -XPOST -d "$alerts1" http://localhost:3000/alert
