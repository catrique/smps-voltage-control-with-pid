# smps-voltage-control-with-pid

## Simulação de uma Fonte Chaveada (SMPS) com Controle PID para Regulação de Tensão

Este projeto apresenta uma simulação detalhada de uma **Fonte de Alimentação Chaveada (SMPS)** do tipo **Conversor Buck**, incorporando um **Controlador PID** avançado para a regulação de tensão em tempo real. A simulação é desenvolvida em **Go** para a lógica de controle e física, e utiliza **Lua com o framework LÖVE** para uma interface gráfica interativa que visualiza o comportamento do sistema.

## Visão Geral
Fontes chaveadas são componentes cruciais em eletrônica de potência, oferecendo alta eficiência na conversão de energia. No entanto, a manutenção de uma tensão de saída estável, especialmente sob variações de carga ou ruído na entrada, exige sistemas de controle robustos. Este simulador explora a aplicação de um controlador PID para atingir essa estabilidade, demonstrando conceitos como anti-windup, filtragem de ruído derivativo e auto-sintonia, essenciais para o desempenho em cenários reais.

## Funcionalidades
- **Simulação de Conversor Buck:** Modelo matemático preciso de um conversor Buck.
- **Controle PID Avançado:** Implementação de um controlador PID com anti-windup e filtro de ruído derivativo.
- **Auto-Sintonia PID:** Algoritmo de busca binária para otimização dos parâmetros PID (Kp, Ki, Kd) com base na resposta transitória da planta.
- **Perturbações Realistas:** Inclusão de ripple de 120Hz e ruído aleatório na tensão de entrada para testar a robustez do controle.
- **Interface Gráfica em Tempo Real:** Visualização dinâmica da tensão de entrada, tensão de saída e erro através de gráficos interativos.
- **Arquitetura Concorrente:** Utilização de *goroutines* e *channels* em Go para desacoplar a geração de ruído e a comunicação com a GUI, otimizando a performance da simulação.

## Conceitos Técnicos

### Conversor Buck (SMPS)
O conversor Buck é um tipo de SMPS que converte uma tensão de entrada DC mais alta em uma tensão de saída DC mais baixa. A simulação modela o comportamento de um circuito Buck ideal, onde a tensão de saída é controlada pelo ciclo de trabalho (duty cycle) do chaveamento, que por sua vez é determinado pelo controlador PID.

### Controle PID (Proporcional-Integral-Derivativo)
O controlador PID é um algoritmo de controle de feedback amplamente utilizado na indústria. Ele calcula um 
sinal de controle baseado em três termos:
- **Proporcional (P):** Responde ao erro atual.
- **Integral (I):** Responde ao acúmulo de erros passados, eliminando o erro em regime permanente.
- **Derivativo (D):** Responde à taxa de variação do erro, prevendo o comportamento futuro do sistema.

Neste projeto, o PID foi aprimorado com:
- **Anti-Windup:** Evita a saturação do termo integral quando o atuador atinge seus limites, melhorando a resposta do sistema após a saturação.
- **Filtro de Ruído Derivativo:** Suaviza a ação derivativa, que é sensível a ruídos de alta frequência, prevenindo oscilações indesejadas.
- **Derivada sobre a Saída:** A ação derivativa é aplicada diretamente à medição da saída, evitando o *derivative kick* quando o ponto de ajuste (setpoint) é alterado.

### Auto-Sintonia
O algoritmo de auto-sintonia busca um parâmetro `alpha` ideal que define a largura de banda do sistema. Ele realiza uma série de simulações e ajusta `alpha` para encontrar o ponto onde o sistema responde rapidamente ao setpoint sem apresentar um *overshoot* excessivo (neste caso, limitado a 5%).

## Como Usar

### Pré-requisitos
Para rodar a simulação, você precisará ter instalado:
- **Go (versão 1.26.1 ou superior):** [Download e Instalação](https://golang.org/doc/install)
- **LÖVE 2D (versão 11.3 ou superior):** [Download e Instalação](https://love2d.org/)

### Instalação
1. Clone o repositório:
   ```bash
   git clone https://github.com/catrique/smps-voltage-control-with-pid.git
   cd smps-voltage-control-with-pid
   ```


### Execução
Execute o simulador a partir da raiz do projeto:
```bash
go run .
```

O programa solicitará os seguintes parâmetros via terminal:
- **Tensão de entrada (Vin):** Tensão de entrada nominal da fonte.
- **Tensão de saída desejada (Vtarget):** Tensão que o controlador PID tentará manter na saída.
- **Indutância (L) em Henrys:** Valor do indutor do conversor Buck.
- **Capacitância (C) em Farads:** Valor do capacitor do conversor Buck.
- **Resistência (R) em Ohms:** Valor da resistência de carga.

Após a entrada dos parâmetros, a simulação será iniciada e a interface gráfica do LÖVE 2D será aberta, exibindo os gráficos em tempo real.

## Estrutura do Projeto
- `main.go`: Ponto de entrada da aplicação Go, orquestra a simulação e a comunicação com a GUI.
- `engine/`: Contém a lógica central da simulação:
    - `plant.go`: Modelo matemático do conversor Buck.
    - `pid.go`: Implementação do controlador PID e algoritmo de auto-sintonia.
    - `power_supply.go`: Simulação da fonte de alimentação com ripple.
    - `noise.go`: Gerador de ruído aleatório.
    - `utils.go`: Funções utilitárias, incluindo a análise de Routh-Hurwitz.
- `gui/`: Contém a interface gráfica em Lua:
    - `main.lua`: Código da GUI que recebe dados via stdin e os plota.

## Contribuição
Contribuições são bem-vindas! Sinta-se à vontade para abrir *issues* ou *pull requests* para melhorias, correções de bugs ou novas funcionalidades.

## Licença
Este projeto está licenciado sob a Licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
