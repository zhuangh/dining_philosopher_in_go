package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Philosopher struct {
	id  int          // my ID
	rnb *Philosopher // right neighbor

	left_fork chan bool
	/* left_fork: this philosopher's left fork, also her/his left neighbor's right fork.
	 *            First pick the left fork for all the philosopher except the one with ID = maxid */

	maxid int // record the number of philosophers

}

// initialize the philosopher-1
func initPhilosopher(assigned_id int, maxid int, left_philosopher *Philosopher) *Philosopher {
	cur_phil := &Philosopher{assigned_id, nil, make(chan bool, 1), maxid}
	cur_phil.left_fork <- true // everyone has an available left fork
	return cur_phil
}

// the dining routine for each philosopher
func (phils *Philosopher) start() {
	ts := time.Duration(1)
	id := phils.id
	pick_dirct := (phils.id < phils.maxid) // true if the id < maxid, direction

	for {
		fPnt(id)
		fmt.Println("Philosopher", id, "is thinking.")
		time.Sleep((time.Duration(rand.Intn(5))*(1e3/ts)*time.Millisecond + (5e3/ts)*time.Millisecond) / ts)

		if pick_dirct == true {
			// the priority for the Philosopher with id < maxid
			select {
			case <-phils.left_fork:
				// get the left_fork because it is not available
				<-phils.rnb.left_fork
			default:
				// can't get the left_fork because it is not available
				phils.left_fork <- true
				continue
			}
		} else {
			// the priority for the Philosopher with id == maxid
			select {
			case <-phils.rnb.left_fork:
				// get the right fork because it is not available
				<-phils.left_fork
			default:
				// can't get the right fork because it is not available
				phils.rnb.left_fork <- true
				continue
			}
		}

		fPnt(id)
		fmt.Println("Philosopher", id, "picks the forks.")
		time.Sleep((time.Duration(rand.Intn(1))*1e3*time.Millisecond + 1e3*time.Millisecond) / ts)

		fPnt(id)
		fmt.Println("Philosopher", id, "is eating.")

		if pick_dirct {
			phils.left_fork <- true
			phils.rnb.left_fork <- true
		} else {
			phils.rnb.left_fork <- true
			phils.left_fork <- true
		}

		fPnt(id)
		fmt.Println("Philosopher", id, "puts down the forks.")
	}
}

func fPnt(id int) {
    /*
	for i := 1; i <= id; i++ {
		fmt.Print("-")
	}
	*/
}



func main() {

	var num int
	fmt.Println("Input the number of Philosophers")
	fmt.Scanf("%d", &num)

	fmt.Println("------------Preparations----------------------")

	num = int(num)

	if num <= 1 {
		fmt.Println(" The input is wrong")
		return
	}

	fmt.Println("There will be ", num, "Philosophers")

	phils := make([]*Philosopher, num)

	for i := 1; i <= num; i++ {
		phils[i-1] = initPhilosopher(i, num, nil)
		fmt.Println("Create Philosopher", i, " done")
	}

	for i := 1; i < num; i++ {
		phils[i].rnb = phils[i-1]
		fmt.Println("Connecting Philosopher", i, " with right neighbor is Philosopher", phils[i].rnb.id)
	}

	phils[0].rnb = phils[num-1]
	fmt.Println("Connecting a philosopher ", 1 , "with right neighbor is ", phils[0].rnb.id)
	fmt.Println("------------ Finish connecting and start to dining----------------------")

	for i := 0; i < num; i++ {
		fmt.Println("Launch Philosopher", phils[i].id,"to start dining")
		go phils[i].start()

	}

	fmt.Println("------------ Dining----------------------")

	var input string
	fmt.Scanln(&input)

	fmt.Println("done")
}
