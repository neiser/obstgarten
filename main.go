package main

import (
	"iter"
	"log"
	"math/rand"
	"slices"
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

type Frucht int
type Rabenteile int

type Obstgarten struct {
	Pflaumen, Birnen, Kirschen, Äpfel Frucht
	Rabenteile
}

func BeginneObstgarten() Obstgarten {
	return Obstgarten{
		Pflaumen: 10,
		Birnen:   10,
		Kirschen: 10,
		Äpfel:    10,
	}
}

func (o Obstgarten) Ergebnis() (gewonnen, verloren bool) {
	gewonnen = o.Pflaumen == 0 && o.Birnen == 0 && o.Kirschen == 0 && o.Äpfel == 0
	verloren = o.Rabenteile == 9
	return
}

func (frucht *Frucht) Nehme() bool {
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
	MachZug(obstgarten *Obstgarten)
}

type Optimierer struct {
}

func (s Optimierer) MachZug(obstgarten *Obstgarten) {
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
		type Paar struct {
			*Frucht
			Anzahl int
		}
		paare := []Paar{
			{&obstgarten.Pflaumen, int(obstgarten.Pflaumen)},
			{&obstgarten.Birnen, int(obstgarten.Birnen)},
			{&obstgarten.Kirschen, int(obstgarten.Kirschen)},
			{&obstgarten.Äpfel, int(obstgarten.Äpfel)},
			{&obstgarten.Pflaumen, int(obstgarten.Pflaumen) - 1},
			{&obstgarten.Birnen, int(obstgarten.Birnen) - 1},
			{&obstgarten.Kirschen, int(obstgarten.Kirschen) - 1},
			{&obstgarten.Äpfel, int(obstgarten.Äpfel) - 1},
		}
		Durcheinander(paare)
		slices.SortStableFunc(paare, func(a, b Paar) int {
			return a.Anzahl - b.Anzahl
		})
		MachObstkorbZug(paare, func(paar Paar) bool {
			return paar.Frucht.Nehme()
		})
	case Rabe:
		obstgarten.Rabenteile.EinsDazu()
	}
}

type Kirschenliebhaber struct {
}

func (s Kirschenliebhaber) MachZug(obstgarten *Obstgarten) {
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
		allesAußerKirschen := Verdopple(Früchte{&obstgarten.Pflaumen, &obstgarten.Äpfel, &obstgarten.Birnen})
		Durcheinander(allesAußerKirschen)
		MachObstkorbZug(append(Verdopple(Früchte{&obstgarten.Kirschen}), allesAußerKirschen...), (*Frucht).Nehme)
	case Rabe:
		obstgarten.Rabenteile.EinsDazu()
	}
}

type Früchte []*Frucht

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

func main() {
	anzahlSpiele := 100000
	anzahlGewonnen := 0
	for spiel := 0; spiel < anzahlSpiele; spiel++ {
		obstgarten := BeginneObstgarten()
		spieler := []Spieler{Optimierer{}, Optimierer{}}
		for _, aktuellerSpieler := range RotiereUnendlich(spieler) {
			aktuellerSpieler.MachZug(&obstgarten)
			if gewonnen, verloren := obstgarten.Ergebnis(); gewonnen {
				anzahlGewonnen++
				break
			} else if verloren {
				break
			}
		}
	}

	log.Printf("Gewinnwahrscheinlichkeit: %.2f", float64(100*anzahlGewonnen)/float64(anzahlSpiele))

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
	return "Optimierer"
}

func (s Kirschenliebhaber) String() string {
	return "Kirschenliebhaber"
}
