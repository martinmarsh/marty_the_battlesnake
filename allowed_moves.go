package main

// Compute which direction is allowed based on a give board

import (
	"log"
	"math/rand"
)

func Iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func is_on_board(x, y int, plan *BoardPlan) bool {
	if y >= plan.Height || y < 0 || x >= plan.Width || x < 0 {
		return false
	}
	return true
}

func passage_width(x, y, direction_x, direction_y int, plan *BoardPlan) int {
	width := 1
	measure := []int{-1, 1}
	for _, direction := range measure {
		for scan := 1; scan < 10; scan++ {
			xs := x + scan*direction_x*direction
			ys := y + scan*direction_y*direction
			if is_on_board(xs, ys, plan) &&  plan.Elements[xs][ys] < 32 {
				width++
			} else {
				break
			}
		}
	}
	return width
}

func move_along_passage(x, y, direction_x, direction_y int, plan *BoardPlan) (int, int, int, int) {
	width := 0
  volume := 0
	max_width := 1
	min_width := 10000
	scan := 0
	for {
		xs := x + scan*direction_x
		ys := y + scan*direction_y
		// width direction is at right angles
		direction_width_y := 1
		direction_width_x := 0
		if direction_x == 0 {
			direction_width_y = 0
			direction_width_x = 1
		}

		if is_on_board(xs, ys, plan) && plan.Elements[xs][ys] < 32 {
			width = passage_width(xs, ys, direction_width_x, direction_width_y, plan)
			volume += width
      if max_width < width {
				max_width = width
			}
			if min_width > width {
				min_width = width
			}
			scan++
		} else {
			break
		}
	}
	return max_width, min_width, scan, volume
}

func not_dead_end(x, y, direction_x, direction_y int, plan *BoardPlan, state *GameState) bool {
  max_width, min_width, length, volume := move_along_passage(x, y, direction_x, direction_y, plan)
  log.Printf("Space for x:%d, y:%d, d: %d %d, %d %d %d %d/n", x, y, direction_x, direction_y, max_width, min_width, length, volume)
  
  if max_width < 2 && length > int(state.You.Length) {
    return false
  }

  if volume < int(state.You.Length) + 1 {
    return false
  }

	return true

}

func food_wighting(x, y int, food []Coord) int {
	weight := 0
	for _, food_pos := range food {
		distance := Iabs(food_pos.X-x) + Iabs(food_pos.Y-y)
		if distance <= 1 {
			weight = 100000
			break
		}
		weight += 100000 / (Iabs(food_pos.X-x) + Iabs(food_pos.Y-y))
	}
	return weight
}

func GetBestMove(state *GameState, plan *BoardPlan) string {

	var nextMove string

	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}
	// Don't let your Battlesnake move back in on it's own neck
	// Don't move out of board

	myHead := state.You.Body[0]
	// Can I move up?
	if myHead.Y >= plan.Height-1 || plan.Elements[myHead.X][myHead.Y+1] > 20 {
		possibleMoves["up"] = false
	} else {
		possibleMoves["up"] = not_dead_end(myHead.X, myHead.Y+1, 0, 1, plan, state)
	}

	// Can I move down?
	if myHead.Y <= 0 || plan.Elements[myHead.X][myHead.Y-1] > 20 {
		possibleMoves["down"] = false
	} else {
		possibleMoves["down"] = not_dead_end(myHead.X, myHead.Y-1, 0, -1, plan, state)
	}

	// Can I move left?
	if myHead.X <= 0 || plan.Elements[myHead.X-1][myHead.Y] > 20 {
		possibleMoves["left"] = false
	} else {
		possibleMoves["left"] = not_dead_end(myHead.X-1, myHead.Y, -1, 0, plan, state)
	}

	// Can I move Right?
	if myHead.X >= plan.Width-1 || plan.Elements[myHead.X+1][myHead.Y] > 20 {
		possibleMoves["right"] = false
	} else {
		possibleMoves["right"] = not_dead_end(myHead.X+1, myHead.Y, 1, 0, plan, state)
	}

	// Avoid dead ends

	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		log.Print(safeMoves)
		// Find move weighted towards food
		weightings := make(map[string]int)
		food := state.Board.Food

		for _, direction := range safeMoves {
			log.Print(direction)
			switch direction {
			case "up":
				weightings[direction] = food_wighting(myHead.X, myHead.Y+1, food)
			case "down":
				weightings[direction] = food_wighting(myHead.X, myHead.Y-1, food)
			case "left":
				weightings[direction] = food_wighting(myHead.X-1, myHead.Y, food)
			case "right":
				weightings[direction] = food_wighting(myHead.X+1, myHead.Y, food)
			}
		}
		log.Print(weightings)
		highest := 0
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		for direction, weight := range weightings {
			if weight > highest {
				nextMove = direction
				highest = weight
			} else if weight == highest && rand.Intn(1) == 1 {
				nextMove = direction
			}
		}
	}

	return nextMove
}
