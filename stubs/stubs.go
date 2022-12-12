package stubs

var StartGameOfLife = "GameOfLifeCalls.StartGameOfLife"
var GetAliveNumber = "GameOfLifeCalls.GetAliveNumber"
var GetCurrentWorld = "GameOfLifeCalls.GetCurrentWorld"
var EndGameOfLife = "GameOfLifeCalls.EndGameOfLife"
var PauseOrContinue = "GameOfLifeCalls.PauseOrContinue"
var QuitServer = "GameOfLifeCalls.QuitServer"

var WorkerGameOfLife = "WorkerCalls.GameOfLife"
var QuitWorker = "WorkerCalls.QuitWorker"

type Request struct {
	World [][]uint8
	Turns int
	ImageWidth  int
	ImageHeight int
	WorkerNum int
}

type Response struct {
	World [][]uint8
	CurTurn int
	AliveNumber int
}

type WorkerRequest struct {
	World [][]uint8
	ImageWidth  int
	ImageHeight int
	HeightStart int
	HeightEnd int
}

type WorkerResponse struct {
	NextWorld [][]uint8
	AliveNumber int
}