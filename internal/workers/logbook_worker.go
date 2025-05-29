package workers

import "fmt"

type LogbookRequest struct {
	FlightId string
}

var LogbookQueue = make(chan LogbookRequest, 100)

func LogbookWorker() {
	for req := range LogbookQueue {
		fmt.Println("Processing ", req.FlightId)
	}
}
