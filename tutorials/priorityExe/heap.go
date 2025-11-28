package main

type JobPQ []Job

func (pq JobPQ) Len() int{
	return len(pq)
}

func (pq JobPQ) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq JobPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *JobPQ) Push(x interface{}){
	*pq = append(*pq, x.(Job))
}

func (pq *JobPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n - 1]
	*pq = old[:n - 1]
	return item
}
