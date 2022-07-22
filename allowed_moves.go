package main
// Compute which direction is allowed based on a give board

import (
	"math/rand"
  "log"
)

func Iabs(x int) int {
	if x < 0 {
		return -x
	} 
	return x
}

func  not_dead_end(x, y int, plan *BoardPlan) bool{
  // see width sideways
  max_width := 1
  max_height := 1
  xp := 0
  yu := 0
  //width up and down
  for ys:= 0; ys < plan.Height - 1; ys++{
    width := 1
    up_right := true
    up_left := true
    down_right := true
    down_left := true
    for xw:= 1; ; xw++ {
      xp = x+xw 
      yu = y+ys
      if yu < plan.Height - 1 && xp < plan.Width - 1 && plan.Elements[xp][yu] <= 20 && up_right {
        width++
      } else {
        up_right = false
      }
      yu = y-ys
      //log.Printf("38: x %d, y %d", xp, yu)
      if yu >=0 && xp < plan.Width - 1 && plan.Elements[xp][yu] <= 20 && down_right {
        width++
      } else {
        down_right = false
      }
      xp = x-xw
      yu = y+ys
      if xp > 0 && yu < plan.Height - 1 && plan.Elements[xp][yu] <= 20 && up_left{
        width++
      } else {
        up_left = false
      }
      yu = y-ys
      if xp > 0 && yu >= 0 && plan.Elements[xp][yu] <= 20 && down_left{
        width++
      } else {
        down_left = false
      }
      
      if !(up_right || up_left || down_right || down_left) {
        break
      }
    }
  
    if width > max_width{
      max_width = width
    } 
  }
  //height left and right
  for xs:= 1; xs < plan.Height - 1; xs++{
    height := 1
    left_up := true
    left_down := true
    right_up := true
    right_down := true
    for yw:= 1; ; yw++ {
      xp = x+xs 
      yu = y+yw
      if yu < plan.Height - 1 && xp < plan.Width - 1 && plan.Elements[xp][yu] <= 20 && right_up {
        height++
      } else {
        right_up = false
      }
      yu = y-yw
      if yu >=0 && xp < plan.Width - 1 && plan.Elements[xp][yu] <= 20 && right_down {
        height++
      } else {
        right_down = false
      }
      xp = x-xs
      yu = y+yw
      if xp > 0 && yu < plan.Height - 1 && plan.Elements[xp][yu] <= 20 && left_up{
        height++
      } else {
        left_up = false
      }
      yu = y-yw
      if xp > 0 && yu >= 0 && plan.Elements[xp][yu] <= 20 && left_down{
        height++
      } else {
        left_down = false
      }
      
      if !(right_down || left_down || right_up || left_up) {
        break
      }
    }
    if height > max_height{
      max_height = height
    } 
  } 
    if max_height <= 1 || max_width <=1{
      return false
    }
  return true
  
}

func food_wighting(x, y int, food []Coord) int {
   weight := 0
   for _, food_pos := range food {
     distance := Iabs(food_pos.X - x) + Iabs(food_pos.Y - y)
     if distance <= 1 {
       weight = 100000
       break;
     }
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
  } else {
    possibleMoves["up"] = not_dead_end(myHead.X, myHead.Y + 1, plan)
  }

  // Can I move down?
  if myHead.Y <= 0 || plan.Elements[myHead.X][myHead.Y - 1] > 20 {
		possibleMoves["down"] = false
  } else {
    possibleMoves["down"] = not_dead_end(myHead.X, myHead.Y-1, plan)
  }

  // Can I move left?
  if myHead.X <= 0 || plan.Elements[myHead.X - 1][myHead.Y] > 20 {
    possibleMoves["left"] = false
  } else {
    possibleMoves["left"] = not_dead_end(myHead.X - 1, myHead.Y, plan)
  }

   // Can I move Right?
  if myHead.X >= plan.Width - 1 || plan.Elements[myHead.X + 1][myHead.Y] > 20 {
    possibleMoves["right"] = false
  } else {
    possibleMoves["right"] = not_dead_end(myHead.X +1 , myHead.Y, plan)
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
