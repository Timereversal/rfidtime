// Package sampling select the best sample for each label/tag [criteria : highest rssi ].
// a goroutine will be created for each label, where each goroutine can receive 1 or more samples

package sampling

import (
	"container/heap"
	"fmt"
	"rfidtime/transport"
	"sync"
	"time"
)

type TagHeap []transport.TagInfo

func (h TagHeap) Len() int { return len(h) }

// min
//func (h TagHeap) Less(i, j int) bool { return h[i].RSSI < h[j].RSSI }

// max
func (h TagHeap) Less(i, j int) bool { return h[i].RSSI > h[j].RSSI }
func (h TagHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TagHeap) Push(x any) { *h = append(*h, x.(transport.TagInfo)) }
func (h *TagHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Broker struct {
	StreamList map[string]chan transport.TagInfo
	Wg         sync.WaitGroup
}

func (b *Broker) StreamGenerator(id string, stream chan<- transport.TagInfo) {

	tagInfoList := &TagHeap{}
	fmt.Println(&tagInfoList)
	heap.Init(tagInfoList)
	b.StreamList[id] = make(chan transport.TagInfo)
	b.Wg.Add(1)
	go func() {
		defer b.Wg.Done()
		//for v := range b.StreamList[id] {
		//	fmt.Printf("tag info %+v inside stream id: %d \n", v, id)
		//}
		for {
			select {
			case v := <-b.StreamList[id]:
				//tagInfoList.Push(v)
				heap.Push(tagInfoList, v)
				//fmt.Printf("tag info %+v inside stream id: %s, tagList %+v \n", v, id, *tagInfoList)

			case <-time.After(10 * time.Second):
				fmt.Printf("stream id %d timeout, EPC: %X  \n", id, (*tagInfoList)[0].EPCData)
				stream <- (*tagInfoList)[0]
				//fmt.Printf(, (*h)[0])
				return
			}
		}
		//b.Wg.Done()
	}()
}
