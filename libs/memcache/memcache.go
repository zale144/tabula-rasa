package memcache

import "log"

// read request
type ReadReq struct{
	Key ItemKey
	Resp chan *[]interface{}
}
// write request
type WriteReq struct{
	Key ItemKey
	Val *[]interface{}
	Resp chan bool
}
// remove request
type RemoveReq struct{
	Key ItemKey
	Resp chan bool
}
// clear cache request
type ClearReq struct{
	Resp chan bool
}
// cache map key
type ItemKey struct {
	Name string
	ID   string
	Col  string
}
// read, write, remove, clear request channels
var (
	reads   = make(chan *ReadReq)
	writes  = make(chan *WriteReq)
	removes = make(chan *RemoveReq)
	clears  = make(chan *ClearReq)
)

func MemCache () { // encapsulate the instance of cache map in a goroutine and access it through channels
	memCache := make(map[ItemKey]*[]interface{}) // the map holding all the cache, while values being pointers
	for { // keep checking for requests
		select { // request can either be read, write or remove, not more than one at the same time
		case read := <-reads:               // read operation requested
			read.Resp <- memCache[read.Key] // read from the cache map and return result through channel
		case write := <-writes:             // write operation requested
			memCache[write.Key] = write.Val // write new key/value pair to cache map
			write.Resp <- true              // return bool response
		case remove := <-removes:           // remove operation requested
			delete(memCache, remove.Key)    // delete value for key from cache map
			_, resp := memCache[remove.Key] // check if key still found in map
			remove.Resp <- !resp            // if key no longer found in map, send true, otherwise false
		case clear := <-clears: // clear operation requested
			memCache = make(map[ItemKey]*[]interface{}) // re-instantiate the cache map
			clear.Resp <- len(memCache) == 0 // send true if map is empty now
		}
	}
}
// read item from cache
func ReadFromCache(entities *[]interface{}, name, id, col string)  { // return interface
	read := &ReadReq{ // instantiate a read request
		Key: ItemKey{
			Name: name,
			ID:   id,
			Col:  col,
		},
		Resp: make(chan *[]interface{}),
	}
	reads <- read // send read request to reads channel
	ret := <- read.Resp // pull return value from the channel
	if ret != nil {
		*entities = *ret
		log.Printf("Read '%s' from memcache\n", read.Key)
	}
}
// write item to cache
func WriteToCache(entities *[]interface{}, name, id, col string)  {
	write := &WriteReq{ // instantiate a write request
		Key: ItemKey{
			Name: name,
			ID:   id,
			Col:  col,
		},
		Val: entities,
		Resp: make(chan bool),
	}
	writes <- write // send write request to writes channel
	<- write.Resp   // pull boolean value from the channel
	log.Printf("Written '%s' to memcache\n", write.Key)
}
// remove item from cache
func RemoveFromCache(name, id, col string)  {
	remove := &RemoveReq{ // instantiate a remove request
		Key: ItemKey{
			Name: name,
			ID:   id,
			Col:  col,
		},
		Resp: make(chan bool),
	}
	removes <- remove // send remove request to removes channel
	<- remove.Resp // pull boolean value from the channel
	log.Printf("Removed '%s' from memcache\n", remove.Key)
}
// clear the cache
func ClearCache()  {
	clear := &ClearReq{ // instantiate a clear request
		Resp: make(chan bool),
	}
	clears <- clear // send clear request to clears channel
	<- clear.Resp // pull boolean value from the channel
	log.Println("Memcache cleared")
}