package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Registro struct {
	pagina_id int
	slot      int
	tamanho   int
	conteudo  string
}

type Pagina struct {
	id        int
	registros []*Registro
	prox      *Pagina
	// esp_disp  int
}

func conectar_db(db_path string, quant_paginas int, quant_bytes_por_pagina int) ([]int, []int) {
	// fmt.Println("Verificando a Existência do Banco de Dados...")

	paginas_ativas := []int{}
	var esp_livre_paginas []int

	if _, err := os.Stat(db_path); os.IsNotExist(err) {

		var op_criar_db int

		fmt.Println("Banco De Dados Não Encontrado!")
		fmt.Println("Deseja Criar o Banco de Dados?")
		fmt.Println(" 1 - Sim\n 2 - Não")
		fmt.Scan(&op_criar_db)

		if op_criar_db == 1 {
			esp_livre_paginas = criar_db(db_path, quant_paginas, quant_bytes_por_pagina)
		} else {
			os.Exit(0)
		}
	} else {
		esp_livre_paginas = ler_esp_livre_paginas(db_path)

		for idx, _ := range esp_livre_paginas {
			if esp_livre_paginas[idx] != quant_bytes_por_pagina {
				paginas_ativas = append(paginas_ativas, idx)
			}
		}

		criar_paginas(db_path, paginas_ativas)

		fmt.Println("Banco de Dados Encontrado...")
	}

	return esp_livre_paginas, paginas_ativas
}

func criar_db(db_path string, quant_paginas int, quant_bytes_por_pagina int) []int {
	var esp_livre_paginas []int

	fmt.Println("Criando Banco de Dados...")
	os.Mkdir(db_path, 0755)

	vetor_ocup := make([]int, quant_bytes_por_pagina)
	for idx, _ := range vetor_ocup {
		vetor_ocup[idx] = -1
	}

	string_vetor := strings.Join(strings.Fields(fmt.Sprint(vetor_ocup)), " ")

	for i := 0; i < quant_paginas; i++ {
		path_comp := db_path + "/" + strconv.Itoa(i) + ".txt"

		pagina, _ := os.OpenFile(path_comp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		pagina.WriteString(string_vetor[1 : len(string_vetor)-1])
	}

	esp_livre_paginas = make([]int, quant_paginas)

	for idx, _ := range esp_livre_paginas {
		esp_livre_paginas[idx] = quant_bytes_por_pagina
	}

	gravar_esp_livre_paginas(db_path, esp_livre_paginas)

	fmt.Println("Banco de Dados Criado com Sucesso!")

	return esp_livre_paginas
}

func gravar_esp_livre_paginas(db_path string, esp_livre_paginas []int) {
	path := db_path + "/esp_livre_paginas.txt"

	string_vetor := ""

	for _, elemento := range esp_livre_paginas {
		string_vetor += fmt.Sprintf("%d\n", elemento)
	}

	err := ioutil.WriteFile(path, []byte(string_vetor), 0644)

	if err != nil {
		os.Exit(0)
	}
}

func ler_esp_livre_paginas(db_path string) []int {
	path := db_path + "/esp_livre_paginas.txt"

	vetor_byte, _ := ioutil.ReadFile(path)

	vetor_string := string(vetor_byte)
	linhas := strings.Split(vetor_string, "\n")

	esp_livre_paginas := []int{}

	for _, linha := range linhas {
		if linha != "" {
			var elemento int

			fmt.Sscanf(linha, "%d", &elemento)

			esp_livre_paginas = append(esp_livre_paginas, elemento)
		}
	}

	// fmt.Println("Vetor Lido:",enderecos)

	return esp_livre_paginas

}

func inserir_registro1(db_path string, enderecos []int, quant_paginas int, quant_bytes_por_pagina int) {
	var registro string

	for {
		fmt.Println("Digite o Registro: ")
		fmt.Scan(&registro)

		if !(len(registro) > 5) {
			break
		} else {
			registro = ""
			fmt.Println("O Tamanho Máximo do Registro é de 5 Bytes")
		}
	}

	quant_vazios := 0

	for i := 0; i < len(enderecos); i++ {
		if enderecos[i] == -1 {
			quant_vazios += 1
			if quant_vazios == len(registro) {
				path := db_path + "/" + strconv.Itoa(i/quant_bytes_por_pagina) + ".txt"

				pagina := i / quant_bytes_por_pagina
				slot := (i % quant_bytes_por_pagina) + 1 - len(registro)

				fmt.Println(fmt.Sprintf("Página %d, Slot %d", pagina, slot))

				string_vetor := ""
				for idx, letra := range registro {
					pos_end := pagina*quant_bytes_por_pagina + slot + idx

					// fmt.Println("pos ",pos_end)

					string_vetor += string(letra) + "\n"

					enderecos[pos_end] = slot
				}

				gravar_esp_livre_paginas(db_path, enderecos)
				fmt.Println("enderecos ", enderecos)

				arquivo, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0664)

				arquivo.WriteString(string_vetor)

				// ioutil.WriteFile(path, []byte(string_vetor), 0644)

			}
		} else {
			quant_vazios = 0
		}
	}

	// for _, i := range(a) {
	//   fmt.Println(string(i))
	// }

}

func ler_conteudo_pagina(db_path string, pagina int) ([]int, []string) {
	var ocupacao []int

	path_pg := db_path + "/" + strconv.Itoa(pagina) + ".txt"

	vetor_bin, _ := ioutil.ReadFile(path_pg)

	spl := func(c rune) bool {
		return c == ' ' || c == '\n'
	}

	vetor_string := string(vetor_bin)
	valores := strings.FieldsFunc(vetor_string, spl)

	for i := 0; i < 5; i++ {
		val_ocup, _ := strconv.Atoi(valores[i])
		ocupacao = append(ocupacao, val_ocup)
	}

	registros := valores[5:]

	return ocupacao, registros
}

func ler_registros_mem(db_path string, pagina int) []*Registro {

	var registros []*Registro
	var ocupacao []int
	var valores_registros []string
	valor_registro := ""
	tamanho_registro := 0

	// path_pg := db_path+"/"+strconv.Itoa(pagina)+".txt"

	ocupacao, valores_registros = ler_conteudo_pagina(db_path, pagina)

	fmt.Println("Ocupação ", ocupacao)
	fmt.Println("valores_Registros ", valores_registros)

	for idx, val := range ocupacao {

		if val != -1 {
			valor_registro += string(valores_registros[idx])
			tamanho_registro += 1

			if idx == 4 || ocupacao[idx+1] != val {

				registro := Registro{
					pagina_id: pagina,
					slot:      val,
					tamanho:   tamanho_registro,
					conteudo:  valor_registro,
				}

				registros = append(registros, &registro)

				// fmt.Print("idx ", idx)
				// fmt.Print(", val ", val)
				// fmt.Print(", reg ", registro.conteudo)
				// fmt.Print(", pg ", registro.pagina_id)
				// fmt.Print(", slot ", registro.slot)
				// fmt.Println(", tamanho ", registro.tamanho)

				valor_registro = ""
				tamanho_registro = 0
			}
		}
	}

	return registros

}

func criar_paginas(db_path string, paginas_ativas []int) {

	for _, pg := range paginas_ativas {
		registros_pg := ler_registros_mem(db_path, pg)
		// fmt.Println(pg)
	}

}

func main() {
	var esp_livre_paginas []int
	var lista_paginas_utilizadas []int

	DB_PATH := "db"
	QUANT_PAGINAS := 20
	QUANT_BYTES_POR_PAGINA := 5

	// fmt.Println(criar_paginas(QUANT_PAGINAS))

	// fmt.Println("Inicializando...")

	esp_livre_paginas, lista_paginas_utilizadas = conectar_db(DB_PATH, QUANT_PAGINAS, QUANT_BYTES_POR_PAGINA)

	fmt.Println(esp_livre_paginas, lista_paginas_utilizadas)

	// inserir_registro1(db_path, enderecos, QUANT_PAGINAS, QUANT_BYTES_POR_PAGINA)

	//  fmt.Scan(&i)
	// fmt.Println("Hello, World!", i, len(i))
}
