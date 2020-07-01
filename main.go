package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
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
//Contact_no
//Contact_yes
//Results

type Persona struct {
	Fever                   float64 `json:"fever"`
	Tiredness               float64 `json:"tiredness"`
	Dry_cough               float64 `json:"dry_cough"`
	Difficulty_in_breathing float64 `json:"difficulty_in_breathing"`
	Sore_throat             float64 `json:"sore_throat"`
	No_sintomas             float64 `json:"no_sintomas"`
	Pains                   float64 `json:"pains"`
	Nasal_congestion        float64 `json:"nasal_congestion"`
	Runny_nose              float64 `json:"runny_nose"`
	Diarrhea                float64 `json:"diarrhea"`
	No_other_sintomas       float64 `json:"no_other_sintomas"`
	Edad_0_9                float64 `json:"edad_0_9"`
	Edad_10_19              float64 `json:"edad_10_19"`
	Edad_20_24              float64 `json:"edad_20_24"`
	Edad_25_59              float64 `json:"edad_25_59"`
	Edad_60_more            float64 `json:"edad_60_more"`
	Female                  float64 `json:"female"`
	Male                    float64 `json:"male"`
	Contact_dk              float64 `json:"contact_dk"`
	Contact_no              float64 `json:"contact_no"`
	Contact_yes             float64 `json:"contact_yes"`
	Results                 float64 `json:"results"`
}

var listPersonas [cols]Persona
var distancias [cols]float64
var vecinos [K]Persona
var prueba Persona

func FloatToString(input_num float64) string {
	// para convertir float a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func leer_dataset() {

	csvFile, err := os.Open("Datasets/dataset_deteccion.csv")
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

		per := Persona{
			Fever:                   fiebre,
			Tiredness:               cansancio,
			Dry_cough:               tos_seca,
			Difficulty_in_breathing: dificultad_respiratoria,
			Sore_throat:             dolor_garganta,
			No_sintomas:             sin_sintomas,
			Pains:                   dolores,
			Nasal_congestion:        congestion_nasal,
			Runny_nose:              nariz_moqueo,
			Diarrhea:                diarrea,
			No_other_sintomas:       sin_otros_sintomas,
			Edad_0_9:                edad_0a9,
			Edad_10_19:              edad_10a19,
			Edad_20_24:              edad_20a24,
			Edad_25_59:              edad_25a59,
			Edad_60_more:            edad_60aMas,
			Female:                  mujer,
			Male:                    hombre,
			Contact_dk:              contacto_nosabe,
			Contact_no:              contacto_no,
			Contact_yes:             contacto_si,
			Results:                 resultados,
		}
		listPersonas[i] = per
		fmt.Println(FloatToString(per.Fever) + " " + FloatToString(per.Tiredness) +
			" " + FloatToString(per.Dry_cough) + " " + FloatToString(per.Difficulty_in_breathing) +
			" " + FloatToString(per.Sore_throat) + " " + FloatToString(per.No_sintomas) +
			" " + FloatToString(per.Pains) + " " + FloatToString(per.Nasal_congestion) +
			" " + FloatToString(per.Runny_nose) + " " + FloatToString(per.Diarrhea) +
			" " + FloatToString(per.No_other_sintomas) + " " + FloatToString(per.Edad_0_9) +
			" " + FloatToString(per.Edad_10_19) + " " + FloatToString(per.Edad_20_24) +
			" " + FloatToString(per.Edad_25_59) + " " + FloatToString(per.Edad_60_more) +
			" " + FloatToString(per.Female) + " " + FloatToString(per.Male) +
			" " + FloatToString(per.Contact_dk) + " " + FloatToString(per.Contact_no) +
			" " + FloatToString(per.Contact_yes) + " " + FloatToString(per.Results))
	}
}

func dist_eucl(per1 Persona, per2 Persona, i int) {
	distancia := 0.0
	distancia = math.Pow(per1.Fever-per2.Fever, 2) +
		math.Pow(per1.Tiredness-per2.Tiredness, 2) +
		math.Pow(per1.Dry_cough-per2.Dry_cough, 2) +
		math.Pow(per1.Difficulty_in_breathing-per2.Difficulty_in_breathing, 2) +
		math.Pow(per1.Sore_throat-per2.Sore_throat, 2) +
		math.Pow(per1.No_sintomas-per2.No_sintomas, 2) +
		math.Pow(per1.Pains-per2.Pains, 2) +
		math.Pow(per1.Nasal_congestion-per2.Nasal_congestion, 2) +
		math.Pow(per1.Runny_nose-per2.Runny_nose, 2) +
		math.Pow(per1.Diarrhea-per2.Diarrhea, 2) +
		math.Pow(per1.No_other_sintomas-per2.No_other_sintomas, 2) +
		math.Pow(per1.Edad_0_9-per2.Edad_0_9, 2) +
		math.Pow(per1.Edad_10_19-per2.Edad_10_19, 2) +
		math.Pow(per1.Edad_20_24-per2.Edad_20_24, 2) +
		math.Pow(per1.Edad_25_59-per2.Edad_25_59, 2) +
		math.Pow(per1.Edad_60_more-per2.Edad_60_more, 2) +
		math.Pow(per1.Female-per2.Female, 2) +
		math.Pow(per1.Male-per2.Male, 2) +
		math.Pow(per1.Contact_dk-per2.Contact_dk, 2) +
		math.Pow(per1.Contact_no-per2.Contact_no, 2) +
		math.Pow(per1.Contact_no-per2.Contact_yes, 2)

	distancias[i] = distancia
}

type ResultData struct {
	Result int `json:"result"`
}

//KNN : Encuetra los vecinos mas cercanos
func KNN(prueba Persona) {
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

func PredecirEndPoint(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var person Persona
	_ = json.NewDecoder(request.Body).Decode(&person)

	//----------------
	//--algoritmo---
	KNN(person)
	var result = predecir()
	var resultData ResultData
	resultData.Result = result
	//dato devuelto
	json.NewEncoder(w).Encode(resultData)
}

//metodo que predice el tipo de cancer dependeiendo de los vecinos mas cercanos
func predecir() int {
	prediccion := 0
	var contadorM int
	var contadorB int
	for _, vecino := range vecinos {
		if vecino.Results == 1 {
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
	prueba := Persona{1.000000, 1.000000, 1.000000, 1.000000, 0.000000, 0.000000, 0.000000,
		1.000000, 1.000000, 1.000000, 0.000000, 0.000000, 0.000000, 0.000000, 0.000000,
		1.000000, 1.000000, 0.000000, 0.000000, 0.000000, 1.000000, 1.000000}
	KNN(prueba)
	fmt.Println(predecir())

	router := mux.NewRouter()
	// endpoints
	router.HandleFunc("/KNN", PredecirEndPoint).Methods("POST", "OPTIONS")
	http.ListenAndServe(":3000", router)
	log.Fatal(http.ListenAndServe(":3000", router))

}
