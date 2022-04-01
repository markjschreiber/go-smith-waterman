package main

import (
	"strings"
)

const (
	gapSymbol  byte = '-'
	gapPenalty int  = -2
)

type MatrixCell struct {
	score       int
	seq1Symbol  byte
	seq2Symbol  byte
	backPointer *MatrixCell
}

type Alignment struct {
	seq1  string
	seq2  string
	score int
}

//BasicScoringFunction A simple scoring function where matches score +3 and mismatches score -3. Symbols are assumed to
//be bytes from the ASCII set (uint8)
func BasicScoringFunction(seq1Symbol byte, seq2Symbol byte) int {
	if seq1Symbol == seq2Symbol {
		return 3
	}
	return -3
}

func Align(seq1 string, seq2 string, scoringFunction func(byte, byte) int, gapPenalty int) Alignment {
	// create matrix
	matrix := createMatrix(len(seq1), len(seq2))

	// fill matrix scores
	highestScoringCell := fillMatrix(matrix, seq1, seq2, scoringFunction, gapPenalty)

	// trace back
	return traceBack(highestScoringCell)
}

func createMatrix(n int, m int) *[][]*MatrixCell {
	matrix := make([][]*MatrixCell, n+1)
	for j := range matrix {
		cellColumn := make([]*MatrixCell, m+1)
		for i := range cellColumn {
			cell := MatrixCell{
				score:       0,
				seq1Symbol:  0,
				seq2Symbol:  0,
				backPointer: nil,
			}
			cellColumn[i] = &cell
		}
		matrix[j] = cellColumn
	}

	return &matrix
}

func fillMatrix(matrix *[][]*MatrixCell, seq1 string, seq2 string, scoringFunction func(uint8, uint8) int, gapPenalty int) *MatrixCell {

	highestScoringCell := (*matrix)[0][0]

	for n, column := range *matrix {
		for m := range column {
			if m == 0 || n == 0 {
				continue
			} else {
				symbol1 := seq1[n-1]
				symbol2 := seq2[m-1]
				diagonalCell := (*matrix)[n-1][m-1]
				upperCell := (*matrix)[n][m-1]
				leftCell := (*matrix)[n-1][m]
				diagonalPathScore := max(diagonalCell.score+scoringFunction(symbol1, symbol2), 0)
				upperPathScore := max(upperCell.score+gapPenalty, 0)
				leftPathScore := max(leftCell.score+gapPenalty, 0)

				cellToUpdate := (*matrix)[n][m]
				if diagonalPathScore > upperPathScore && diagonalPathScore > leftPathScore {
					cellToUpdate.score = diagonalPathScore
					cellToUpdate.backPointer = diagonalCell
					cellToUpdate.seq1Symbol = symbol1
					cellToUpdate.seq2Symbol = symbol2
					highestScoringCell = updateHighestScoringCell(cellToUpdate, highestScoringCell)
				} else if upperPathScore > diagonalPathScore && upperPathScore > leftPathScore {
					cellToUpdate.score = upperPathScore
					cellToUpdate.backPointer = upperCell
					cellToUpdate.seq1Symbol = gapSymbol
					cellToUpdate.seq2Symbol = symbol2
					highestScoringCell = updateHighestScoringCell(cellToUpdate, highestScoringCell)
				} else {
					cellToUpdate.score = leftPathScore
					cellToUpdate.backPointer = leftCell
					cellToUpdate.seq1Symbol = symbol1
					cellToUpdate.seq2Symbol = gapSymbol
					highestScoringCell = updateHighestScoringCell(cellToUpdate, highestScoringCell)
				}
			}
		}
	}

	return highestScoringCell
}

func updateHighestScoringCell(currentCell *MatrixCell, cellWithMaxScore *MatrixCell) *MatrixCell {
	if currentCell.score > cellWithMaxScore.score {
		cellWithMaxScore = currentCell
	}
	return cellWithMaxScore
}

func traceBack(highestScoringCell *MatrixCell) Alignment {
	var builder1, builder2 strings.Builder
	alignment := Alignment{}
	currentCell := highestScoringCell
	totalScore := 0

	//trace back until we find a cell with a score of 0
	stop := false
	for stop == false {
		if currentCell.score == 0 {
			stop = true
			continue
		}
		totalScore += currentCell.score
		builder1.WriteByte(currentCell.seq1Symbol)
		builder2.WriteByte(currentCell.seq2Symbol)
		currentCell = currentCell.backPointer
	}

	//fill in the alignment struct
	alignment.score = totalScore
	alignment.seq1 = reverse(builder1.String())
	alignment.seq2 = reverse(builder2.String())

	return alignment
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
