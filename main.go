package main

import "fmt"

func main() {
	println("Go Smith Waterman")

	alignment := Align("TGTTACGG", "GGTTGACTA", BasicScoringFunction, gapPenalty)
	fmt.Printf("%s\n%s\n", alignment.seq1, alignment.seq2)
	fmt.Printf("score = '%d'", alignment.score)
}
