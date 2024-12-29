package main

import (
	"fmt"
	"iter"
	"math"
	"math/rand"
	"slices"
	"sync"
)

type WürfelErgebnis int

const (
	Pflaume WürfelErgebnis = iota
	Birne
	Kirsche
	Apfel
	Obstkorb
	Rabe
)

func WerfeWürfel() WürfelErgebnis {
	return WürfelErgebnis(rand.Intn(6))
}

type Früchte int
type Rabenteile int

type Obstgarten struct {
	Pflaumen, Birnen, Kirschen, Äpfel Früchte
	Rabenteile
}

func (o Obstgarten) Ergebnis() (gewonnen, verloren bool) {
	gewonnen = o.Pflaumen == 0 && o.Birnen == 0 && o.Kirschen == 0 && o.Äpfel == 0
	verloren = o.Rabenteile == 9
	return
}

func SpieleObstgarten(spieler []Spieler) (gewonnen bool) {
	obstgarten := Obstgarten{
		Pflaumen: 10,
		Birnen:   10,
		Kirschen: 10,
		Äpfel:    10,
	}
	for _, aktuellerSpieler := range RotiereUnendlich(spieler) {
		switch WerfeWürfel() {
		case Pflaume:
			obstgarten.Pflaumen.Nehme()
		case Birne:
			obstgarten.Birnen.Nehme()
		case Kirsche:
			obstgarten.Kirschen.Nehme()
		case Apfel:
			obstgarten.Äpfel.Nehme()
		case Obstkorb:
			aktuellerSpieler.MachObstkorbZug(&obstgarten)
		case Rabe:
			obstgarten.Rabenteile.EinsDazu()
		}

		if gewonnen, verloren := obstgarten.Ergebnis(); gewonnen {
			return true
		} else if verloren {
			return false
		}
	}
	panic("should never be reached")
}

func main() {
	const anzahlSpiele = 500000
	fmt.Printf("Anzahl Spiele: %d, Sigma = %.2f%%\n", anzahlSpiele, 100.0/math.Sqrt(float64(anzahlSpiele)))
	var wg sync.WaitGroup
	for _, spieler := range [][]Spieler{
		{Optimierer{GuteStrategie}},
		{Optimierer{SchlechteStrategie}},
		{Optimierer{GuteStrategie}, Optimierer{GuteStrategie}},
		{Optimierer{GuteStrategie}, Optimierer{SchlechteStrategie}},
		{Optimierer{GuteStrategie}, Kirschenliebhaber{}},
		{Kirschenliebhaber{}, Kirschenliebhaber{}},
		{Kirschenliebhaber{}},
		{Kirschenliebhaber{}, Optimierer{SchlechteStrategie}},
		{Kirschenliebhaber{}, Kirschenliebhaber{}, Optimierer{GuteStrategie}},
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			anzahlGewonnen := 0
			for spiel := 0; spiel < anzahlSpiele; spiel++ {
				if gewonnen := SpieleObstgarten(spieler); gewonnen {
					anzahlGewonnen++
				}
			}
			fmt.Printf("Gewinnwskt %.2f%% mit %s\n", float64(100*anzahlGewonnen)/float64(anzahlSpiele), spieler)
		}()
	}
	wg.Wait()
}

func (frucht *Früchte) Nehme() bool {
	if *frucht > 0 {
		*frucht--
		return true
	} else {
		return false
	}
}

func (r *Rabenteile) EinsDazu() {
	*r++
}

type Spieler interface {
	MachObstkorbZug(obstgarten *Obstgarten)
}

type Kirschenliebhaber struct {
}

func (s Kirschenliebhaber) MachObstkorbZug(obstgarten *Obstgarten) {
	allesAußerKirschen := Verdopple([]*Früchte{&obstgarten.Pflaumen, &obstgarten.Äpfel, &obstgarten.Birnen})
	Durcheinander(allesAußerKirschen)
	MachObstkorbZug(append(Verdopple([]*Früchte{&obstgarten.Kirschen}), allesAußerKirschen...), (*Früchte).Nehme)
}

var (
	GuteStrategie      = Strategie{"VieleZuerst", (*FrüchteMitAnzahl).VieleZuerst}
	SchlechteStrategie = Strategie{"WenigeZuerst", (*FrüchteMitAnzahl).WenigeZuerst}
)

type Strategie struct {
	Name              string
	FrüchteSortierung func(a, b *FrüchteMitAnzahl) int
}

type Optimierer struct {
	Strategie Strategie
}

type FrüchteMitAnzahl struct {
	*Früchte
	Anzahl int
}

func (a *FrüchteMitAnzahl) VieleZuerst(b *FrüchteMitAnzahl) int {
	return b.Anzahl - a.Anzahl
}

func (a *FrüchteMitAnzahl) WenigeZuerst(b *FrüchteMitAnzahl) int {
	return a.Anzahl - b.Anzahl
}

func (s Optimierer) MachObstkorbZug(obstgarten *Obstgarten) {
	var früchteMitAnzahl []*FrüchteMitAnzahl
	fügeHinzu := func(frucht *Früchte, delta int) {
		früchteMitAnzahl = append(früchteMitAnzahl, &FrüchteMitAnzahl{frucht, int(*frucht) + delta})
	}
	fügeHinzu(&obstgarten.Pflaumen, 0)
	fügeHinzu(&obstgarten.Birnen, 0)
	fügeHinzu(&obstgarten.Kirschen, 0)
	fügeHinzu(&obstgarten.Äpfel, 0)
	fügeHinzu(&obstgarten.Pflaumen, -1)
	fügeHinzu(&obstgarten.Birnen, -1)
	fügeHinzu(&obstgarten.Kirschen, -1)
	fügeHinzu(&obstgarten.Äpfel, -1)
	Durcheinander(früchteMitAnzahl)
	slices.SortStableFunc(früchteMitAnzahl, s.Strategie.FrüchteSortierung)
	MachObstkorbZug(früchteMitAnzahl, (*FrüchteMitAnzahl).Nehme)
}

func MachObstkorbZug[T any](früchte []T, nimm func(T) bool) {
	genommen := 0
	for _, frucht := range früchte {
		if genommen == 2 {
			break
		} else if hatGeklappt := nimm(frucht); hatGeklappt {
			genommen++
		}
	}
}

func RotiereUnendlich[T any](s []T) iter.Seq2[int, T] {
	i := 0
	return func(yield func(int, T) bool) {
		for {
			if yield(i, s[i%len(s)]) {
				i++
			} else {
				break
			}
		}
	}
}

func Verdopple[T any](s []T) []T {
	return append(s, s...)
}

func Durcheinander[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func (w WürfelErgebnis) String() string {
	switch w {
	case Rabe:
		return "Rabe"
	case Pflaume:
		return "Pflaume"
	case Birne:
		return "Birne"
	case Kirsche:
		return "Kirsche"
	case Apfel:
		return "Apfel"
	case Obstkorb:
		return "Obstkorb"
	default:
		panic("unbekanntes Würfelergebnis")
	}
}

func (s Optimierer) String() string {
	return fmt.Sprintf("Optimierer(%s)", s.Strategie.Name)
}

func (s Kirschenliebhaber) String() string {
	return "Kirschenliebhaber"
}
