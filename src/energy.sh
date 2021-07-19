!#/bin/bash

curl https://monitoringapi.solaredge.com/site/1368783/energyDetails?api_key=$SOLAR_EDGE_API_KEY&meters=PRODUCTION,CONSUMPTION&timeUnit=DAY&startTime=2021-04-19%2000:00:00&endTime=2021-07-18%2000:00:00 > ../data/energy.json
js
