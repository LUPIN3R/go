package gol

import (
	"strconv"
	"flag"
	"net/rpc"
	"time"
	"fmt"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}


var server = flag.String("server", "127.0.0.1:12345", "IP:PORT of the server")
// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {
	// Get input from IO by filename
	c.ioCommand <- ioInput
	filename := strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight)
	c.ioFilename <- filename
	
	// TODO: Create a 2D slice to store the world.
	world := make([][]uint8, p.ImageHeight)
	for i := range world {
		world[i] = make([]uint8, p.ImageWidth)
	}
	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {
			world[i][j] = <-c.ioInput
		}
	}
	
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	// turn := 0
	
	// used for blocking when 'p' is pressed
	isPause := false
	keyPressesCheckDone := make(chan bool)
	// check keyPresses
	go func() {
		for true {
			select {
			case <-keyPressesCheckDone:
				return
			case curKey := <-keyPresses:
				switch curKey {
				case 's':
					curWorldResponse := rpcGetCurrentWorld(client, p)
					// Send ioOutput to IO to save current world
					c.ioCommand <- ioOutput
					newFilename := filename + "x" + strconv.Itoa(curWorldResponse.CurTurn)
					c.ioFilename <- newFilename
					for i := 0; i < p.ImageHeight; i++ {
						for j:= 0; j < p.ImageWidth; j++ {
							c.ioOutput <- curWorldResponse.World[i][j]
						}
					}
					// Send ImageOutputComplete to SDL
					c.events <- ImageOutputComplete{curWorldResponse.CurTurn, newFilename}

				case 'q':
					endGameOfLifeResponse := rpcEndGameOfLife(client, p)
					// Send ioOutput to IO to save current world
					c.ioCommand <- ioOutput
					newFilename := filename + "x" + strconv.Itoa(endGameOfLifeResponse.CurTurn)
					c.ioFilename <- newFilename
					for i := 0; i < p.ImageHeight; i++ {
						for j:= 0; j < p.ImageWidth; j++ {
							c.ioOutput <- endGameOfLifeResponse.World[i][j]
						}
					}
					// Send ImageOutputComplete and FinalTurnComplete to SDL
					c.events <- ImageOutputComplete{endGameOfLifeResponse.CurTurn, newFilename}
					c.events <- FinalTurnComplete{endGameOfLifeResponse.CurTurn, gatherAliveCells(p, endGameOfLifeResponse.World)}

				case 'p':
					pauseOrContinueResponse := rpcPauseOrContinue(client, p)
					if isPause {
						fmt.Println("Continuing")
					} else {
						fmt.Printf("Paused, turn = %d\n", pauseOrContinueResponse.CurTurn)
					}
					isPause = !isPause

				case 'k':
					rpcQuitServer(client, p)
				}
			}
		}
	}()
	
	// Use tickerDone to exit the thread, otherwise it will cause a fault
	tickerDone := make(chan bool)
	go func() {
		ticker := time.NewTicker(time.Second << 1)
		for true {
			select {
			case <-tickerDone:
				return
			case <-ticker.C:
				aliveNumResponse := rpcGetAliveNumber(client, p)
				c.events <- AliveCellsCount{aliveNumResponse.CurTurn, aliveNumResponse.AliveNumber}
			}
		}
	}()
		
	// TODO: Execute all turns of the Game of Life.
	response := rpcStartGameOfLife(client, world, p)

	// TODO: Report the final state using FinalTurnCompleteEvent.
	keyPressesCheckDone <- true
	tickerDone <- true

	c.ioCommand <- ioOutput
	newFilename := filename + "x" + strconv.Itoa(response.CurTurn)
	c.ioFilename <- newFilename
	for i := 0; i < p.ImageHeight; i++ {
		for j:= 0; j < p.ImageWidth; j++ {
			c.ioOutput <- response.World[i][j]
		}
	}
	// Send ImageOutputComplete and FinalTurnComplete to SDL
	c.events <- ImageOutputComplete{response.CurTurn, newFilename}
	c.events <- FinalTurnComplete{response.CurTurn, gatherAliveCells(p, response.World)}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{response.CurTurn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

