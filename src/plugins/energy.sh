#!/bin/bash

. ./load_secrets.sh

DATE_BIN=date
if [[ $OSTYPE == 'darwin'* ]]; then
  DATE_BIN=gdate
fi

OUT=$1
DATE_FMT="+%Y-%m-%d %H:%M:%S"
end=$($DATE_BIN "$DATE_FMT")
start=$($DATE_BIN -d "$END_DATE -30 days" "$DATE_FMT")

url="https://monitoringapi.solaredge.com/site/1368783/energyDetails"

energy=$(curl -G -s \
  --data-urlencode "api_key=$SOLAR_EDGE_API_KEY" \
  --data-urlencode "meters=PRODUCTION,CONSUMPTION" \
  --data-urlencode "timeUnit=DAY" \
  --data-urlencode "startTime=$start" \
  --data-urlencode "endTime=$end" \
  $url)

productionWh=$(echo $energy |
  jq '.energyDetails.meters[] | [select(.type | contains("Production")).values[].value] | add | select(. != null)')

echo "<dl>" >$OUT
echo "<dt>30-day solar status</dt>" >>$OUT
echo "<dd>$(expr $productionWh / 1000) kWh produced ⚡</dd>" >>$OUT
echo "</dl>" >>$OUT
