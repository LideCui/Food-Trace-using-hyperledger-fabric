/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a fruit
type SmartContract struct {
	contractapi.Contract
}

// fruit describes basic details of what makes up a fruit
type Fruit struct {
	SerialNum   string `json:"SerialNum"`
	Name  string `json:"Name"`
	Origin string `json:"Origin"`
	Date  string `json:"Date"`
	Owner string `json:"Owner`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Fruit
}

// InitLedger adds a base set of fruits to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fruits := []Fruit{
		Fruit{SerialNum: "1", Name: "Grape", Origin: "South Africa", Date: "2020/07/01", Owner: "Walmart"},
		Fruit{SerialNum: "2", Name: "Banana", Origin: "Filipino", Date: "2020/07/01", Owner: "Provigo"},
		Fruit{SerialNum: "3", Name: "Apple", Origin: "China", Date: "2020/07/01", Owner: "Walmart"},
		Fruit{SerialNum: "4", Name: "Strawberry", Origin: "USA", Date: "2020/07/01", Owner: "Walmart"},
		Fruit{SerialNum: "5", Name: "Blueberry", Origin: "Canada", Date: "2020/07/01", Owner: "Provigo"},
	}

	for i, fruit := range fruits {
		fruitAsBytes, _ := json.Marshal(fruit)
		err := ctx.GetStub().PutState("fruit"+strconv.Itoa(i), fruitAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Createfruit adds a new fruit to the world state with given details
func (s *SmartContract) Createfruit(ctx contractapi.TransactionContextInterface, fruitNumber string, serialNum string, name string, origin string, date string, owner string) error {
	fruit := Fruit{
		SerialNum:   serialNum,
		Name:  name,
		Origin: origin,
		Date:  date,
		Owner: owner,
	}

	fruitAsBytes, _ := json.Marshal(fruit)

	return ctx.GetStub().PutState(fruitNumber, fruitAsBytes)
}

// Queryfruit returns the fruit stored in the world state with given id
func (s *SmartContract) Queryfruit(ctx contractapi.TransactionContextInterface, fruitNumber string) (*Fruit, error) {
	fruitAsBytes, err := ctx.GetStub().GetState(fruitNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if fruitAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", fruitNumber)
	}

	fruit := new(Fruit)
	_ = json.Unmarshal(fruitAsBytes, fruit)

	return fruit, nil
}

// QueryAllfruits returns all fruits found in world state
func (s *SmartContract) QueryAllfruits(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "fruit0"
	endKey := "fruit99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		fruit := new(Fruit)
		_ = json.Unmarshal(queryResponse.Value, fruit)

		queryResult := QueryResult{Key: queryResponse.Key, Record: fruit}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangefruitOwner updates the owner field of fruit with given id in world state
func (s *SmartContract) ChangefruitOwner(ctx contractapi.TransactionContextInterface, fruitNumber string, newOwner string) error {
	fruit, err := s.Queryfruit(ctx, fruitNumber)

	if err != nil {
		return err
	}

	ownerHistory := fruit.Owner
	ownerHistory += "-" 
	ownerHistory += newOwner
	fruit.Owner = ownerHistory

	fruitAsBytes, _ := json.Marshal(fruit)

	return ctx.GetStub().PutState(fruitNumber, fruitAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabfruit chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabfruit chaincode: %s", err.Error())
	}
}
