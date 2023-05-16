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

func conectar_db(db_path string, quant_paginas int, quant_bytes_por_pagina int) ([]int, []*Pagina) {
	// fmt.Println("Verificando a Existência do Banco de Dados...")

	var paginas []*Pagina
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

		paginas = criar_paginas(db_path, paginas_ativas)

		fmt.Println("Banco de Dados Encontrado...")
	}

	return esp_livre_paginas, paginas
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

	return esp_livre_paginas

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

	ocupacao, valores_registros = ler_conteudo_pagina(db_path, pagina)

	for idx, _ := range ocupacao {

		if ocupacao[idx] != -1 {
			valor_registro += string(valores_registros[idx])
			tamanho_registro += 1

			if idx == 4 || ocupacao[idx+1] != ocupacao[idx] {

				registro := Registro{
					pagina_id: pagina,
					slot:      ocupacao[idx],
					tamanho:   tamanho_registro,
					conteudo:  valor_registro,
				}

				registros = append(registros, &registro)

				valor_registro = ""
				tamanho_registro = 0
			}
		}
	}

	return registros

}

func criar_paginas(db_path string, paginas_ativas []int) []*Pagina {

	var paginas []*Pagina

	for idx, _ := range paginas_ativas {
		registros_pg := ler_registros_mem(db_path, paginas_ativas[idx])
		fmt.Println(registros_pg)

		pagina := Pagina{
			id:        paginas_ativas[idx],
			registros: registros_pg,
			prox:      nil,
		}

		paginas = append(paginas, &pagina)
	}

	for i := (len(paginas) - 1); i > 1; i-- {
		paginas[i-1].prox = paginas[i]
	}

	return paginas
}

func inserir_registro(db_path string, paginas_utilizadas *[]*Pagina, esp_livre_paginas []int) {
	var registro_string string

	for {
		fmt.Println("Digite o Registro_string: ")
		fmt.Scan(&registro_string)

		if !(len(registro_string) > 5) {
			break
		} else {
			registro_string = ""
			fmt.Println("O Tamanho Máximo do Registro é de 5 Bytes")
		}
	}

	quant_vazios := 0
	slot_gravacao := -1

	for _, pg_ativa := range *paginas_utilizadas {
		if esp_livre_paginas[pg_ativa.id] >= len(registro_string) {
			ocupacao_slots, _ := ler_conteudo_pagina(db_path, pg_ativa.id)

			for i := 0; i < len(ocupacao_slots); i++ {
				if ocupacao_slots[i] == -1 {
					quant_vazios += 1
					if quant_vazios == len(registro_string) {
						slot_gravacao = i - len(registro_string) + 1
						break
					}
				} else {
					quant_vazios = 0
				}
			}

			if slot_gravacao != -1 {
				registro := Registro{
					pagina_id: pg_ativa.id,
					slot:      slot_gravacao,
					tamanho:   len(registro_string),
					conteudo:  registro_string,
				}

				pg_ativa.registros = append(pg_ativa.registros, &registro)

				fmt.Println("pagina", pg_ativa.id)
				fmt.Println("slot", slot_gravacao)
				break
			}

		}

	}

	if slot_gravacao == -1 {
		for idx, _ := range esp_livre_paginas {
			if esp_livre_paginas[idx] == 5 {
				slot_gravacao = 0

				registro := Registro{
					pagina_id: idx,
					slot:      slot_gravacao,
					tamanho:   len(registro_string),
					conteudo:  registro_string,
				}

				var regs []*Registro

				regs = append(regs, &registro)

				pagina := Pagina{
					id:        idx,
					registros: regs,
					prox:      nil,
				}

				(*paginas_utilizadas)[len(*paginas_utilizadas)-1].prox = &pagina

				*paginas_utilizadas = append(*paginas_utilizadas, &pagina)

				fmt.Println("pagina ", registro.pagina_id)

				break
			}
		}
	}

	if slot_gravacao == -1 {
		fmt.Println("Banco de Dados Cheio!!")
	}

	// for i := 0; i < len(enderecos); i++ {
	// 	if enderecos[i] == -1 {
	// 		quant_vazios += 1
	// 		if quant_vazios == len(registro) {
	// 			path := db_path + "/" + strconv.Itoa(i/quant_bytes_por_pagina) + ".txt"

	// 			pagina := i / quant_bytes_por_pagina
	// 			slot := (i % quant_bytes_por_pagina) + 1 - len(registro)

	// 			fmt.Println(fmt.Sprintf("Página %d, Slot %d", pagina, slot))

	// 			string_vetor := ""
	// 			for idx, letra := range registro {
	// 				pos_end := pagina*quant_bytes_por_pagina + slot + idx

	// 				// fmt.Println("pos ",pos_end)

	// 				string_vetor += string(letra) + "\n"

	// 				enderecos[pos_end] = slot
	// 			}

	// 			gravar_esp_livre_paginas(db_path, enderecos)
	// 			fmt.Println("enderecos ", enderecos)

	// 			arquivo, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0664)

	// 			arquivo.WriteString(string_vetor)

	// 			// ioutil.WriteFile(path, []byte(string_vetor), 0644)

	// 		}
	// 	} else {
	// 		quant_vazios = 0
	// 	}

}

func main() {
	var esp_livre_paginas []int
	var paginas_utilizadas []*Pagina

	DB_PATH := "db"
	QUANT_PAGINAS := 20
	QUANT_BYTES_POR_PAGINA := 5

	// fmt.Println(criar_paginas(QUANT_PAGINAS))

	// fmt.Println("Inicializando...")

	esp_livre_paginas, paginas_utilizadas = conectar_db(DB_PATH, QUANT_PAGINAS, QUANT_BYTES_POR_PAGINA)

	fmt.Println(esp_livre_paginas)

	for _, i := range paginas_utilizadas {
		fmt.Println(i.id)

		for _, reg := range i.registros {
			fmt.Println(reg.conteudo)
		}
	}

	inserir_registro(DB_PATH, &paginas_utilizadas, esp_livre_paginas)

	for _, i := range paginas_utilizadas {
		fmt.Println(i.id)

		for _, reg := range i.registros {
			fmt.Println(reg.conteudo)
		}
	}

}
