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
		fmt.Printf("Au worker %d de faire son taff...", workerID)
		doOK(val)
	}
	res := <-c
	fmt.Println(res)

}

func main() {
	records := readCsvFile("../dataset.csv")

	c := make(chan [][]string, 1000)

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go writeToChannel(c, records, i)
	}

	close(c) //close the channel

	wg.Wait() //wait for all writing to be done

}

//comment faire pour répartir le travail entre les workers (batchs de 250 valeurs)

//but ici : créer un channel qui va dispatcher le travail entre les workers
//pourquoi ici il faut un channel : car sinon les goroutines vont s'éxecuter aléatoirement => les données du fichier vont être traitées aléatoirement
//il faut typer le channel pour assurer qu'un seul type de données y passe
//il faut donner un sens au channel, pour que les données ne puissent circuler que dans un sens, pour ne pas risquer d'envoyer les données dans le sens inverse (et avoir des résultats non voulus)
//une fois que toutes les données sont passées, il ne faut pas oublier de fermer le channel avec la fonction close()
//quand on transmet une donnée x à un channel c (c <- x), les deux doivent être du même type
