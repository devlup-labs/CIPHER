package scheduler

import (
	"sync"
	"cipher/internal/content/core"
)

type ChunkTask struct {
	Index    int
	ChunkID  core.ChunkID
	Attempts int
}

type ChunkQueue struct {
	tasks []ChunkTask
	mu    sync.Mutex
}

func NewChunkQueue(tasks []ChunkTask) *ChunkQueue {
	return &ChunkQueue{
		tasks: tasks,
	}
}

func (q *ChunkQueue) Next() (ChunkTask, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return ChunkTask{}, false
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task, true
}

func (q *ChunkQueue) Push(task ChunkTask) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
}

func (q *ChunkQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.tasks)
}
