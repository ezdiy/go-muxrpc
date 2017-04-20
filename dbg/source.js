var MRPC = require('muxrpc')
var pull = require('pull-stream')
var toPull = require('stream-to-pull-stream')
var pullNext = require('pull-next')

var api = {
  stuff: 'source'
}

var rpc = MRPC(api, null)()
var a = rpc.createStream()
pull(a, toPull.sink(process.stdout))
pull(toPull.source(process.stdin), a)

pull(rpc.stuff({test:1}), pull.drain(console.error))
pull(rpc.stuff({test:2}), pull.drain(console.error))

var state = false;
pullNext(function () {
    if (!state) {
	state = true;
	return pull(rpc.stuff({test:1}), pull.through(console.error))
    } else {
	return pull(rpc.stuff({test:2}), pull.through(console.error))
    }
})

