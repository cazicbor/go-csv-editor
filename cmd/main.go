package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
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

func writeToChannel(c chan [][]string, data [][]string, workerID int) {
	c <- data
	close(c)

}

func main() {
	records := readCsvFile("../dataset.csv")

	c := make(chan [][]string)

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)


	go writeToChannel(c, records, i)
	time.Sleep(1 * time.Second)

	res := <-c

	fmt.Println("Read:", res)
	time.Sleep(1 * time.Second)

	_, ok := <-c //check if channel is open or not
	if ok {
		fmt.Println("Channel is open!")
	} else {
		fmt.Println("Channel is closed!")
	}
}

//but ici : créer un channel qui va dispatcher le travail entre les workers
//pourquoi ici il faut un channel : car sinon les goroutines vont s'éxecuter aléatoirement => les données du fichier vont être traitées aléatoirement
//il faut typer le channel pour assurer qu'un seul type de données y passe
//il faut donner un sens au channel, pour que les données ne puissent circuler que dans un sens, pour ne pas risquer d'envoyer les données dans le sens inverse (et avoir des résultats non voulus)
//une fois que toutes les données sont passées, il ne faut pas oublier de fermer le channel avec la fonction close()
//quand on transmet une donnée x à un channel c (c <- x), les deux doivent être du même type
