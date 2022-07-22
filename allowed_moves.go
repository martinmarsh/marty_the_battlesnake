package main
// Compute which direction is allowed based on a give board

import (
	"math/rand"
  "log"
)

func Iabs(x int) int {
	if x < 0 {
		return -x
	} else if x == 0 {
    return 1
  }
	return x
}

func food_wighting(x, y int, food []Coord) int {
   weight := 0
   for _, food_pos := range food {
     weight += 100000 / (Iabs(food_pos.X - x) + Iabs(food_pos.Y - y))
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
  if  myHead.Y >= plan.Height - 1 || plan.Elements[myHead.X][myHead.Y + 1] > 20 {
    possibleMoves["up"] = false
  }

  // Can I move down?
  if myHead.Y <= 0 || plan.Elements[myHead.X][myHead.Y - 1] > 20 {
		possibleMoves["down"] = false
  }

  // Can I move left?
  if myHead.X <= 0 || plan.Elements[myHead.X - 1][myHead.Y] > 20 {
    possibleMoves["left"] = false
  } 

   // Can I move Right?
  if myHead.X >= plan.Width - 1 || plan.Elements[myHead.X + 1][myHead.Y] > 20 {
    possibleMoves["right"] = false
  }

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
    nextMove = safeMoves[0]
    for direction, weight := range weightings{
      if weight > highest {
         nextMove = direction
         highest = weight
      } else if weight == highest && rand.Intn(1) == 1{
         nextMove = direction   
      }  
    }
  }

  return nextMove
}
