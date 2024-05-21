package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Oda struct {
	Adi    string
	KoordX int
	KoordY int
}

type Tünel struct {
	Oda1 string
	Oda2 string
}

type KarincaÇiftliği struct {
	KarincaSayisi int
	BaşlangiçOda  Oda
	BitişOda      Oda
	Odalar        map[string]Oda
	Tüneller      []Tünel
}

func main() {
	startTime := time.Now()
	if len(os.Args) != 2 {
		fmt.Println("Lütfen bir dosya adi belirtin: go run main.go <dosya_adi>")
		return
	}

	dosyaAdi := os.Args[1]
	dosya, err := os.Open(dosyaAdi)
	if err != nil {
		fmt.Println("Dosya okunamadı:", err)
		return
	}
	defer dosya.Close()

	var karincaÇiftliği KarincaÇiftliği
	karincaÇiftliği.Odalar = make(map[string]Oda)
	var startOdaBelirlenmis, endOdaBelirlenmis bool

	tarayici := bufio.NewScanner(dosya)
	for tarayici.Scan() {
		satir := tarayici.Text()

		if karincaÇiftliği.KarincaSayisi == 0 {
			karincaSayisi, err := strconv.Atoi(satir)
			if err != nil || karincaSayisi <= 0 {
				fmt.Println("ERROR: invalid data format")
				return
			}
			karincaÇiftliği.KarincaSayisi = karincaSayisi
			continue
		}

		if strings.HasPrefix(satir, "##start") {
			if tarayici.Scan() {
				başlangiçOdaBilgisi := strings.Fields(tarayici.Text())
				if len(başlangiçOdaBilgisi) != 3 {
					fmt.Println("ERROR: invalid data format")
					return
				}
				karincaÇiftliği.BaşlangiçOda.Adi = başlangiçOdaBilgisi[0]
				karincaÇiftliği.BaşlangiçOda.KoordX, err = strconv.Atoi(başlangiçOdaBilgisi[1])
				karincaÇiftliği.BaşlangiçOda.KoordY, err = strconv.Atoi(başlangiçOdaBilgisi[2])
				if err != nil || strings.Contains(karincaÇiftliği.BaşlangiçOda.Adi, " ") {
					fmt.Println("ERROR: invalid room name format")
					return
				}
				karincaÇiftliği.Odalar[karincaÇiftliği.BaşlangiçOda.Adi] = karincaÇiftliği.BaşlangiçOda
				startOdaBelirlenmis = true
			}
			continue
		}

		if strings.HasPrefix(satir, "##end") {
			if tarayici.Scan() {
				bitişOdaBilgisi := strings.Fields(tarayici.Text())
				if len(bitişOdaBilgisi) != 3 {
					fmt.Println("ERROR: invalid data format")
					return
				}
				karincaÇiftliği.BitişOda.Adi = bitişOdaBilgisi[0]
				karincaÇiftliği.BitişOda.KoordX, err = strconv.Atoi(bitişOdaBilgisi[1])
				karincaÇiftliği.BitişOda.KoordY, err = strconv.Atoi(bitişOdaBilgisi[2])
				if err != nil || strings.Contains(karincaÇiftliği.BitişOda.Adi, " ") {
					fmt.Println("ERROR: invalid room name format")
					return
				}
				karincaÇiftliği.Odalar[karincaÇiftliği.BitişOda.Adi] = karincaÇiftliği.BitişOda
				endOdaBelirlenmis = true
			}
			continue
		}

		if strings.Contains(satir, " ") && !strings.HasPrefix(satir, "##") {
			odaBilgisi := strings.Fields(satir)
			if len(odaBilgisi) != 3 {
				fmt.Println("ERROR: invalid data format")
				return
			}
			koordinatX, err := strconv.Atoi(odaBilgisi[1])
			koordinatY, err := strconv.Atoi(odaBilgisi[2])
			if err != nil || strings.Contains(odaBilgisi[0], " ") {
				fmt.Println("ERROR: invalid room coordinates")
				return
			}
			if _, mevcut := karincaÇiftliği.Odalar[odaBilgisi[0]]; mevcut { //Aynı ada sahip birden fazla oda eklenmesini engellemek.
				fmt.Println("ERROR: duplicate room names")
				return
			}
			oda := Oda{Adi: odaBilgisi[0], KoordX: koordinatX, KoordY: koordinatY}
			karincaÇiftliği.Odalar[oda.Adi] = oda
			continue
		}

		if strings.Contains(satir, "-") {
			tünelBilgisi := strings.Split(satir, "-")
			if len(tünelBilgisi) != 2 {
				fmt.Println("ERROR: invalid data format")
				return
			}
			if _, odaVar := karincaÇiftliği.Odalar[tünelBilgisi[0]]; !odaVar {
				fmt.Println("ERROR: unknown room in link")
				return
			}
			if _, odaVar := karincaÇiftliği.Odalar[tünelBilgisi[1]]; !odaVar {
				fmt.Println("ERROR: unknown room in link")
				return
			}
			if tünelBilgisi[0] == tünelBilgisi[1] {
				fmt.Println("ERROR: loops in the ant farm")
				return
			}
			tünel := Tünel{Oda1: tünelBilgisi[0], Oda2: tünelBilgisi[1]}
			karincaÇiftliği.Tüneller = append(karincaÇiftliği.Tüneller, tünel)
			continue
		}
	}

	if err := tarayici.Err(); err != nil {
		fmt.Println("Dosya okunurken bir hata oluştu:", err)
		return
	}

	if !startOdaBelirlenmis || !endOdaBelirlenmis || karincaÇiftliği.KarincaSayisi == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Başlangıç odasından hedef odaya giden tüm yolları bul
	başlangiçOda := karincaÇiftliği.BaşlangiçOda.Adi
	hedefOda := karincaÇiftliği.BitişOda.Adi
	tümYollar := TümYollar(karincaÇiftliği.Tüneller, başlangiçOda, hedefOda)
	fmt.Println("Başlangiç odasindan hedef odaya giden tüm yollar:")
	for i, yol := range tümYollar {
		fmt.Printf("Yol %d: %v\n", i+1, yol)
	}
	sort.Slice(tümYollar, func(i, j int) bool {
		return len(tümYollar[i]) < len(tümYollar[j])
	})

	// Yolları filtreleyelim
	filtrelenmişYollar := FiltreleYollar(tümYollar, karincaÇiftliği.KarincaSayisi)
	fmt.Println("Filtrelenmiş Yollar:")
	for i, yol := range filtrelenmişYollar {
		fmt.Printf("Yol %d: %v\n", i+1, yol)
	}
	// En kısa yolu al ve yazdır
	enKisaYol := filtrelenmişYollar[0]
	fmt.Println("En Kısa Yol:")
	fmt.Printf("%v\n", enKisaYol)
	// Karıncaların hareketlerini simüle edelim
	karincaHareketleri := KarincaHareketSimulasyonu(filtrelenmişYollar, karincaÇiftliği.KarincaSayisi, başlangiçOda, hedefOda, enKisaYol)
	// Çıktı formatına uygun olarak çıktının oluşturulması
	fmt.Println(karincaÇiftliği.KarincaSayisi)
	fmt.Println("##start", karincaÇiftliği.BaşlangiçOda.Adi, karincaÇiftliği.BaşlangiçOda.KoordX, karincaÇiftliği.BaşlangiçOda.KoordY)
	fmt.Println("##end", karincaÇiftliği.BitişOda.Adi, karincaÇiftliği.BitişOda.KoordX, karincaÇiftliği.BitişOda.KoordY)

	// Odaların listesi
	for _, oda := range karincaÇiftliği.Odalar {
		fmt.Printf("%s %d %d\n", oda.Adi, oda.KoordX, oda.KoordY)
	}

	// Tünellerin listesi
	for _, tünel := range karincaÇiftliği.Tüneller {
		fmt.Printf("%s-%s\n", tünel.Oda1, tünel.Oda2)
	}

	// Karınca hareketlerinin listesi
	for _, hareket := range karincaHareketleri {
		fmt.Println(hareket)
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Kodun çalışması %.8f dakika sürdü.\n", elapsed.Minutes())
}

// Odaların bağlantılarını kullanarak tüm olası yolları bulan fonksiyon
func TümYollar(tüneller []Tünel, başlangiç, hedef string) [][]string {
	var tümYollar [][]string
	var yolBul func(mevcutOda string, yol []string)
	ziyaretEdilen := make(map[string]bool)

	yolBul = func(mevcutOda string, yol []string) {
		yol = append(yol, mevcutOda)
		if mevcutOda == hedef {
			kopyaYol := make([]string, len(yol))
			copy(kopyaYol, yol)
			tümYollar = append(tümYollar, kopyaYol)
			return
		}
		ziyaretEdilen[mevcutOda] = true
		for _, tünel := range tüneller {
			if tünel.Oda1 == mevcutOda && !ziyaretEdilen[tünel.Oda2] {
				yolBul(tünel.Oda2, yol)
			}
			if tünel.Oda2 == mevcutOda && !ziyaretEdilen[tünel.Oda1] {
				yolBul(tünel.Oda1, yol)
			}
		}
		ziyaretEdilen[mevcutOda] = false
	}

	yolBul(başlangiç, []string{})
	return tümYollar
}

// Yolları filtreler ve çakışan odaları çıkarır
func FiltreleYollar(yollar [][]string, karincaSayisi int) [][]string {
	var filtrelenmisYollar [][]string

	// İki yolun ara odalarda çakışıp çakışmadığını kontrol eden yardımcı fonksiyon
	yollarCakisiyor := func(yol1, yol2 []string) bool {
		kume := make(map[string]bool)
		for _, oda := range yol1[1 : len(yol1)-1] {
			kume[oda] = true
		}
		for _, oda := range yol2[1 : len(yol2)-1] {
			if kume[oda] {
				return true
			}
		}
		return false
	}

	// Çakışmayan yol kombinasyonlarını bulmak için tüm kombinasyonları dene
	var kombinasyonlar func([][]string, int, []int)
	var enIyiKombinasyon []int
	maxYol := 0

	kombinasyonlar = func(yollar [][]string, indeks int, secili []int) {
		if len(secili) > maxYol {
			maxYol = len(secili)
			enIyiKombinasyon = make([]int, len(secili))
			copy(enIyiKombinasyon, secili)
		}

		for i := indeks; i < len(yollar); i++ {
			cakisiyor := false
			for _, s := range secili {
				if yollarCakisiyor(yollar[s], yollar[i]) {
					cakisiyor = true
					break
				}
			}
			if !cakisiyor {
				secili = append(secili, i)
				kombinasyonlar(yollar, i+1, secili)
				secili = secili[:len(secili)-1] //Mevcut yol seçilmiş yollardan çıkarılır ve diğer kombinasyonlar için yeni aramalar yapılır.
			}
		}
	}

	kombinasyonlar(yollar, 0, []int{})

	for _, indeks := range enIyiKombinasyon {
		filtrelenmisYollar = append(filtrelenmisYollar, yollar[indeks])
		if len(filtrelenmisYollar) == karincaSayisi {
			break
		}
	}

	return filtrelenmisYollar
}

func KarincaHareketSimulasyonu(yollar [][]string, karincaSayisi int, başlangiç, hedef string, enKisaYol []string) []string {
	var hareketler []string
	karincaPozisyonu := make(map[int]int)
	karincaHedefte := make(map[int]bool)
	karincaYollari := make(map[int][]string)
	aktifKarincaSayisi := karincaSayisi

	for i := 1; i <= karincaSayisi; i++ {
		if i == karincaSayisi {
			karincaYollari[i] = enKisaYol //Son karınca en kısa yolu takip eder
		} else {
			karincaYollari[i] = yollar[(i-1)%len(yollar)] // diğerleri sırayla takip eder.
		}
		karincaPozisyonu[i] = 0   // bütün karıncalar başlangıç pozisyonundadır.
		karincaHedefte[i] = false // hiçbir karınca bitişe ulaşmamıştır.
	}

	tur := 0
	for aktifKarincaSayisi > 0 {
		tur++
		var turHareketi []string
		tünelKullanimi := make(map[string]bool)

		for i := 1; i <= karincaSayisi; i++ {
			if karincaHedefte[i] {
				continue
			}

			şuAnkiOda := karincaYollari[i][karincaPozisyonu[i]]
			sonrakiOda := karincaYollari[i][karincaPozisyonu[i]+1]
			tünel := fmt.Sprintf("%s-%s", şuAnkiOda, sonrakiOda)
			tersTünel := fmt.Sprintf("%s-%s", sonrakiOda, şuAnkiOda)

			if !tünelKullanimi[tünel] && !tünelKullanimi[tersTünel] {
				turHareketi = append(turHareketi, fmt.Sprintf("L%d-%s", i, sonrakiOda))
				tünelKullanimi[tünel] = true
				tünelKullanimi[tersTünel] = true
				karincaPozisyonu[i]++
				if sonrakiOda == hedef {
					karincaHedefte[i] = true
					aktifKarincaSayisi--
				}
			}
		}

		if len(turHareketi) > 0 {
			hareketler = append(hareketler, strings.Join(turHareketi, " "))
		} else {
			break // Eğer bu turda hareket eden karınca yoksa, döngüyü kır
		}
	}
	return hareketler
}
