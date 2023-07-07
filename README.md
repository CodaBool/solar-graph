# Ref
- [source](https://github.com/lux4rd0/sense-collector/issues/9)
- [fork of source](https://github.com/lux4rd0/sense-collector/issues/9) 

# Usage
1. start docker compose of influxdb
1. write sensitive data into gen.sh file
2. run gen.sh to generate values in docker compose
3. write all required env lines

```sh
# sense container
SENSE_COLLECTOR_HOST_HOSTNAME=from gen.sh
SENSE_COLLECTOR_INFLUXDB_PASSWORD=from influxdb
SENSE_COLLECTOR_INFLUXDB_URL=http://influxdb:8086/write?db=codabool
SENSE_COLLECTOR_INFLUXDB_USERNAME=codabool
SENSE_COLLECTOR_PERF_INTERVAL=60
SENSE_COLLECTOR_DISABLE_HEALTH_CHECK=true
SENSE_COLLECTOR_TOKEN=from gen.sh
SENSE_COLLECTOR_MONITOR_ID=from gen.sh
TZ=America/New_York

# GO
TOKEN=from influxdb
PG_URI=
```

4. combine docker compose and run both services
5. go get -d .
6. go run .
