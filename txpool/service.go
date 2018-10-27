/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package txpool

import "github.com/it-chain/engine/common"

const SendTransactionsToLeader = "SendTransactionsToLeaderProtocol"

type TxpoolQueryService interface {
	FindUncommittedTransactions() ([]Transaction, error)
}

func filter(vs []Transaction, f func(Transaction) bool) []Transaction {
	vsf := make([]Transaction, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func IsLeader(nodeId common.NodeID, leader Leader) bool {
	if nodeId != leader.Id {
		return false
	}

	return true
}

type EventService interface {
	Publish(topic string, event interface{}) error
}
