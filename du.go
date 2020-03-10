package main
import (
	"path/filepath"
	"sync"
	"os"
	"fmt"
)
func main() {
	roots := os.Args[1:]
	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()
	var nfiles, nbytes int64
	loop:
	for {
		select {
			case <-done:
				// drain fileSizes to allow existing goroutines to finish.
				for range fileSizes {
				}
				return
			case size, ok := <-fileSizes:
				if !ok {
					break loop // fileSizes was closed
				}
				nfiles++
				nbytes += size
			}
		}
	printDiskUsage(nfiles, nbytes) // final totals
	}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f MB\n", nfiles, float64(nbytes)/1e6)
	}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			fmt.Println(subdir)
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
			//fmt.Println(filepath.Join(dir, entry.Name()), entry.Size()/1e3, "KB")
		}
	}
}

// concurrency-limiting counting semaphore
var sema = make(chan struct{}, 20)
func dirents(dir string) []os.FileInfo {
	select {
		case sema <- struct{}{}: // acquire token
		case <-done:
			return nil // cancelled
		}
	defer func() { <-sema }() // release token
	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()
	// 0 => no limit; read all entries
	entries, err := f.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
	}
	return entries
}

var done = make(chan struct{})
func cancelled() bool {
	select {
		case <-done:
			return true
		default:
			return false
		}
}