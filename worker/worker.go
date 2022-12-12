package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"fmt"
	"os"
	"uk.ac.bris.cs/gameoflife/stubs"
)

const ALIVE = 255
const DEAD  = 0
var isEndGOL bool = false
type WorkerCalls struct {}

func calculateAliveNeighbors(world [][]uint8, width int, height int, y, x int) int {
	var aliveNeighborNum = 0
	posBias := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for i := 0; i < 8; i++ {
		curY := (y + posBias[i][0] + height) % height
		curX := (x + posBias[i][1] + width) % width
		if world[curY][curX] == ALIVE {
			aliveNeighborNum++
		}
	}
	return aliveNeighborNum
}


func nextGeneration(world [][]uint8, startHeight int, endHeight int, width int, height int, myNextWorld [][]uint8) int {
	curAliveNum := 0
	for i := 0; i < endHeight - startHeight; i++ {
		for j := 0; j < width; j++ {
			actualHeight := i + startHeight
			curState := world[actualHeight][j]
			aliveNeighborNum := calculateAliveNeighbors(world, width, height, actualHeight, j)
			if curState == ALIVE && aliveNeighborNum < 2 {
				myNextWorld[i][j] = DEAD
			}
			if (curState == ALIVE && aliveNeighborNum == 2) || (curState == ALIVE && aliveNeighborNum == 3) {
				myNextWorld[i][j] = ALIVE
				curAliveNum++
			}
			if curState == ALIVE && aliveNeighborNum > 3 {
				myNextWorld[i][j] = DEAD
			}
			if curState == DEAD && aliveNeighborNum == 3 {
				myNextWorld[i][j] = ALIVE
				curAliveNum++
			}
		}
	}
	return curAliveNum
}

func (s *WorkerCalls) GameOfLife(req stubs.WorkerRequest, res *stubs.WorkerResponse) (err error) {
	isEndGOL = false
	nextWorld := make([][]uint8, req.HeightEnd - req.HeightStart)
	for i := range nextWorld {
		nextWorld[i] = make([]uint8, req.ImageWidth)
	}
	curAliveNum := nextGeneration(req.World, req.HeightStart, req.HeightEnd, req.ImageWidth, req.ImageHeight, nextWorld[:][:])
	res.NextWorld = nextWorld
	res.AliveNumber = curAliveNum
	isEndGOL = true
	return
}

func (s *WorkerCalls) QuitWorker(req stubs.WorkerRequest, res *stubs.WorkerResponse) (err error) {
	for !isEndGOL {

	}
	fmt.Println("The Worker Has Quitted")
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)
	return
}

func main(){
	// NOTE
	pAddr := flag.String("port","12346","Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&WorkerCalls{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
