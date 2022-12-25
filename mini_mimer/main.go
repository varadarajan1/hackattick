package minimimer

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/bits"
	"math/rand"
	"net/http"
)

type Problem struct {
	Difficulty int   `json:"difficulty"`
	Block      Block `json:"block"`
}

type Block struct {
	Data  interface{} `json:"data"`
	Nonce int         `json:"nonce"`
}

/*
SHA 256 - Hashing algorithim - produces a 256 bit -> 32 byte representation of any data
1. Fetch the problem
2. Generate a random number
3. Compute hash of the Block
4. Get leading number of zeros
5. Evaluate if it satisfies difficulty,
	if yes - submit
	if no - repeat 2-5
*/

func Execute() {
	problem := fetchProblem()
	isNonceIdentified := false
	nonce := 0

	for !isNonceIdentified {
		nonce = rand.Int()
		problem.Block.Nonce = nonce
		data, _ := json.Marshal(problem.Block)
		s := sha256.Sum256(data)
		leadingZeros := 0
		for _, b := range s {
			if b == 0 {
				leadingZeros += 8
				continue
			}
			leadingZeros += bits.LeadingZeros8(uint8(b))
			break
		}
		isNonceIdentified = leadingZeros >= problem.Difficulty
	}
	postResult(nonce)
}

func postResult(nonce int) {
	body := map[string]int{
		"nonce": nonce,
	}
	json_data, _ := json.Marshal(body)
	resp, err := http.Post("https://hackattic.com/challenges/mini_miner/solve?access_token=a2d5a457aa62fdcb&playground=1", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res)
}

func fetchProblem() Problem {
	response, err := http.Get("https://hackattic.com/challenges/mini_miner/problem?access_token=a2d5a457aa62fdcb")
	if err != nil {
		fmt.Print(err.Error())
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	problem := Problem{}
	json.Unmarshal(responseData, &problem)
	return problem
}
