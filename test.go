package main

import (
	"fmt"
	"go-steem/rpc"
	"time"
)

func main() {
	// Instantiate the WebSocket transport.
	t, _ := websocket.NewTransport("ws://localhost:8090")

	// Use the transport to create an RPC client.
	client, _ := rpc.NewClient(t)
	defer client.Close()

	// Call "get_config".
	config, _ := client.Database.GetConfig()

	// Start processing blocks.
	lastBlock := 1800000
	for {
		// Call "get_dynamic_global_properties".
		props, _ := client.Database.GetDynamicGlobalProperties()

		for props.LastIrreversibleBlockNum-lastBlock > 0 {
			// Call "get_block".
			block, _ := client.Database.GetBlock(lastBlock)

			// Process the transactions.
			for _, tx := range block.Transactions {
				for _, op := range tx.Operations {
					switch body := op.Data().(type) {
					// Comment operation.
					case *types.CommentOperation:
						content, _ := client.Database.GetContent(body.Author, body.Permlink)
						fmt.Printf("COMMENT @%v %v\n", content.Author, content.URL)

					// Vote operation.
					case *types.VoteOperation:
						fmt.Printf("VOTE @%v @%v/%v\n", body.Voter, body.Author, body.Permlink)

						// You can add more cases, it depends on what
						// operations you actually need to process.
					}
				}
			}

			lastBlock++
		}

		time.Sleep(time.Duration(config.SteemitBlockInterval) * time.Second)
	}
}
