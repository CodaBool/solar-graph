# Ref
- [source](https://github.com/lux4rd0/sense-collector/issues/9)
- [fork of source](https://github.com/lux4rd0/sense-collector/issues/9) 

# Usage
1. start docker compose of influxdb
1. write sensitive data into gen.sh file
2. run gen.sh to generate values in docker compose
3. combine docker compose data to form docker compose of both services
4. go get -d .
5. go run .