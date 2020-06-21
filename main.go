package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

const (
	cols = 300
	rows = 8
	K    = 10
)

//Tiredness
//Dry-Cough
//Difficulty-in-Breathing
//Sore-Throat
//None_Sympton
//Pains
//Nasal-Congestion
//Runny-Nose
//Diarrhea
//None_Experiencing
//Age_0-9
//Age_10-19
//Age_20-24
//Age_25-59
//Age_60+
//Gender_Female
//Gender_Male
//Contact_Dont-Know
//Contact_No
//Contact_Yes
//Results

type persona struct {
	fever                   float64
	tiredness               float64
	dry_cough               float64
	difficulty_in_breathing float64
	sore_throat             float64
	no_sintomas             float64
	pains                   float64
	nasal_congestion        float64
	runny_nose              float64
	diarrhea                float64
	no_other_sintomas       float64
	edad_0_9                float64
	edad_10_19              float64
	edad_20_24              float64
	edad_25_59              float64
	edad_60_more            float64
	female                  float64
	male                    float64
	contact_dk              float64
	contact_no              float64
	contact_yes             float64
	results                 float64
}

var listPersonas [cols]persona
var distancias [cols]float64
var vecinos [K]persona
var prueba persona

func FloatToString(input_num float64) string {
	// para convertir float a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func leer_dataset() {

	csvFile, err := os.Open("dataset_N.csv")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for i, line := range csvLines {
		fiebre, _ := strconv.ParseFloat(line[0], 64)
		cansancio, _ := strconv.ParseFloat(line[1], 64)
		tos_seca, _ := strconv.ParseFloat(line[2], 64)
		dificultad_respiratoria, _ := strconv.ParseFloat(line[3], 64)
		dolor_garganta, _ := strconv.ParseFloat(line[4], 64)
		sin_sintomas, _ := strconv.ParseFloat(line[5], 64)
		dolores, _ := strconv.ParseFloat(line[6], 64)
		congestion_nasal, _ := strconv.ParseFloat(line[7], 64)
		nariz_moqueo, _ := strconv.ParseFloat(line[8], 64)
		diarrea, _ := strconv.ParseFloat(line[9], 64)
		sin_otros_sintomas, _ := strconv.ParseFloat(line[10], 64)
		edad_0a9, _ := strconv.ParseFloat(line[11], 64)
		edad_10a19, _ := strconv.ParseFloat(line[12], 64)
		edad_20a24, _ := strconv.ParseFloat(line[13], 64)
		edad_25a59, _ := strconv.ParseFloat(line[14], 64)
		edad_60aMas, _ := strconv.ParseFloat(line[15], 64)
		mujer, _ := strconv.ParseFloat(line[16], 64)
		hombre, _ := strconv.ParseFloat(line[17], 64)
		contacto_nosabe, _ := strconv.ParseFloat(line[18], 64)
		contacto_no, _ := strconv.ParseFloat(line[19], 64)
		contacto_si, _ := strconv.ParseFloat(line[20], 64)
		resultados, _ := strconv.ParseFloat(line[21], 64)

		per := persona{
			fever:                   fiebre,
			tiredness:               cansancio,
			dry_cough:               tos_seca,
			difficulty_in_breathing: dificultad_respiratoria,
			sore_throat:             dolor_garganta,
			no_sintomas:             sin_sintomas,
			pains:                   dolores,
			nasal_congestion:        congestion_nasal,
			runny_nose:              nariz_moqueo,
			diarrhea:                diarrea,
			no_other_sintomas:       sin_otros_sintomas,
			edad_0_9:                edad_0a9,
			edad_10_19:              edad_10a19,
			edad_20_24:              edad_20a24,
			edad_25_59:              edad_25a59,
			edad_60_more:            edad_60aMas,
			female:                  mujer,
			male:                    hombre,
			contact_dk:              contacto_nosabe,
			contact_no:              contacto_no,
			contact_yes:             contacto_si,
			results:                 resultados,
		}
		listPersonas[i] = per
		fmt.Println(FloatToString(per.fever) + " " + FloatToString(per.tiredness) +
			" " + FloatToString(per.dry_cough) + " " + FloatToString(per.difficulty_in_breathing) +
			" " + FloatToString(per.sore_throat) + " " + FloatToString(per.no_sintomas) +
			" " + FloatToString(per.pains) + " " + FloatToString(per.nasal_congestion) +
			" " + FloatToString(per.runny_nose) + " " + FloatToString(per.diarrhea) +
			" " + FloatToString(per.no_other_sintomas) + " " + FloatToString(per.edad_0_9) +
			" " + FloatToString(per.edad_10_19) + " " + FloatToString(per.edad_20_24) +
			" " + FloatToString(per.edad_25_59) + " " + FloatToString(per.edad_60_more) +
			" " + FloatToString(per.female) + " " + FloatToString(per.male) +
			" " + FloatToString(per.contact_dk) + " " + FloatToString(per.contact_no) +
			" " + FloatToString(per.contact_yes) + " " + FloatToString(per.results))
	}
}

func dist_eucl(per1 persona, per2 persona, i int) {
	distancia := 0.0
	distancia = math.Pow(per1.fever-per2.fever, 2) +
		math.Pow(per1.tiredness-per2.tiredness, 2) +
		math.Pow(per1.dry_cough-per2.dry_cough, 2) +
		math.Pow(per1.difficulty_in_breathing-per2.difficulty_in_breathing, 2) +
		math.Pow(per1.sore_throat-per2.sore_throat, 2) +
		math.Pow(per1.no_sintomas-per2.no_sintomas, 2) +
		math.Pow(per1.pains-per2.pains, 2) +
		math.Pow(per1.nasal_congestion-per2.nasal_congestion, 2) +
		math.Pow(per1.runny_nose-per2.runny_nose, 2) +
		math.Pow(per1.diarrhea-per2.diarrhea, 2) +
		math.Pow(per1.no_other_sintomas-per2.no_other_sintomas, 2) +
		math.Pow(per1.edad_0_9-per2.edad_0_9, 2) +
		math.Pow(per1.edad_10_19-per2.edad_10_19, 2) +
		math.Pow(per1.edad_20_24-per2.edad_20_24, 2) +
		math.Pow(per1.edad_25_59-per2.edad_25_59, 2) +
		math.Pow(per1.edad_60_more-per2.edad_60_more, 2) +
		math.Pow(per1.female-per2.female, 2) +
		math.Pow(per1.male-per2.male, 2) +
		math.Pow(per1.contact_dk-per2.contact_dk, 2) +
		math.Pow(per1.contact_no-per2.contact_no, 2) +
		math.Pow(per1.contact_no-per2.contact_yes, 2)

	distancias[i] = distancia
}

//KNN : Encuetra los vecinos mas cercanos
func KNN(prueba persona) {
	ch := make(chan int, len(listPersonas))
	var wg sync.WaitGroup
	wg.Add(len(listPersonas))
	var index [cols]int
	for i := 0; i < len(listPersonas); i++ {
		ch <- i
		go func() {
			p := <-ch
			dist_eucl(prueba, listPersonas[p], p)
			wg.Done()
		}()

		index[i] = i
	}
	wg.Wait()
	tmp := 0.0
	tmp2 := 0
	for x := 0; x < len(distancias); x++ {
		for y := 0; y < len(distancias); y++ {
			if distancias[x] < distancias[y] {
				tmp = distancias[y]
				distancias[y] = distancias[x]
				distancias[x] = tmp
				tmp2 = index[y]
				index[y] = index[x]
				index[x] = tmp2
			}
		}
	}
	for i := 0; i < len(vecinos); i++ {
		vecinos[i] = listPersonas[index[i]]
	}

}

//metodo que predice el tipo de cancer dependeiendo de los vecinos mas cercanos
func predecir() int {
	prediccion := 0
	var contadorM int
	var contadorB int
	for _, vecino := range vecinos {
		if vecino.results == 1 {
			contadorM++
		} else {
			contadorB++
		}
	}
	if contadorM > contadorB {
		prediccion = 1
	}
	return prediccion
}

func main() {
	leer_dataset()
	prueba := persona{1.000000, 1.000000, 1.000000, 1.000000, 0.000000, 0.000000, 0.000000,
		1.000000, 1.000000, 1.000000, 0.000000, 0.000000, 0.000000, 0.000000, 0.000000,
		1.000000, 1.000000, 0.000000, 0.000000, 0.000000, 1.000000, 1.000000}
	KNN(prueba)
	fmt.Println(predecir())
}
