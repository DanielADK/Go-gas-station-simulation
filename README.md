# GOLANG Simulace čerpací stanice [česky]

Tento projekt je implementací simulace čerpací stanice. Cílem je simulovat proces tankování vozidel různých typů paliva 
(gas, diesel, lpg, electric, hydrogen) a následného placení na pokladnách. Projekt je vytvořen jako součást předmětu CTC
na Fakultě mechatroniky, informatiky a mezioborových studií na Technické univerzitě v Liberci (FM TUL).

## Zadání
Úkolem je implementovat systém, kde hlavní program generuje auta s různými typy paliva. Každé auto je následně přiřazeno
k čerpací stanici na základě typu paliva a délky fronty. Po tankování je vozidlo přesunuto do fronty u pokladny, kde se 
simuluje proces placení. Simulace zahrnuje správu front a paralelní zpracování událostí.

### Start
1. Načtení konfigurace
2. Start goroutin se stanicemi a pokladnami
### Tankování vozidla a placení
3. Generování vozidel a přiřazení k čerpací stanici
5. Přesun vozidla k pokladně (blokování stanice)
6. Zaplacení = uvolnění pokladny i stanice
### Konec
7. Jakmile je generování dokončeno a všechny fronty stanic jsou prázdné, uzavírají se všechna vlánka
8. Počítání a výpis statistik

## Technologie
Projekt je napsán v jazyce Go, využívající goroutiny pro paralelní zpracování a kanály pro komunikaci mezi částmi 
systému. Je zaměřen na efektivní správu souběžnosti a poskytuje realistickou simulaci provozu na čerpací stanici.

## Použití
Pro spuštění simulace je potřeba mít nainstalované prostředí Go. Po klonování repozitáře můžete simulaci spustit 
příkazem:

```bash
go run main.go
```
Nebo pomocí Dockeru:

```bash
# Sestavte obraz v Dockeru
docker build -t gas-station-simulation .

# Spusťte kontejner z obrazu
docker run -d --name gas-station-simulation-container gas-station-simulation
```

Projekt dále obsahuje konfigurační soubory, které umožňují nastavit parametry simulace, jako jsou typy a počet stanic a 
pokladen.

# GOLANG Gas Station Simulation [english]

This project is an implementation of a gas station simulation. The goal is to simulate the fueling process for vehicles 
with various types of fuel (gas, diesel, lpg, electric, hydrogen) and the subsequent payment process at the cash 
registers. The project is created as part of the CTC course at the Faculty of Mechatronics, Informatics, and 
Interdisciplinary Studies at the Technical University of Liberec (FM TUL).

## Assignment
The task is to implement a system where the main program generates cars with different fuel types. Each car is then 
assigned to a fueling station based on the fuel type and the length of the queue. After fueling, the vehicle is moved 
to the queue at the cash register, where the payment process is simulated. The simulation includes queue management and 
parallel processing of events.

### Start
1. Configuration loading 
2. Starting goroutines with stations and cash registers
### Vehicle refueling and payment
3. Generating vehicles and assigning them to a gas station
4. Moving the vehicle to the cash register (blocking the station)
5. Payment = release of the cash register and station
### End
6. Once the generation is complete and all station queues are empty, all threads are closed
7. Calculation and output of statistics

## Technology
The project is written in Go, utilizing goroutines for parallel processing and channels for communication between parts 
of the system. It focuses on efficient concurrency management and provides a realistic simulation of operations at a 
gas station.

## Usage
To run the simulation, you need to have the Go environment installed. After cloning the repository, you can start the 
simulation with the command:

```bash
go run main.go
```
Or with Docker:
```bash
# Build Docker image
docker build -t gas-station-simulation .

# Run container from image
docker run -d --name gas-station-simulation-container gas-station-simulation
```
The project also includes configuration files that allow setting the parameters of the simulation, such as the types 
and the number of stations and cash registers.