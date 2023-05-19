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

		paginas = criar_paginas(db_path, paginas_ativas, quant_bytes_por_pagina)

		fmt.Println("Banco de Dados Encontrado...")
	}

	return esp_livre_paginas, paginas
}

func criar_db(db_path string, quant_paginas int, quant_bytes_por_pagina int) []int {
	var esp_livre_paginas []int
	var empty_string []string

	fmt.Println("Criando Banco de Dados...")
	os.Mkdir(db_path, 0755)

	vetor_ocup := make([]int, quant_bytes_por_pagina)
	for idx, _ := range vetor_ocup {
		vetor_ocup[idx] = -1
	}

	for i := 0; i < quant_paginas; i++ {
		gravar_conteudo_pagina(db_path, i, vetor_ocup, empty_string)
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

func gravar_conteudo_pagina(db_path string, pagina_id int, slots []int, registros []string) {
	// fmt.Println("slots", slots)

	path_comp := db_path + "/" + strconv.Itoa(pagina_id) + ".txt"

	string_vetor_slots := strings.Join(strings.Fields(fmt.Sprint(slots)), " ")

	string_vetor_comp := string_vetor_slots[1 : len(string_vetor_slots)-1]
	string_vetor_comp += "\n"

	for idx, _ := range registros {
		string_vetor_comp += registros[idx]
		string_vetor_comp += "\n"
	}

	err := ioutil.WriteFile(path_comp, []byte(string_vetor_comp), 0644)

	if err != nil {
		os.Exit(0)
	}
}

func ler_conteudo_pagina(db_path string, pagina_id int, quant_bytes_por_pagina int) ([]int, []string) {
	var ocupacao []int

	registros := make([]string, quant_bytes_por_pagina)

	path_pg := db_path + "/" + strconv.Itoa(pagina_id) + ".txt"

	vetor_bin, _ := ioutil.ReadFile(path_pg)

	spl := func(c rune) bool {
		return c == ' ' || c == '\n'
	}

	vetor_string := string(vetor_bin)
	valores := strings.FieldsFunc(vetor_string, spl)

	for i := 0; i < quant_bytes_por_pagina; i++ {
		val_ocup, _ := strconv.Atoi(valores[i])
		ocupacao = append(ocupacao, val_ocup)
	}

	// fmt.Println(len(valores))
	// fmt.Println(valores)

	for i := 0; i < len(valores)-quant_bytes_por_pagina; i++ {
		registros[i] = valores[quant_bytes_por_pagina+i]
	}

	return ocupacao, registros
}

func ler_registros_mem(db_path string, pagina int, quant_bytes_por_pagina int) []*Registro {

	var registros []*Registro
	var ocupacao []int
	var valores_registros []string
	valor_registro := ""
	tamanho_registro := 0

	ocupacao, valores_registros = ler_conteudo_pagina(db_path, pagina, quant_bytes_por_pagina)

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

func criar_paginas(db_path string, paginas_ativas []int, quant_bytes_por_pagina int) []*Pagina {

	var paginas []*Pagina

	for idx, _ := range paginas_ativas {
		registros_pg := ler_registros_mem(db_path, paginas_ativas[idx], quant_bytes_por_pagina)

		pagina := Pagina{
			id:        paginas_ativas[idx],
			registros: registros_pg,
			prox:      nil,
		}

		paginas = append(paginas, &pagina)
	}

	for i := 0; i < (len(paginas) - 1); i++ {
		paginas[i].prox = paginas[i+1]
	}

	return paginas
}

func insert(db_path string, paginas_utilizadas *[]*Pagina, esp_livre_paginas []int, quant_bytes_por_pagina int) {
	var registro_string string

	for {
		fmt.Println("Digite o Registro: ")
		fmt.Scan(&registro_string)

		if !(len(registro_string) > quant_bytes_por_pagina) {
			break
		} else {
			registro_string = ""
			fmt.Println("O Tamanho Máximo do Registro é de %d Bytes", quant_bytes_por_pagina)
		}
	}

	quant_vazios := 0
	slot_gravacao := -1
	pagina_gravacao := -1

	for _, pg_ativa := range *paginas_utilizadas {
		if esp_livre_paginas[pg_ativa.id] >= len(registro_string) {
			ocupacao_slots, _ := ler_conteudo_pagina(db_path, pg_ativa.id, quant_bytes_por_pagina)

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

				pagina_gravacao = pg_ativa.id

				pg_ativa.registros = append(pg_ativa.registros, &registro)

				break
			}

		}

	}

	if slot_gravacao == -1 {
		for idx, _ := range esp_livre_paginas {
			if esp_livre_paginas[idx] == quant_bytes_por_pagina {
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

				pagina_gravacao = idx

				if len(*paginas_utilizadas) != 0 {
					(*paginas_utilizadas)[len(*paginas_utilizadas)-1].prox = &pagina
				}

				*paginas_utilizadas = append(*paginas_utilizadas, &pagina)

				break
			}
		}
	}

	if slot_gravacao == -1 {

		fmt.Println("Banco de Dados Cheio!!")

	} else {

		vetorSlots, vetorRegistros := ler_conteudo_pagina(db_path, pagina_gravacao, quant_bytes_por_pagina)

		vetor_registro_string := strings.Split(registro_string, "")

		for i := 0; i < len(registro_string); i++ {
			vetorSlots[slot_gravacao+i] = slot_gravacao
			vetorRegistros[slot_gravacao+i] = vetor_registro_string[i]
		}

		esp_livre_paginas[pagina_gravacao] -= len(registro_string)

		gravar_conteudo_pagina(db_path, pagina_gravacao, vetorSlots, vetorRegistros)
		gravar_esp_livre_paginas(db_path, esp_livre_paginas)

	}

}

func scan(paginas_ativas *[]*Pagina) ([]*Registro, string) {
	var RegistrosEncontrados []*Registro
	var log string

	if len(*paginas_ativas) == 0 {
		log = "Nenhum registro encontrado"
		return RegistrosEncontrados, log
	}

	paginaAtual := (*paginas_ativas)[0]

	for {
		if paginaAtual == nil {
			break
		}
		for _, registro := range (*paginaAtual).registros {
			RegistrosEncontrados = append(RegistrosEncontrados, registro)
		}
		paginaAtual = paginaAtual.prox
	}

	log = "Número de registros encontrados: " + strconv.Itoa(len(RegistrosEncontrados))

	return RegistrosEncontrados, log
}

func seek(paginas_ativas *[]*Pagina, valor_a_pesquisar string) ([]*Registro, string) {
	var RegistrosAretornar []*Registro
	var log string

	if len(*paginas_ativas) == 0 {
		log = "Nenhum registro encontrado"
		return RegistrosAretornar, log
	}

	paginaAtual := (*paginas_ativas)[0]

	for {
		if paginaAtual == nil {
			break
		}
		for _, registro := range (*paginaAtual).registros {
			if registro.conteudo == valor_a_pesquisar {
				RegistrosAretornar = append(RegistrosAretornar, registro)
			}
		}
		paginaAtual = paginaAtual.prox
	}

	log = "Número de registros encontrados: " + strconv.Itoa(len(RegistrosAretornar))

	return RegistrosAretornar, log
}

func delete(db_path string, paginas_ativas *[]*Pagina, espaco_livre_paginas []int, valor_a_pesquisar string, quant_bytes_por_pagina int) string {
	var registrosAdeletar []*Registro
	var log string

	registrosAdeletar, log = seek(paginas_ativas, valor_a_pesquisar)

	if log == "Nenhum registro encontrado" {
		return log
	}

	for _, registro := range registrosAdeletar {
		fmt.Println("O REGISTRO É ====")
		fmt.Println(registro)
		for indexPagina, pagina := range *paginas_ativas {
			fmt.Println("ESTOU NA PAGINA DE INDEX ===")
			fmt.Println(indexPagina)
			fmt.Println("=======================")
			fmt.Println(pagina)
			if (*pagina).id == registro.pagina_id {

				fmt.Println("ENCONTREI A PAGINA DO REGISTRO ====")
				fmt.Println("É A PAGINA DE INDEX ===")
				fmt.Println(indexPagina)

				index := 0
				tamanhoLista := len((*pagina).registros)
				fmt.Println("PEGUEI A LISTA DE REGISTROS DA PÁGINA")
				for {
					if index >= tamanhoLista {
						break
					}
					fmt.Println("AAAAAAAAA")
					if (*pagina).registros[index].slot == registro.slot {
						fmt.Println("CONTEUDO DO REGISTRO ENCONTRADO")
						fmt.Println((*pagina).registros[index].conteudo)
						fmt.Println("=================")
						(*pagina).registros[index] = (*pagina).registros[len((*pagina).registros)-1]
						(*pagina).registros = (*pagina).registros[:len((*pagina).registros)-1]
						fmt.Println("aaaaaaaskdjhasjdhaksjdaaa")

					} else {
						index += 1
					}

					tamanhoLista = len((*pagina).registros)
					fmt.Println("TAMANHO DA LISTA DE REGI APÓS A REMOÇÃO")
					fmt.Println(tamanhoLista)
					fmt.Println("==============================")

				}

				fmt.Println("ESPAÇO LIVRE ANTES ")
				fmt.Println(espaco_livre_paginas[indexPagina])
				fmt.Println("TAMANHO DO REGISTRO ")
				fmt.Println(registro.tamanho)
				fmt.Println("ESPAÇO LIVRE DEPOIS")
				espaco_livre_paginas[indexPagina] += registro.tamanho
				fmt.Println(espaco_livre_paginas[indexPagina])

				gravar_esp_livre_paginas(db_path, espaco_livre_paginas)

				if espaco_livre_paginas[indexPagina] == quant_bytes_por_pagina {
					/* (*paginas_ativas)[indexPagina] = (*paginas_ativas)[len(*paginas_ativas)]
					*paginas_ativas = (*paginas_ativas)[:len(*(paginas_ativas))] */

					(*paginas_ativas)[indexPagina] = (*paginas_ativas)[len(*paginas_ativas)-1]
					*paginas_ativas = (*paginas_ativas)[:len(*(paginas_ativas))-1]

					fmt.Println("TAMANHO DA LISTA DE PAGINAS ATIVAS APOS A REMOÇÃO")
					fmt.Println(len(*paginas_ativas))

					/* paginaAtual := (*paginas_ativas)[0]

					for {
						if paginaAtual == nil {
							break
						}
						if (*paginaAtual).prox == pagina {

							(*paginaAtual).prox = pagina.prox
							break
						}
						paginaAtual = paginaAtual.prox
					} */

					fmt.Println("VOU LER A PÁGINA DO DISCO")

					vetorSlots, vetorRegistros := ler_conteudo_pagina(db_path, (*pagina).id, quant_bytes_por_pagina)
					fmt.Println("=============================")
					fmt.Println("PRINTAR O VETOR DE SLOTS")
					for _, i := range vetorSlots {
						fmt.Println(i)
					}

					for i := 0; i < registro.tamanho; i++ {
						vetorSlots[registro.slot+i] = -1
					}

					fmt.Println("Vou escrever na página")

					gravar_conteudo_pagina(db_path, (*pagina).id, vetorSlots, vetorRegistros)

				}

				break
			}
		}
	}

	return "Regristros deletados com sucesso"
}

func main() {
	var esp_livre_paginas []int
	var paginas_utilizadas []*Pagina
	//vector created to test scan

	DB_PATH := "db"
	QUANT_PAGINAS := 20
	QUANT_BYTES_POR_PAGINA := 5

	// fmt.Println(criar_paginas(QUANT_PAGINAS))

	// fmt.Println("Inicializando...")

	esp_livre_paginas, paginas_utilizadas = conectar_db(DB_PATH, QUANT_PAGINAS, QUANT_BYTES_POR_PAGINA)

	// insert(DB_PATH, &paginas_utilizadas, esp_livre_paginas, QUANT_BYTES_POR_PAGINA)
	// insert(DB_PATH, &paginas_utilizadas, esp_livre_paginas, QUANT_BYTES_POR_PAGINA)
	// insert(DB_PATH, &paginas_utilizadas, esp_livre_paginas, QUANT_BYTES_POR_PAGINA)
	// insert(DB_PATH, &paginas_utilizadas, esp_livre_paginas, QUANT_BYTES_POR_PAGINA)

	fmt.Println("Fazendo scan")

	registros, log := scan(&paginas_utilizadas)

	fmt.Println(log)

	for _, regi := range registros {
		fmt.Println(regi)
	}

	fmt.Println("Fazendo o delete")

	delete(DB_PATH, &paginas_utilizadas, esp_livre_paginas, "haha", QUANT_BYTES_POR_PAGINA)

	fmt.Println("Fazendo scan após delete")

	registros2, log2 := scan(&paginas_utilizadas)

	fmt.Print(log2)

	insert(DB_PATH, &paginas_utilizadas, esp_livre_paginas, QUANT_BYTES_POR_PAGINA)

	for i := range registros2 {
		fmt.Println(i)
	}

	/* fmt.Println(esp_livre_paginas)

	for _, i := range paginas_utilizadas {
		fmt.Println(i.id)

		for _, reg := range i.registros {
			fmt.Println(reg.conteudo)
		}
	} */

	for _, i := range paginas_utilizadas {
		fmt.Println(i.id)

		for _, reg := range i.registros {
			fmt.Print(len(reg.conteudo), " ")
			fmt.Println(reg.conteudo)
		}
	}

}
