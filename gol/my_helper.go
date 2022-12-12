package gol

import (
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/util"
	"uk.ac.bris.cs/gameoflife/stubs"
)

const ALIVE = 255
const DEAD = 0

func rpcStartGameOfLife(client *rpc.Client, world [][]uint8, p Params) *stubs.Response {
	request := stubs.Request{
		World: world, 
		Turns: p.Turns, 
		ImageWidth: p.ImageWidth, 
		ImageHeight: p.ImageHeight,
		WorkerNum: p.Threads,
	}
	response := new(stubs.Response)
	client.Call(stubs.StartGameOfLife, request, response)
	return response
}

func rpcGetAliveNumber(client *rpc.Client, p Params) *stubs.Response{
	request := stubs.Request{
		World: nil, 
		Turns: p.Turns, 
		ImageWidth: p.ImageWidth, 
		ImageHeight: p.ImageHeight,
		WorkerNum: p.Threads,
	}
	response := new(stubs.Response)
	client.Call(stubs.GetAliveNumber, request, response)
	return response
}

func rpcGetCurrentWorld(client *rpc.Client, p Params) *stubs.Response{
	request := stubs.Request{
		World: nil, 
		Turns: p.Turns, 
		ImageWidth: p.ImageWidth, 
		ImageHeight: p.ImageHeight,
		WorkerNum: p.Threads,
	}
	response := new(stubs.Response)
	// rpc : remote procedure call
	client.Call(stubs.GetCurrentWorld, request, response)
	return response
}

func rpcEndGameOfLife(client *rpc.Client, p Params) *stubs.Response{
	request := stubs.Request{
		World: nil, 
		Turns: p.Turns, 
		ImageWidth: p.ImageWidth, 
		ImageHeight: p.ImageHeight,
		WorkerNum: p.Threads,
	}
	response := new(stubs.Response)
	client.Call(stubs.EndGameOfLife, request, response)
	return response
}

func rpcPauseOrContinue(client *rpc.Client, p Params) *stubs.Response{
	request := stubs.Request{
		World: nil, 
		Turns: p.Turns, 
		ImageWidth: p.ImageWidth, 
		ImageHeight: p.ImageHeight,
		WorkerNum: p.Threads,
	}
	response := new(stubs.Response)
	client.Call(stubs.PauseOrContinue, request, response)
	return response
}

func rpcQuitServer(client *rpc.Client, p Params) *stubs.Response{
	request := stubs.Request{
		World: nil, 
		Turns: p.Turns, 
		ImageWidth: p.ImageWidth, 
		ImageHeight: p.ImageHeight,
		WorkerNum: p.Threads,
	}
	response := new(stubs.Response)
	client.Call(stubs.QuitServer, request, response)
	return response
}

func gatherAliveCells(p Params, world [][]uint8) []util.Cell {
	var aliveCells []util.Cell
	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {
			if world[i][j] == ALIVE {
				newCell := util.Cell{
					X: j,
					Y: i,
				}
				aliveCells = append(aliveCells, newCell)
			}
		}
	}
	return aliveCells
}