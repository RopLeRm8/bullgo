package main

import "bullgo/examples"

func main() {
	examples.BasicError()
}

// Workers -> always working per queue
// Jobs -> All stack up in 1 Queue (the queue is created by the user TOO!)
// Queues -> are the root of bullgo system
