package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const alphabetSize = 'z' - 'a' + 1
const FILENAME = "test2.txt"

type conjunct struct {
	positive []byte
	negative []byte

	banned []int
}

type solution struct {
	mu        *sync.Mutex
	bannedMap map[int]bool
}

type fileOutput struct {
	mu *sync.Mutex
	f  *os.File
}

func main() {

	t1 := time.Now().UnixNano()
	files, err := ioutil.ReadDir("tests")
	if err != nil {
		log.Fatal(err)
	}
	var wgOut sync.WaitGroup
	for _, f := range files {
		wgOut.Add(1)
		go func(fileName string) {
			tin := time.Now().UnixNano()
			defer wgOut.Done()
			entry, err := read("tests/" + fileName)
			if err != nil {
				log.Println(err.Error())
				return
			}
			var wg sync.WaitGroup
			m := make(map[int]bool)
			s := solution{mu: &sync.Mutex{}, bannedMap: m}
			for _, c := range entry {
				wg.Add(1)
				go func(con conjunct) {
					defer wg.Done()
					arr := solver(con)
					for _, el := range arr {
						s.mu.Lock()
						s.bannedMap[el] = true
						s.mu.Unlock()
					}
				}(c)
			}
			wg.Wait()

			t2 := time.Now().UnixNano()
			fmt.Println(fmt.Sprintf("Файл: %s : %f секунд", fileName, float64(t2-tin)/1000000000))
		}(f.Name())

	}

	wgOut.Wait()
	t2 := time.Now().UnixNano()
	fmt.Println(fmt.Sprintf("Прошло %f секунд", float64(t2-t1)/1000000000))
}

func check(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(file)
	var p []int = make([]int, 0, 50)
	for scanner.Scan() {
		p = append(p, len(scanner.Text())-1)
	}
	var oc int = 0
	for _, v := range p {
		oc += 1 << (alphabetSize - v)
	}
	if oc < (1 << alphabetSize) {
		return true, nil
	}
	return false, nil
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
			if char > 'z' || char < 'a' {
				continue
			}
			positive = append(positive, char)
		}

		for _, char := range []byte(separation[1]) {
			if char > 'z' || char < 'a' {
				continue
			}
			negative = append(negative, char)
		}
		el.positive = positive
		el.negative = negative
		output = append(output, el)
	}

	return output, nil
}

func solver(input conjunct) []int {

	var solComb []byte = make([]byte, alphabetSize)
	for _, v := range input.positive {
		solComb[v-'a'] = '0'
	}
	for _, v := range input.negative {
		solComb[v-'a'] = '1'
	}
	var arrSol [][]byte = make([][]byte, 1)
	arrSol[0] = solComb
	flag := true
	for flag == true {
		flag = false
		for _, v := range arrSol {
			for i, u := range v {
				if u == 0 {
					flag = true
					var nval1 []byte
					var nval2 []byte
					nval1 = append(nval1, v[0:i]...)
					nval1 = append(nval1, '0')
					nval1 = append(nval1, v[i+1:]...)

					nval2 = append(nval2, v[0:i]...)
					nval2 = append(nval2, '1')
					nval2 = append(nval2, v[i+1:]...)
					arrSol = remove(arrSol, v)
					arrSol = append(arrSol, nval1, nval2)
					break
				}
			}
		}
	}
	var res []int = make([]int, 0, len(arrSol))
	for _, v := range arrSol {
		res = append(res, systemToInt(v))
	}
	return res
}

func remove(l [][]byte, item []byte) [][]byte {
	for i, other := range l {
		if equal(other, item) {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func equal(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func systemToInt(system []byte) int {
	var res int = 0
	for i, v := range system {
		if v == '1' {
			res += 1 << (len(system) - i - 1)
		}
	}
	return res
}
