package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type conjunct struct {
	positive []byte
	negative []byte
	
	banned []int
}

type solution struct {
	mu        *sync.Mutex
	bannedMap map[int]bool
}

func main() {
	entry, err := read("test.txt")
	if err != nil {
		log.Println(err.Error())
		return
	}
	var wg sync.WaitGroup
	m := make(map[int]bool)
	s := solution{mu: &sync.Mutex{}, bannedMap: m}
	for i, c := range entry {
		wg.Add(1)
		go func(con conjunct, j int) {
			defer wg.Done()
			arr := solver(con, j)
			for _, el := range arr {
				s.mu.Lock()
				s.bannedMap[el] = true
				s.mu.Unlock()
			}
		}(c, i)
	}
	wg.Wait()
	fmt.Println(s)
}

func read(filename string) ([]conjunct, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	
	scanner := bufio.NewScanner(file)
	var out []string
	for scanner.Scan() {
		data := scanner.Text()
		out = append(out, data)
	}
	var output []conjunct
	for _, val := range out {
		if val == "#" {
			break
		}
		var el conjunct
		separation := strings.Split(val, "_")
		var positive, negative []byte
		for _, char := range []byte(separation[0]) {
			positive = append(positive, char)
		}
		for _, char := range []byte(separation[1]) {
			negative = append(negative, char)
		}
		el.positive = positive
		el.negative = negative
		output = append(output, el)
	}
	
	return output, nil
}

func solver(input conjunct, i int) []int {
	return []int{i + 1, (i + 1) * 10}
}