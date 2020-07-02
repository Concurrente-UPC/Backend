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
	cols = 100
	rows = 10
	K    = 10
)

//Edad
//Sexo
//Insuf_resp
//Neumonia
//Hipertension
//Asma
//Obesidad
//Diabetes
//Enf_cardiacas
//Diagnostico

type Persona_GrupoRiesgo struct {
	Edad          float64 `json:"edad"`
	Sexo          float64 `json:"sexo"`
	Insuf_resp    float64 `json:"insuf_resp"`
	Neumonia      float64 `json:"neumonia"`
	Hipertension  float64 `json:"hipertension"`
	Asma          float64 `json:"asma"`
	Obesidad      float64 `json:"obesidad"`
	Diabetes      float64 `json:"diabetes"`
	Enf_cardiacas float64 `json:"enf_cardiacas"`
	Diagnostico   float64 `json:"diagnostico"`
}

var listPersonas_GR [cols]Persona_GrupoRiesgo
var distancias_GR [cols]float64
var vecinos_GR [K]Persona_GrupoRiesgo
var prueba Persona_GrupoRiesgo

func FloatToString(input_num float64) string {
	// para convertir float a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func leer_dataset_gruporiesgo() {

	csvFile, err := os.Open("Datasets/dataset_gruporiesgo.csv")
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
		age, _ := strconv.ParseFloat(line[0], 64)
		sex, _ := strconv.ParseFloat(line[1], 64)
		insufresp, _ := strconv.ParseFloat(line[2], 64)
		neumon, _ := strconv.ParseFloat(line[3], 64)
		hipert, _ := strconv.ParseFloat(line[4], 64)
		asmha, _ := strconv.ParseFloat(line[5], 64)
		obesid, _ := strconv.ParseFloat(line[6], 64)
		diabet, _ := strconv.ParseFloat(line[7], 64)
		enfcardi, _ := strconv.ParseFloat(line[8], 64)
		diagnos, _ := strconv.ParseFloat(line[9], 64)

		per := Persona_GrupoRiesgo{
			Edad:          age,
			Sexo:          sex,
			Insuf_resp:    insufresp,
			Neumonia:      neumon,
			Hipertension:  hipert,
			Asma:          asmha,
			Obesidad:      obesid,
			Diabetes:      diabet,
			Enf_cardiacas: enfcardi,
			Diagnostico:   diagnos,
		}
		listPersonas_GR[i] = per
		fmt.Println(FloatToString(per.Edad) + " " + FloatToString(per.Sexo) +
			" " + FloatToString(per.Insuf_resp) + " " + FloatToString(per.Neumonia) +
			" " + FloatToString(per.Hipertension) +
			" " + FloatToString(per.Asma) + " " + FloatToString(per.Obesidad) +
			" " + FloatToString(per.Diabetes) + " " + FloatToString(per.Enf_cardiacas) +
			" " + FloatToString(per.Diagnostico))
	}
}

func dist_eucl_gruporiesgo(per1 Persona_GrupoRiesgo, per2 Persona_GrupoRiesgo, i int) {
	distancia := 0.0
	distancia = math.Pow((per1.Edad-per2.Edad)*10, 2) +
		math.Pow(per1.Sexo-per2.Sexo, 2) +
		math.Pow(per1.Insuf_resp-per2.Insuf_resp, 2) +
		math.Pow(per1.Neumonia-per2.Neumonia, 2) +
		math.Pow(per1.Hipertension-per2.Hipertension, 2) +
		math.Pow(per1.Asma-per2.Asma, 2) +
		math.Pow(per1.Obesidad-per2.Obesidad, 2) +
		math.Pow(per1.Diabetes-per2.Diabetes, 2) +
		math.Pow(per1.Enf_cardiacas-per2.Enf_cardiacas, 2)

	distancias_GR[i] = distancia
	println(distancia)
}

type ResultData struct {
	Result int `json:"result"`
}

//KNN_gruporiesgo : Encuetra los vecinos_GR mas cercanos
func KNN_gruporiesgo(prueba Persona_GrupoRiesgo) {
	ch := make(chan int, len(listPersonas_GR))
	var wg sync.WaitGroup
	wg.Add(len(listPersonas_GR))
	var index [cols]int
	for i := 0; i < len(listPersonas_GR); i++ {
		ch <- i
		go func() {
			p := <-ch
			dist_eucl_gruporiesgo(prueba, listPersonas_GR[p], p)
			wg.Done()
		}()

		index[i] = i
	}
	wg.Wait()
	tmp := 0.0
	tmp2 := 0
	for x := 0; x < len(distancias_GR); x++ {
		for y := 0; y < len(distancias_GR); y++ {
			if distancias_GR[x] < distancias_GR[y] {
				tmp = distancias_GR[y]
				distancias_GR[y] = distancias_GR[x]
				distancias_GR[x] = tmp
				tmp2 = index[y]
				index[y] = index[x]
				index[x] = tmp2
			}
		}
	}
	for i := 0; i < len(vecinos_GR); i++ {
		vecinos_GR[i] = listPersonas_GR[index[i]]
	}

}

func GrupoRiesgoEndPoint(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var person Persona_GrupoRiesgo
	_ = json.NewDecoder(request.Body).Decode(&person)
	person.Edad = (person.Edad - 96) / (96 - 19)
	print(person.Edad)
	//----------------
	//---algoritmo----
	KNN_gruporiesgo(person)
	var result = definir_gruporiesgo()
	var resultData ResultData
	resultData.Result = result
	//dato devuelto
	json.NewEncoder(w).Encode(resultData)
}

//metodo que predice si pertence al grupo de riesgo de los vecinos_GR mas cercanos
func definir_gruporiesgo() int {
	prediccion := 0
	var contadorM int
	var contadorB int
	for _, vecino := range vecinos_GR {
		if vecino.Diagnostico == 1 {
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
	leer_dataset_gruporiesgo()
	prueba := Persona_GrupoRiesgo{0.5844, 1, 0, 0, 0, 0, 0, 0, 0, 1}
	KNN_gruporiesgo(prueba)
	fmt.Println(definir_gruporiesgo())

	router := mux.NewRouter()
	//endpoints
	router.HandleFunc("/KNN_gruporiesgo", GrupoRiesgoEndPoint).Methods("POST", "OPTIONS")
	http.ListenAndServe(":3000", router)
	log.Fatal(http.ListenAndServe(":3000", router))

}
