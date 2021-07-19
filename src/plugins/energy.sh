#!/bin/bash

. ./load_secrets.sh

OUT=./data/energy.html
DATE_FMT="+%Y-%m-%d %H:%M:%S"
end=$(date "$DATE_FMT")
start=$(date -d "$END_DATE -10 days" "$DATE_FMT")

url="https://monitoringapi.solaredge.com/site/1368783/energyDetails"

energy=$(curl -G -s \
    --data-urlencode "api_key=$SOLAR_EDGE_API_KEY" \
    --data-urlencode "meters=PRODUCTION,CONSUMPTION" \
    --data-urlencode "timeUnit=DAY" \
    --data-urlencode "startTime=$start" \
    --data-urlencode "endTime=$end" \
    $url)

productionWh=$(echo $energy \
    | jq '.energyDetails.meters[] | [select(.type | contains("Production")).values[].value] | add | select(. != null)')
consumptionWh=$(echo $energy \
    | jq '.energyDetails.meters[] | [select(.type | contains("Consumption")).values[].value] | add | select(. != null)')

echo "<dl>" > $OUT
echo "<dt>10-day energy status</dt>" >> $OUT
echo "<dd>$(expr $productionWh / 1000) kWh produced ⚡</dd>" >> $OUT
echo "<dd>$(expr $consumptionWh / 1000) kWh consumed 🔌</dd>" >> $OUT
echo "</dl>" >> $OUT
