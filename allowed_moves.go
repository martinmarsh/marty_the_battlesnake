package main
// Compute which direction is allowed based on a give board

import (
	"log"
)


func GetAllowedMoves(state GameState) map[string]bool {

	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}
  // Don't let your Battlesnake move back in on it's own neck
  // Don't move out of board
  myHead := state.You.Body[0]
	myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")
	
  boardWidth := state.Board.Width
	boardHeight := state.Board.Height
  log.Printf("board %d x %d head %d, %d \n", boardWidth, boardHeight,  myHead.X, myHead.Y)
  
  // Can I move up?
  if myNeck.Y > myHead.Y || myHead.Y >= boardHeight - 1 {
    possibleMoves["up"] = false
  }

  // Can I move down?
  if myNeck.Y < myHead.Y ||  myHead.Y <= 0 {
		possibleMoves["down"] = false
  }

  // Can I move left?
  if myNeck.X < myHead.X || myHead.X <= 0 {
    possibleMoves["left"] = false
  } 

   // Can I move Right?
  if myNeck.X > myHead.X || myHead.X >= boardWidth -1 {
    possibleMoves["right"] = false
  }
  

  return possibleMoves
}
