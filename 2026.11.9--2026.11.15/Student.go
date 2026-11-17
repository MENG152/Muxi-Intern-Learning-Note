package main

import (
	"fmt"
)

type Student struct {
	Name   string
	Scores []int // 记录若干项成绩
}

func main() {
	s1 := Student{Name: "A", Scores: []int{80, 90, 100}}
	s2 := Student{Name: "B", Scores: []int{70}}

	fmt.Println("test Average", s1.Average())
	s1.AddScore(60)
	fmt.Println("test AddScore", s1.Scores)
	fmt.Println("test Top", s1.Top(2))
	fmt.Println("test Merge", s1.Merge(s2))

}

func (s Student) Average() float64 {
	if len(s.Scores) == 0 {
		return 0.0
	}
	sum := 0
	for _, score := range s.Scores {
		sum += score
	}

	return float64(sum) / float64(len(s.Scores))
}

func (s *Student) AddScore(score int) {
	s.Scores = append(s.Scores, score)
}

func (s Student) Top(n int) []int {
	if n <= 0 {
		return nil
	}
	length := len(s.Scores)

	temp := make([]int, length)
	copy(temp, s.Scores)

	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if temp[j] < temp[j+1] {
				temp[j], temp[j+1] = temp[j+1], temp[j]
			}
		}
	}

	if n > length {
		n = length
	}

	return temp[:n]
}

func (s Student) Merge(other Student) Student {
	NewScores := make([]int, 0, len(s.Scores)+len(other.Scores))
	NewScores = append(NewScores, s.Scores...)
	NewScores = append(NewScores, other.Scores...)

	return Student{
		Name:   s.Name,
		Scores: NewScores,
	}
}
