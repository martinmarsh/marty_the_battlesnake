package main
// Compute which direction is allowed based on a give board

import (
	
)


func GetAllowedMoves(state GameState, plan BoardPlan) map[string]bool {

	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}
  // Don't let your Battlesnake move back in on it's own neck
  // Don't move out of board
  //myBody := state.You.Body
  myHead := state.You.Body[0]
	//myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")
	

  // Can I move up?
  if  myHead.Y >= plan.Height - 1 || plan.Elements[myHead.X][myHead.Y + 1] > 20 {
    possibleMoves["up"] = false
  }

  // Can I move down?
  if myHead.Y <= 0 || plan.Elements[myHead.X][myHead.Y - 1] > 20{
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
  

  return possibleMoves
}
