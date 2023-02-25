package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var numTotalPlayer = 28
var numVipTable = 3
var numVip = 2
var numVipLimit = 4

var numTotalRound = 5
var numTotalTable = numTotalPlayer / 4

var tables = [10][10][4]int{}
var foundTables = [10][10][4]int{}
var balanceTables = [10][10][4]int{}
var bolFoundTablePlayer = false
var bolFoundBalance = false

const MAX_PLAYER = 50
var flagLiveTable = [MAX_PLAYER]int{}
var collisions = [MAX_PLAYER][MAX_PLAYER]int{}
var countVipTable = [MAX_PLAYER]int{}
var foundCountVipTable = [MAX_PLAYER]int{}
var positions = [MAX_PLAYER][MAX_PLAYER]int{}


func genSearchQueue(beginStep int, endStep int) [50]int{
	var searchQueue = [50]int{}

	for i:=0;i<numTotalPlayer;i++ {
		searchQueue[i] = 0
	}

	var idxQueue = beginStep

	for i:=0;i<=numTotalRound;i++ {
		for j:=beginStep;j<=endStep;j++ {
			if countVipTable[j] == i {
				searchQueue[idxQueue] = j;
				idxQueue ++
			}
		}
	}

	return searchQueue
	//fmt.Println(searchQueue)
}

func dfsTablePlayer(round int, step int, beginStep int, endStep int, endTable int, limitVipTable int, searchQueue [50]int) {
	// exit
	if bolFoundTablePlayer {
		return
	}

	// found
	if round + 1 > numTotalRound {
		bolFoundTablePlayer = true
		foundTables = tables
		foundCountVipTable = countVipTable
		return
	}

	var curPlayer = searchQueue[step]
	// search
	for idxTable:=0;idxTable < endTable;idxTable++ {

		// live table
		if idxTable == 0 && flagLiveTable[curPlayer] == 1 {
			continue
		}

		// vip table limit
		if countVipTable[curPlayer] + 1 > limitVipTable {
			continue
		}

		var idxPosition = -1

		for i:=0;i < 4;i++{
			if tables[round][idxTable][i] == 0 {
				var bolPosition = true

				// collision
				for j:=0;j < i;j++ {
					if collisions[curPlayer][tables[round][idxTable][j]] != 0 {
						bolPosition = false
						break
					}
				}

				if bolPosition {
					idxPosition = i
					break
				}
			}
		}

		// position invalid
		if idxPosition < 0 {
			continue
		}

		// set
		tables[round][idxTable][idxPosition] = curPlayer
		if idxTable == 0 {
			flagLiveTable[curPlayer] = 1
		}
		for j:=0;j < 4;j++ {
			if tables[round][idxTable][j] != 0 {
				collisions[tables[round][idxTable][j]][curPlayer] = 1
				collisions[curPlayer][tables[round][idxTable][j]] = 1
			}
		}
		if idxTable < numVipTable {
			countVipTable[curPlayer] += 1
		}

		// next
		if step + 1 <= endStep {
			dfsTablePlayer(round, step + 1, beginStep, endStep, endTable, limitVipTable, searchQueue)
		} else {
			var curQueue = genSearchQueue(beginStep, endStep)
			dfsTablePlayer(round + 1, beginStep, beginStep, endStep, endTable, limitVipTable, curQueue)
		}

		// unset
		tables[round][idxTable][idxPosition] = 0
		if idxTable == 0 {
			flagLiveTable[curPlayer] = 0
		}
		for j:=0;j < 4;j++ {
			if tables[round][idxTable][j] != 0 {
				collisions[tables[round][idxTable][j]][curPlayer] = 0
				collisions[curPlayer][tables[round][idxTable][j]] = 0
			}
		}
		if idxTable < numVipTable {
			countVipTable[curPlayer] -= 1
		}
	}
}

func rebalanceTable(round int, step int, flagPlayer [50]int) {
	if bolFoundBalance {
		return
	}

	if round+1 > numTotalRound {
		foundTables = balanceTables
		bolFoundBalance = true
		return
	}

	var idxPosition = step % 4
	var idxTable = step / 4

	for i := 0; i < 4; i++ {
		var positionPlayer = foundTables[round][idxTable][i]
		if flagPlayer[positionPlayer] != 0 {
			continue
		}

		// max position count
		/*
		var maxPositionCount = -1
		var idxMaxPosition = -1
		for j:=0;j<4;j++{
			if positions[positionPlayer][j] > maxPositionCount {
				maxPositionCount = positions[positionPlayer][j]
				idxMaxPosition = j
			}
		}
		if maxPositionCount >= 2 {
			if idxPosition == idxMaxPosition {
				continue
			} else if positions[positionPlayer][idxPosition] + 1 >= 2 {
				continue
			}
		}

		 */

		// min position

		var minPositionCount = 50
		for j:=0;j<4;j++ {
			if positions[positionPlayer][j] < minPositionCount {
				minPositionCount = positions[positionPlayer][j]
			}
		}
		if positions[positionPlayer][idxPosition] > minPositionCount {
			continue
		}

		// set
		balanceTables[round][idxTable][idxPosition] = positionPlayer
		positions[positionPlayer][idxPosition] += 1
		flagPlayer[positionPlayer] = 1

		// next
		if step + 1 >= numTotalPlayer {
			rebalanceTable(round+1, 0, [50]int{})
		} else {
			rebalanceTable(round, step+1, flagPlayer)
		}

		// unset
		balanceTables[round][idxTable][idxPosition] = 0
		positions[positionPlayer][idxPosition] -= 1
		flagPlayer[positionPlayer] = 0
	}
}

func printCountTable() {
	var countPlayerTable = [50]int{}

	for i:=0;i < numTotalRound;i++ {
		for j:=0;j < numTotalTable;j++ {
			for k:=0;k < 4;k++ {

				if j < numVipTable {
					countPlayerTable[foundTables[i][j][k]] += 1
				}
			}
		}
	}

	for i:=1;i <= numTotalPlayer;i++ {
		fmt.Println(i, countPlayerTable[i])
	}
}

func printTable()  {
	for i:=0;i < numTotalRound;i++ {
		for j:=0;j < numTotalTable;j++ {
			for k:=0;k < 4;k++ {
				fmt.Printf("%d ", foundTables[i][j][k])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}


func saveToFile()  {
	var fileName = "table1.txt"
	var byteContent []byte
	byteContent , _ = json.Marshal(123)
	os.WriteFile(fileName, byteContent, 0644)
}

func loadFromFile() {
	var fileName = "table1.txt"
	var byteContent []byte
	byteContent,_ = os.ReadFile(fileName)
	fmt.Println(byteContent)
}

func validateFoundTable() bool {
	var bolValid = true

	for i:=0;i<numTotalRound;i++ {
		var flagPlayer = [50]int{}

		for j:=0;j<numTotalTable;j++ {
			for k:=0;k<4;k++{
				flagPlayer[ foundTables[i][j][k] ] += 1
			}
		}

		for l:=1;l<=numTotalPlayer;l++ {
			if flagPlayer[l] != 1 {
				bolValid = false
				return bolValid
			}
		}
	}

	return bolValid
}

func main()  {
	fmt.Println( "search player =", numTotalPlayer, "vip =", numVip, "vip table =", numVipTable, "vip table limit=", numVipLimit)

	bolFoundTablePlayer = false
	if numVip > 0 {
		var curQueue = genSearchQueue(1, numVip)
		dfsTablePlayer(0 ,1, 1, numVip, numVipTable, numTotalRound, curQueue)
	} else {
		bolFoundTablePlayer = true
	}

	if bolFoundTablePlayer {
		fmt.Println("found vip table")
	} else {
		fmt.Println("not found vip table")
		return
	}

	tables = foundTables
	bolFoundTablePlayer = false
	var curQueue = genSearchQueue(numVip + 1, numTotalPlayer)
	dfsTablePlayer(0, numVip + 1, numVip + 1, numTotalPlayer, numTotalTable, numVipLimit, curQueue)
	if bolFoundTablePlayer && validateFoundTable() {
		fmt.Println("found total table")
	} else {
		fmt.Println("not found total table")
		return
	}

	tables = foundTables
	bolFoundBalance = false
	rebalanceTable(0, 0, [50]int{})
	if bolFoundBalance && validateFoundTable() {
		fmt.Println("found balance table")
	} else {
		fmt.Println("not found balance table")
		return
	}

	printTable()
	printCountTable()

	//saveToFile()
	//loadFromFile()
}
