package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	flux "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	buildLogger()
	client := flux.NewClient("http://localhost:8086", os.Getenv("TOKEN"))
	defer client.Close()
	queryAPI := client.QueryAPI("codabool")
	usage := usage(queryAPI)
	solar := solar(queryAPI)
	log.Info().Msg("usage = " + strconv.Itoa(int(usage)))
	log.Info().Msg("solar = " + strconv.Itoa(int(solar)))
	if usage == 0 || solar == 0 {
		log.Info().Msg("no data found, exiting")
		return
	}
	dbInit(true) // set to true to migrate db
	now := time.Now()
	day := uint16(now.Weekday())
	write(Sense{Usage: usage, Solar: solar, Day: day})

	// clear rows that are over a week old
	// only perform this check at the beginning of the day
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	threeMinutesPastMidnight := midnight.Add(time.Minute * 3)
	if now.Before(threeMinutesPastMidnight) {
		clearPrevWeek(day)
	} else {
		log.Print("skip deleting old records")
	}
}

// curl "localhost:8086/query?db=codabool" --header "Authorization: Token token_here" --data-urlencode "q=SELECT time,leg,watts FROM codabool..sense_mains ORDER BY time DESC"
func usage(flux api.QueryAPI) uint16 {
	result, err := flux.Query(context.Background(), `from(bucket: "codabool")
  |> range(start: - 2m)
  |> filter(fn: (r) => r["_measurement"] == "sense_mains")
  |> filter(fn: (r) => r["_field"] == "watts")
  |> limit(n:1)`)
	if err == nil {
		sum := 0.0
		// Iterate over query response
		for result.Next() {
			// fmt.Print("iterate")
			// Notice when group key has changed
			if result.TableChanged() {
				// fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			// fmt.Println(result.Record().Value())
			if f, ok := result.Record().Value().(float64); ok {
				// fmt.Println("usage ", f)
				sum += f
			} else {
				fmt.Println("failed type assert to float64")
			}
			// fmt.Printf("value: %v\n", result.Record().Value())
		}

		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
		log.Print("usage before cast ", int(sum))
		return uint16(sum)
	} else {
		panic(err)
	}
}

func solar(flux api.QueryAPI) uint16 {
	result, err := flux.Query(context.Background(), `from(bucket: "codabool")
	|> range(start: -2m) 
	|> filter(fn: (r) => r["_field"] == "current_watts")
	|> filter(fn: (r) => r["id"] == "solar")
	|> limit(n:1)`)
	if err == nil {
		var sum float64
		// Iterate over query response
		for result.Next() {
			// fmt.Print("iterate")q
			// Notice when group key has changed
			// if result.TableChanged() {
			// 	fmt.Printf("table: %s\n", result.TableMetadata().String())
			// }
			// Access data
			// fmt.Printf("Solar: %v\n", result.Record().Value())
			if f, ok := result.Record().Value().(float64); ok {
				// fmt.Println("float ", f)
				sum = f
			} else {
				fmt.Println("failed type assert to float64")
			}
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
		log.Print("solar before cast ", int(sum))
		return uint16(sum)
	} else {
		panic(err)
	}
}
