# ksp-proboj-2025-jesen

# Čo je to Proboj a ako funguje?

Proboj, skratka pre programátorský boj, je aktivita z KSP sústredení, kde hráči (vy) programujú vlastného bota, ktorý
súťazí v predom pripravenej hre. K hre je taktiež pripravený template bota, ktorý zvláda komunikáciu so serverom a nejaké
užitočné funkcie. Taktiež obsahuje veľmi jednoduchý príklad jednoduchého bota, ktorého môžete ďalej upravovať.

## Štruktúra/harmonogram

Počas proboja bežia hry (matche), ktoré sa skladajú z niekoľko stovák kôl, v ktorých vaši boti hrajú. Počas tejto hry
(matchu) sa nemení mapa, na ktorej hráte a ani ostatní boti, proti ktorým hráte.

Po každej hre (matchy) sa náhodne vygeneruje mapa a boti, ktorí na nej budú hrať a spustí sa hra (match).

**Začiatok**: 3.12.2025 19:30
**Koniec**: 4.12.2025 18:00

## Ciele

Zabaviť sa a vyskúšať si niečo pekné nakódiť.

A pre tých kompetetívnejších z vás: Počas hry (matchu) bude váš bot získavať body za rôzne úkony (vid. Pravidlá) počas hry.
Tieto body sa sčítavajú medzi hrami (matchmi). Kto bude mať na konci najviac bodov, vyhráva.

# Pravidlá hry

## Krátky opis hry

Každý hráč riadi vesmírnu flotilu lodí v 2D priestore. Cieľom je získať body pomocou ovládania asteroidov, ťažby zdrojov, stavby nových lodí a boja proti protivníkom. Hra prebieha v reálnom čase, kde pohyb lodí je fyzikálne realistický - ak je loď v pohybe, zostáva v pohybe.

### Herné prostredie

- **Mapa**: Kruhový priestor s polomerom 15 000 jednotiek
- **Asteroidy**: 500 náhodne generovaných asteroidov (palivové a kamenné)
- **Červie diery**: 25 párov teleportačných bodov pre strategický presun
- **Počet kôl**: Maximálne 1 000 kôl na hru

### Základné zdroje

- **Palivo**: Potrebné pre pohyb lodí a stavbu nových plavidiel
- **Kameň**: Potrebný pre stavbu lodí a opravy poškodených plavidiel
- **Počiatočné zdroje**: Každý hráč začína s 1 000 paliva a 1 000 kameňa

## Lodné typy

### MotherShip (Materská loď)
- **Funkcia**: Základňa hráča, ukladanie zdrojov, opravy iných lodí
- **Schopnosti**: Kúpa ostatné lode, spravovať poškodené lode v blízkosti
- **Zdravie**: 100 HP
- **Špeciálne**: Nezničiteľná, chráni ostatné lode v okolí 50 jednotiek

### SuckerShip (Cucač)
- **Funkcia**: Ťažba palivových asteroidov
- **Efektivita**: Špeciálne navrhnutá na získavanie paliva
- **Dosah ťažby**: 50 jednotiek
- **Množstvo ťažby**: 10 paliva za kolo

### DrillShip (Vŕtačka)
- **Funkcia**: Ťažba kamenných asteroidov
- **Efektivita**: Špeciálne navrhnutá na získavanie kameňa
- **Dosah ťažby**: 50 jednotiek
- **Množstvo ťažby**: 10 kameňa za kolo

### TankerShip (Cisterna)
- **Funkcia**: Transport paliva medzi loďami
- **Špeciálne**: Zvýšená efektivita pohybu (3x menšia spotreba paliva)
- **Dosah presunu**: 20 jednotiek

### TruckShip (Tatrovka)
- **Funkcia**: Transport kameňa medzi loďami
- **Špeciálne**: Zvýšená efektivita pohybu (3x menšia spotreba paliva)
- **Dosah presunu**: 20 jednotiek

### BattleShip (Bojová loď)
- **Funkcia**: Útok na nepriateľské lode
- **Zbrane**: Dosah streľby 500 jednotiek, damage 25 HP
- **Obmedzenia**: Nemôže útočiť na MotherShip, nemôže útočiť na lode v ochrannom polomere MotherShip

## Herné mechaniky

### Pohyb lodí
- **Fyzika**: Lodě si udržujú svoj pohybový vektor
- **Základný pohyb**: 1.0 jednotky zdarma bez spotreby paliva
- **Spotreba paliva**: Za pohyb nad 1.0 jednotku sa platí palivo podľa typu lode
- **Špeciálne lode**: TankerShip a TruckShip majú 3x nižšiu spotrebu paliva

### Ťažba asteroidov
- **FuelAsteroid**: SuckerShip môže ťažiť v dosahu 50 jednotiek
- **RockAsteroid**: DrillShip môže ťažiť v dosahu 50 jednotiek
- **Pohyb asteroidov**: Asteroidy sa pomaly pohybujú, čo vytvára dynamické prostredie
- **Vyčerpanie**: Asteroidy po úplnom vyčerpaní zmiznú z mapy

### Červie diery (Wormholes)
- **Teleportácia**: Ak loď vstúpi do rádiusu 5 jednotiek okolo červiej diery, okamžite sa teleportuje
- **Páry**: 25 párov červích dier, každý má svoj protipól
- **Bezpečná vzdialenosť**: Loď sa objaví v minimálnej vzdialenosti 10 jednotiek od cieľovej diery
- **Stratégia**: Umožňujú rýchly presun medzi vzdialenými časťami mapy

### Bojový systém
- **Povolené lode**: Iba BattleShip môže útočiť
- **Dosah útoku**: 500 jednotiek
- **Poškodenie**: 25 HP za zásah
- **Zničenie**: Loď po dosiahnutí 0 HP sa zničí a zanechá za sebou asteroidy s palivom a kameňom
- **Ochranný polomer**: Lode v dosahu 50 jednotiek od svojej MotherShip sú chránené pred útokmi

## Hracie príkazy

V každom kole môže hráč vykonať niekoľko z týchto príkazov:

### Buy (Nákup lode)
- **Cena**: 250 kameňa + 100 paliva
- **Dostupné typy**: SuckerShip, DrillShip, TankerShip, TruckShip, BattleShip
- **Zjavenie**: Nová loď sa objaví na pozícii MotherShip

### Move (Pohyb)
- **Mechanizmus**: Zmena pohybového vektora lode
- **Spotreba**: Palivo sa stráca podľa vzdialenosti a typu lode
- **Limitácia**: Maximálna zmena vektora je obmedzená pre vyváženosť hry

### Load (Presun kameňa)
- **Podmienky**: Vzdialenosť medzi loďami maximálne 20 jednotiek
- **Mechanizmus**: Presun kameňa zo zdrojovej lode do cieľovej lode
- **Ohraničenie**: Iba funkčné lode môžu presúvať zdroje

### Siphon (Presun paliva)
- **Podmienky**: Vzdialenosť medzi loďami maximálne 20 jednotiek
- **Mechanizmus**: Presun paliva zo zdrojovej lode do cieľovej lode
- **Obmedzenia**: MotherShip a TankerShip musia participovať v transferi paliva aspoň na jednej strane

### Shoot (Útok)
- **Podmienky**: Vzdialenosť medzi loďami maximálne 500 jednotiek
- **Povolené**: Iba BattleShip môže útočiť
- **Obmedzenia**: Nemôže útočiť na MotherShip ani na lode v ochrannom polomere

### Repair (Oprava)
- **Podmienky**: Loď v dosahu 50 jednotiek od MotherShip
- **Cena**: 15 kameňa za operáciu
- **Efekt**: Obnoví 30 HP (maximálne do 100 HP)

## Ovládanie asteroidov a bodovanie

### Získavanie kontroly
- **Mechanizmus**: Pasívne zaberanie asteroidov prítomnosťou lode v dosahu ťažby
- **Rýchlosť**: Asteroid sa pomaly zaberá, ak je v okolí loď daného hráča
- **Súťaženie**: Viacerí hráči môžu súťažiť o rovnaký asteroid, ale kontrolovať ho môže vždy len jeden

### Bodovací systém za asteroidy
- **Povrchová plocha**: Body sa prideľujú na základe dobytej povrchovej plochy asteroidu
- **Veľkosť asteroidu**: Čím väčší asteroid, tým viac bodov je možné získať za jeho ovládnutie
- **Súťaženie o kontrolu**: Ak viacerí hráči súťažia o rovnaký asteroid, body sa najprv odpočítajú od aktuálneho ovládajúceho hráča a potom sa pripočítajú novému ovládajúcemu hráčovi
- **Zničenie asteroidu**: Ak je asteroid úplne vyťažený, zničí sa a všetky body za tento asteroid sa stratia
- **Strategická rovnováha**: Vzniká tak strategické napätie medzi ťažbou zdrojov a udržaním kontroly nad územím

### Skórovací systém
- **Body**: Získavajú sa za ovládané asteroidy (na základe povrchovej plochy), ťažbu zdrojov, zničenie nepriateľských lodí a ďalšie herné akcie
- **Kontrola asteroidov**: Primárny zdroj bodov založený na dobytej povrchovej ploche
- **Dlhodobá strategia**: Udržanie kontroly nad veľkými asteroidmi prináša stabilný príjem bodov

## Tip pre vývojárov botov

### Základné stratégie
- **Rozvoj**: Začať s ťažbou a stavbou základných lodí
- **Logistika**: Efektívne presúvať zdroje medzi loďami
- **Boj**: Chrániť svoje lode a útočiť na oslabené ciele
- **Presun**: Využívať červie diery pre strategické manévre
