# Sistemas Gerenciadores de Bancos de Dados - Trabalho 1

## Especificação do Trabalho
Considere uma organização de páginas (blocos) de tamanho fixo, onde cada página possui
5 bytes. Implemente uma classe para representar um documento de tamanho variável 
pertencente a um SGBD não relacional, onde cada documento ocupa no mínimo 1 e no 
máximo 5 bytes. Note que as páginas (blocos) possuem tamanho fixo, enquanto os 
documentos (registros) possuem tamanhos variáveis. Logo, páginas distintas podem 
armazenar quantidades diferentes de documentos. Cada documento será representado por 
uma tupla do tipo <int page id, int seq, int tam> chamada DID. Implemente também uma 
classe para representar uma página de disco que contém uma coleção de documentos. Cada 
documento será simulado por uma sequência de caracteres. O sistema de armazenamento 
deverá ter um total de 20 páginas. É recomendado que seja criada uma interface para a 
classe página expondo os métodos de manipulação dos documentos (registros).
O sistema de armazenamento deverá implementar as seguintes operações:

### Scan():
faz uma varredura de todas as páginas e, consequentemente, retorna todos os 
documentos existentes nelas.
### Seek([]char):
busca um documento através do seu conteúdo e retorna o DID correspondente 
a ele, representado pela tupla <int page id, int seq, int tam>. Se houver mais de um 
documento com o mesmo conteúdo, a primeira ocorrência deve ser retornada. Caso o 
documento não seja encontrado, apresentar mensagem na tela.
### Delete([]char):
deleta o documento que contém o conteúdo que foi informado como 
parâmetro. Se houver mais de um documento com o mesmo conteúdo, remover a primeira 
ocorrência. Caso a página fique vazia, manipular ponteiro da página anterior para apontar 
para a página posterior à que ficou vazia, e direcionar a página que ficou vazia para o 
diretório de páginas em branco. Se o documento não existir, retornar mensagem de erro.
### Insert([]char):
insere o documento na primeira página, caso a mesma possua espaço 
disponível. Caso contrário, busca nas páginas seguintes a primeira que possua espaço 
disponível suficiente para armazenar o documento. Em seguida, insere o documento na 
página encontrada. Se todas as páginas em utilização não possuírem espaço disponível, 
retirar uma página do diretório de páginas em branco, adicioná-la no diretório de páginas 
em utilização e inserir o documento nesta página. Caso não seja possível encontrar o espaço 
necessário em nenhuma das 20 páginas, exibir mensagem de erro.
