package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
)

type Bow map[string]int

const THRESHOLD = 100

func tokenize(message string) []string {
	tokens := strings.Fields(message)
	for i := range tokens {
		tokens[i] = strings.ToUpper(tokens[i])
	}
	return tokens
}

func addFileToBow(path string, bow map[string]int) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	for _, token := range tokenize(string(content)) {
		bow[token] += 1
	}
	return nil
}

func addDirToBow(path string, bow map[string]int) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		addFileToBow(path, bow)
		return nil
	})
}

func totalCount(bow Bow) int {
	count := 0
	for word := range bow {
		if bow[word] < THRESHOLD {
			continue
		}
		count += bow[word]
	}
	return count
}

func classifyFile(ham, spam Bow, hamTotalCount, spamTotalCount int, path string) (float64, float64, error) {
	totalCount := hamTotalCount + spamTotalCount

	emailBow := Bow{}
	err := addFileToBow(path, emailBow)
	if err != nil {
		return 0, 0, err
	}

	emailProb := 0.0
	hamEmailProb := 0.0
	spamEmailProb := 0.0
	hamProb := math.Log(float64(hamTotalCount) / float64(totalCount))
	spamProb := math.Log(float64(spamTotalCount) / float64(totalCount))

	for token := range emailBow {
		n := spam[token] + ham[token]
		if n < THRESHOLD {
			continue
		}
		if spam[token] != 0 {
			spamEmailProb += math.Log(float64(spam[token]) / float64(spamTotalCount))
		}
		if ham[token] != 0 {
			hamEmailProb += math.Log(float64(ham[token]) / float64(hamTotalCount))
		}
		if n != 0 {
			emailProb += math.Log(float64(n) / float64(totalCount))
		}
	}

	//fmt.Printf("hamProb = %v\n", hamProb)
	//fmt.Printf("spamProb = %v\n", spamProb)
	//fmt.Printf("hamEmailProb = %v\n", hamEmailProb)
	//fmt.Printf("spamEmailProb = %v\n", spamEmailProb)
	//fmt.Printf("emailProb = %v\n", emailProb)
	//fmt.Printf("probability of spam %v\n", spamEmailProb+spamProb-emailProb)
	//fmt.Printf("probability of ham %v\n", hamEmailProb+hamProb-emailProb)
	return spamEmailProb + spamProb - emailProb, hamEmailProb + hamProb - emailProb, nil
}

func classifyDir(ham, spam Bow, hamTotalCount, spamTotalCount int, dirPath string) (int, int, error) {
	hamOutcomeCount := 0
	spamOutcomeCount := 0
	err := filepath.WalkDir(dirPath, func(filePath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		spamProb, hamProb, err := classifyFile(ham, spam, hamTotalCount, spamTotalCount, filePath)
		if err != nil {
			return err
		}
		if spamProb > hamProb {
			spamOutcomeCount += 1
		} else {
			hamOutcomeCount += 1
		}
		return nil
	})
	return spamOutcomeCount, hamOutcomeCount, err
}

func main() {
	fmt.Println("Training...")
	ham := Bow{}
	spam := Bow{}
	for i := 1; i <= 5; i++ {
		err := addDirToBow(fmt.Sprintf("./enron%v/ham", i), ham)
		if err != nil {
			panic(err)
		}

		err = addDirToBow(fmt.Sprintf("./enron%v/spam", i), spam)
		if err != nil {
			panic(err)
		}
	}

	hamTotalCount := totalCount(ham)
	spamTotalCount := totalCount(spam)

	fmt.Println("Clasifying ham...")
	spamOutcomeCount, hamOutcomeCount, err := classifyDir(ham, spam, hamTotalCount, spamTotalCount, "./enron6/ham")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Spamcount = %v\n HamCount = %v\n", spamOutcomeCount, hamOutcomeCount)

	fmt.Println("Clasifying spam...")
	spamOutcomeCount, hamOutcomeCount, err = classifyDir(ham, spam, hamTotalCount, spamTotalCount, "./enron6/spam")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Spamcount = %v\n HamCount = %v\n", spamOutcomeCount, hamOutcomeCount)
}
