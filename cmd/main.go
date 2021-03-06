package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

const WORKERS = 4

func readCsvFile(filePath string) [][]string {
	f, err := os.Open("/home/coddity/projets/TEST_GOROUTINES/data/dataset.csv")
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func doOK(data [][]string) [][]string {
	for i := range data {
		data[i] = append(data[i], "-OK")
	}
	return data
}

func worker(fanout <-chan [][]string, fanin chan<- [][]string, workerID int) { //this func gets data from a channel, and modifies the input string
	defer wg.Done()

	for val := range fanout {
		res := doOK(val)
		fanin <- res
		fmt.Println(res)
	}
}

func main() {
	records := readCsvFile("../dataset.csv")

	fanin := make(chan [][]string, 1000) //converging channel in which we will store the data,

	fanout := make(chan [][]string, 1000) //diverging channel in which we will store the data, which will be distributed to the workers

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go worker(fanout, fanin, i)

	}
	fanout <- records //adding data to the diverging channel

	close(fanout) //close the channel

	wg.Wait() //wait for all writing to be done

}

//but ici : créer un channel qui va dispatcher le travail entre les workers
//pourquoi ici il faut un channel : car sinon les goroutines vont s'éxecuter aléatoirement => les données du fichier vont être traitées aléatoirement
//il faut typer le channel pour assurer qu'un seul type de données y passe
//il faut donner un sens au channel, pour que les données ne puissent circuler que dans un sens, pour ne pas risquer d'envoyer les données dans le sens inverse (et avoir des résultats non voulus)
//une fois que toutes les données sont passées, il ne faut pas oublier de fermer le channel avec la fonction close()
//quand on transmet une donnée x à un channel c (c <- x), les deux doivent être du même type
