package main

import (
	"fmt"
  "time"
  "math/rand"
  "os"
  "os/exec"
)

type coords struct{
  x int
  y int
}

const boardSize = 12 // 2 units are taken for borders

// includeHead bool determines if we treat head of a snake as it's part. It helps with checking lose conditions
func isIn(x int, y int, snake []coords, includeHead bool) bool{ 
  for i, v := range snake {
    if v.x == x && v.y == y{
      if i != 0{ 
        return true
      }else if includeHead{
        return true
      }
    }
  }
  return false
}

func printBoard(board[boardSize][boardSize] bool, snake []coords, apple coords) {
  fmt.Print("\033[H\033[2J") // clear screen
  for y := 0; y < boardSize; y++ {
    for x := 0; x < boardSize; x++{
      if y == 0 || x == 0 || y == boardSize - 1 || x == boardSize - 1 { // borders
        fmt.Print("# ")
      }else if isIn(x, y, snake, true){ // snake
        fmt.Print("* ")
      }else if x == apple.x && y == apple.y{ // apple
        fmt.Print("Ó ")
      }else{
        fmt.Print(". ") // empty field
      }
    }
    fmt.Println()
  }
} 

func posApple(apple *coords, snake []coords){
  for{
    apple.x = rand.Intn(boardSize - 2) + 1
    apple.y = rand.Intn(boardSize - 2) + 1
    if !isIn(apple.x, apple.y, snake, true){
      return
    }
  }
}

func isMoveValid(snake []coords) bool{
  return true
}

func main() {
  rand.Seed(time.Now().UnixNano())
  cont := true
  var board[boardSize][boardSize] bool
  snake := make([]coords, 3, (boardSize - 2) * (boardSize - 2))
  snake[0].x, snake[1].x, snake[2].x = 5, 5, 5
  snake[0].y, snake[1].y, snake[2].y = 8, 9, 10
  
  dir := 'N' // directions: north, west...
  apple := coords{}
  posApple(&apple, snake)

  tick := time.Tick(200 * time.Millisecond)
  ch := make(chan string)

  // input handler
  go func(ch chan string) {
    // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    var b []byte = make([]byte, 1)
    for {
      os.Stdin.Read(b)
      ch <- string(b)
    }
  }(ch)
  
  printBoard(board, snake, apple)
  for cont{ // game loop
    select{
      case <-tick:
        for i := len(snake) - 1; i > 0; i-- {
          snake[i] = snake[i - 1]
        }
        switch dir{
          case 'N':
            snake[0].y--
          case 'E':
            snake[0].x++
          case 'S':
            snake[0].y++
          case 'W':
            snake[0].x--
        }
	      printBoard(board, snake, apple)
        if isIn(snake[0].x, snake[0].y, snake, false) || snake[0].y == 0 || snake[0].x == 0 || snake[0].y == boardSize - 1 || snake[0].x == boardSize - 1 { // lose check
          cont = false
          fmt.Println("GG! Your score is:", len(snake))
        }else if snake[0].x == apple.x && snake[0].y == apple.y { // eat apple
          snake = append(snake, snake[len(snake) - 1])
          posApple(&apple, snake)
        }
      
      case stdin, _ := <-ch:
        switch stdin{
          // changing directions with keyboard input and prevent snake going back into itself
          case "w":
            if snake[0].y - 1 != snake[1].y {
              dir = 'N'
            } 
          case "a":
            if snake[0].x - 1 != snake[1].x {
              dir = 'W'
            } 
          case "s":
            if snake[0].y + 1 != snake[1].y {
              dir = 'S'
            } 
          case "d":
            if snake[0].x + 1 != snake[1].x {
              dir = 'E'
            } 
        }
    }
  }
}
