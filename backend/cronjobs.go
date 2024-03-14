// cronjobs.go

package main

import (
    "time"
    "github.com/mriusd/game-contracts/drop"
)


func SecondlyCronJob() {
    // Execute tasks immediately upon startup
    executeSecondlyTasks()

    // Set up a ticker that fires every second
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Execute your tasks here at the start of each second
            executeSecondlyTasks()
        }
    }
}

func executeSecondlyTasks() {
    // Your tasks to be executed every second
    // Add your specific task code here
    // Example: fmt.Println("Task executed at", time.Now())
    drop.DroppedItems.CleanDroppedItems()
    broadcastDropMessage()
}



func TenSecondsCronJob() {
    // Calculate the time to wait until the start of the next 10-second interval
    now := time.Now()
    secondsUntilNextTenSeconds := (10 - now.Second()%10)
    nextTenSeconds := now.Add(time.Duration(secondsUntilNextTenSeconds) * time.Second)
    time.Sleep(nextTenSeconds.Sub(now))

    // Execute tasks at the start of the next 10-second interval
    executeTenSecondsTasks()

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Execute your tasks here
            executeTenSecondsTasks()
        }
    }
}

func executeTenSecondsTasks() {
    // Your tasks to be executed every 10 seconds
    // Add your code here
    // connections.RemoveQuietConnections()

}




func HourlyCronJob() {
    // Calculate the time to wait until the start of the next hour
    now := time.Now()
    nextHour := now.Truncate(time.Hour).Add(time.Hour)
    time.Sleep(nextHour.Sub(now))

    executeHourlyTasks()

    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Execute your task here
            executeHourlyTasks()
        }
    }
}



func executeHourlyTasks() {
	// setDailyAssetStats()
    // loadNewsToMap()
	// deductFundingRates()	
}

func MinutelyCronJob() {
    // Calculate the time to wait until the start of the next minute
    now := time.Now()
    nextMinute := now.Truncate(time.Minute).Add(time.Minute)
    time.Sleep(nextMinute.Sub(now))

    // Execute tasks at the start of the next minute
    go  executeMinutelyTasks()

    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Execute your tasks here
            go executeMinutelyTasks()
        }
    }
}


func executeMinutelyTasks() {
    // Your task implementation
    // RecordLastCandlesToDB()
    // deductCrossFundingRates()    
    // loadGlobalParamsFromDB()
    // updateAssetsFundingRates()
    // GetEthEstimatedGas()
    // GetSuggestedBtcFee()
    // GetSuggestedTronFee()
    // RemoveOrdersWithPriceOutOfRange()
    //loadNewsToMap()
    //deductFundingRates()
}


