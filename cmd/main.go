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

func doOK(val [][]string) {
	for i := range val {
		val[i] = append(val[i], "-OK")
	}
}

func writeToChannel(c chan [][]string, data [][]string, workerID int) {
	defer wg.Done()

	for val := range c {
		fmt.Printf("Au worker %d de traiter ses données...", workerID)
		doOK(val)
	}
	res := <-c
	fmt.Println(res)

}

func readResults(c chan [][]string) {
	for i := range c {
		fmt.Println(i)
	}
}

func main() {
	records := readCsvFile("../dataset.csv")

	fanin := make(chan [][]string, 1000) //converging channel in which we will store the data,

	fanout := make(chan [][]string, 1000) //diverging channel in which we will store the data, which will be distributed to the workers

	// var lastElementPerWorker *int

	// lastElementPerWorker = new(int)

	// *lastElementPerWorker = len(c)

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go writeToChannel(fanout, records, i)

	}

	//ajout de taches dans la queue
	for i := 1; i <= 1000; i++ {
		fanout <- records
	}

	close(fanout) //close the channel

	wg.Wait() //wait for all writing to be done

	readResults(fanin)

}

//comment faire pour répartir le travail entre les workers (batchs de 250 valeurs) ?
//doit-on faire du Fan-In / Fn-Out avec chacque channel ?

//but ici : créer un channel qui va dispatcher le travail entre les workers
//pourquoi ici il faut un channel : car sinon les goroutines vont s'éxecuter aléatoirement => les données du fichier vont être traitées aléatoirement
//il faut typer le channel pour assurer qu'un seul type de données y passe
//il faut donner un sens au channel, pour que les données ne puissent circuler que dans un sens, pour ne pas risquer d'envoyer les données dans le sens inverse (et avoir des résultats non voulus)
//une fois que toutes les données sont passées, il ne faut pas oublier de fermer le channel avec la fonction close()
//quand on transmet une donnée x à un channel c (c <- x), les deux doivent être du même type

//OK CE QU'ON VA FAIRE : ON VA APPLIQUER L'ALGORITHME DE FAN IN : ON VA ASSIGNER UN CHANNEL À CHAQUE WORKER, ET ON VA TRANSMETTRE A CHAQUE WORKER 250 ELEMENTS DU DATASET. CHAQUE WORKER VA DONC AVOIR SES 250
//ELEMENTS TRANSMIS DANS SON CHANNEL PERSO, ON VA AINSI AVOIR UN TRAITEMENT SIMULTANÉ DES 1000 DONNÉES, QU'ON TRANSMETTRA ENSUITE DANS LE MÊME CHANNEL
