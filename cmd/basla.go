package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	//	b "yandor2/bin"

	//h "yandor2/helper"
	md "yandor2/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//	Conn = "mongodb://localhost:27017"
	//	Conn = "mongodb+srv://msamidag:msd2095msd@cluster0.ssfgi.mongodb.net/test"
	db    = "teverpan"
	coll  = "kisiler"
	colli = "kullanicilar"
)

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

var Conn = "mongodb+srv://msamidag:msd2095msd@cluster0.ssfgi.mongodb.net/test"
var Aranan string
var Kisiler []md.Calisan
var Kisi md.Calisan
var tpl *template.Template
var e error
var Kullanici User
var Kullanicis User
var ID int
var Kullanicilar []User
var Ch, Hc chan string
var glID int
var kisiKontrol md.Calisan

func init() { //ilk çalışan fonksiyon bu, hazırlık fonksiyonu
	tpl = template.Must(template.ParseGlob("pages/*.html"))

}

func Start() {
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/anasayfa", Anasayfa)   //ok
	http.HandleFunc("/yenigiris", Yenigiris) //ok
	//	http.HandleFunc("/yenikayit", Yenikayit)                   //ok?
	http.HandleFunc("/kisigir", Kisigir)                       //ok
	http.HandleFunc("/kisibil", Kisibil)                       //ok
	http.HandleFunc("/calisanlistesi", CalisanListesi)         //ok
	http.HandleFunc("/kayitbul", KayitBul)                     //ok
	http.HandleFunc("/guncelle", Guncelle)                     //ok
	http.HandleFunc("/arananbul", ArananBul)                   //ok
	http.HandleFunc("/yenigirisgoruntule", YeniGirisGoruntule) //ok
	http.HandleFunc("/kayitguncelle", KayitGuncelle)           //ok
	http.HandleFunc("/kayitsil", KayitSil)                     //ok
	http.HandleFunc("/kullanicikaydet", Kullanicikaydet)       //ok
	http.HandleFunc("/kayit", Kayit)                           //ok
	http.HandleFunc("/kisigun", Kisigun)                       //ok
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/galeri", Galeri) //ok
	http.HandleFunc("/login", login)   //ok
	http.HandleFunc("/parsielArama", ParsielArama)

	//css kodlarının çalışması için css klasörünün go ya tanıtılması
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	//resimlerin gösterilebilmesi için img klasörünün go ya tanıtılması
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.ListenAndServe(":9000", nil)

}

func YeniGirisGoruntule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}
	var kisiKontrol md.Calisan
	Kisi = EkrandanVeriAl(w, r)

	ID := r.FormValue("id")
	fmt.Println(ID)
	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	collection := client.Database("teverpan").Collection("kisiler")
	//ctx := context.Background()

	//mükerrer kayıt kontrolü
	ArananID := Kisi.ID //r.FormValue("advesoyadi")
	filter := bson.D{{"id", ArananID}}
	e = collection.FindOne(context.TODO(), filter).Decode(&kisiKontrol)
	if ArananID == kisiKontrol.ID {
		KayitVar(w, r)
		return
	} else {
		// ekrana girilen bilgilerin kontrolü için webde gösterilmesi
		tpl.ExecuteTemplate(w, "yenigirisgoruntule.html", Kisi)
	}

	//-------------------------------------

	//verilerin jsona kaydı için --> yenigiris()
	//	Yenigiris(w, r)

}

/*
func Kayit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.Ilkkontrol(w, r)
		return
	}
	tpl.ExecuteTemplate(w, "kayit.html", nil)
}
*/
func Kullanicikaydet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}
	kullaniclarVeri, e := ioutil.ReadFile("json/kullanicilar.json")
	Check(e)
	e = json.Unmarshal(kullaniclarVeri, &Kullanicilar)
	glUsername := r.FormValue("username")
	glPassword := r.FormValue("password")
	Kullanici = User{
		Username: glUsername,
		Password: glPassword,
	}

	Kullanicilar = append(Kullanicilar, Kullanici)
	file, e := json.Marshal(Kullanicilar)
	Check(e)
	e = ioutil.WriteFile("json/kullanicilar.json", file, 0666)

	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	Collection := client.Database(db).Collection(colli)

	_, e = Collection.InsertOne(context.TODO(), Kullanici)
	Check(e)
	Welcome(w, r)

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "welcome.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	glUsername := r.FormValue("username")
	glPassword := r.FormValue("password")
	Kullanicis = User{
		Username: glUsername,
		Password: glPassword,
	}

	//kontrol
	clientOptions := options.Client().ApplyURI(Conn)
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	collection := client.Database("teverpan").Collection("kullanicilar")
	//	ctx := context.Background()
	findOptions := options.Find()
	cursor, e := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	Check(e)
	for cursor.Next(context.TODO()) {
		e = cursor.Decode(&Kullanici)
		Check(e)
		if Kullanici == Kullanicis {
			Anasayfa(w, r)
			go func() {
				Ch <- glUsername
				Hc <- glPassword
			}()
			cursor.Close(context.TODO())
			return
		}
	}
	Gecersizkullanici(w, r)
	cursor.Close(context.TODO())
}

func Anasayfa(w http.ResponseWriter, r *http.Request) {
	if len(Kullanicis.Username) != 0 && len(Kullanicis.Password) != 0 {
		if Kullanici == Kullanicis {
			tpl.ExecuteTemplate(w, "anasayfa.html", nil)
		} else {
			Welcome(w, r)
		}
	} else {
		Welcome(w, r)
	}
}

// ekrana bilgileri giriş için sayfa açılması
func Kisigir(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	Kisi = YeniKisisi(w, r)
	ID = Yenikayitno()
	Kisi.ID = ID
	glID := ID  //strconv.Itoa(ID)
	glID = glID //---------------------------------
	tpl.ExecuteTemplate(w, "yenigiris.html", Kisi)
}

func Kisibil(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	//webden kişi verilerinin alınması
	Kisi = YeniKisisi(w, r)

	// ekrana girilen bilgilerin kontrolü için webde gösterilmesi
	tpl.ExecuteTemplate(w, "kisi.html", Kisi)

	//verilerin jsona kaydı için --> yenigiris()
	Yenigiris(w, r)
}

func Kisigun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}
	//webden kişi verilerinin alınması
	ID := Yenikayitno()
	glID := ID
	glID = glID //glID kızmasın diye öylesine
	//fmt.Println(glID, ID)
	Kisi = EkrandanVeriAl(w, r)
	// ekrana girilen bilgilerin kontrolü için webde gösterilmesi
	tpl.ExecuteTemplate(w, "kisiguncel.html", Kisi)

}

//ekrana girilen bilgilerin alınması
func YeniKisisi(w http.ResponseWriter, r *http.Request) md.Calisan {
	ID = Yenikayitno()
	glID = ID //strconv.Itoa(ID)
	//glID = glID //glID kızmasın diye öylesine
	fmt.Println(glID)           //---------------------------------
	Kisi = EkrandanVeriAl(w, r) //ekrandan verileri al
	fmt.Println("ykisinde", Kisi)
	Kisi.ID = glID
	return Kisi
}

//yeni girişin json ve mongodb ye kaydı
func Yenigiris(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}
	//Kisim := YeniKisisi(w, r)

	EkrandanVeriAl(w, r)
	fmt.Println("yyyy", Kisi)
	kisilerVeri, e := ioutil.ReadFile("json/kisiler.json")
	Check(e)
	e = json.Unmarshal(kisilerVeri, &Kisiler)
	Check(e)
	Kisiler = append(Kisiler, Kisi)
	kisilerFile, e := json.Marshal(Kisiler)
	Check(e)
	e = ioutil.WriteFile("json/kisiler.json", kisilerFile, 0666)
	Check(e)

	//mongo
	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	Collection := client.Database(db).Collection(coll)
	//kisi := YeniKisisi(w, r)
	fmt.Println("xxxxxx", Kisi)
	_, e = Collection.InsertOne(context.TODO(), Kisi)
	Check(e)
	Anasayfa(w, r)
}

//kisisil.go
func KayitSil(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, e := template.ParseFiles("./pages/welcome.html")
		t.Execute(w, nil)
		Check(e)
		return
	}
	Aranan := Kisi.NameSurname
	//	Aranan = r.FormValue("aranan")
	fmt.Println("KayitSil () de Aranan", Aranan)

	filter := bson.D{{"nameSurname", Aranan}}
	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	Collection := client.Database(db).Collection(coll)
	//delete tek okiiiii
	//	mongoStart()
	_, e = Collection.DeleteOne(context.TODO(), filter)
	Check(e)
	fmt.Println("Kayıt silindi")
	t, e := template.ParseFiles("./pages/anasayfa.html")
	Check(e)
	t.Execute(w, nil)
}

//kayitbul.go
//var Kisi md.Calisan
var Tpl *template.Template

//var Aranan string
var Tusa []md.Calisan
var Turna []md.Calisan
var Tuska md.Calisan

//var Ch chan string

//--------------------------------
func ParsielArama(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	Tusa = Turna //Tusa yı boşaltmak için

	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	collection := client.Database("teverpan").Collection("kisiler")
	ctx := context.Background()
	Aranano := r.FormValue("advesoyadi")
	filter := bson.D{{"nameSurname", bson.D{{"$regex", Aranano}}}}

	cursor, e := collection.Find(ctx, filter)

	for cursor.Next(ctx) {
		if err := cursor.Decode(&Tuska); err != nil {
			log.Fatal("cursor. Decode ERROR:", err)
		}
		Tusa = append(Tusa, Tuska)
	}
	if len(Tusa) == 0 {
		KayitYok(w, r)
	} else {
		for k := range r.Form {
			delete(r.Form, k)
		}

		t, e := template.ParseFiles("./pages/bulunanlar.html")
		Check(e)
		t.Execute(w, &Tusa)
	}

}

//------------------------------------------------

func KayitBul(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	//	mongodb+srv://msamidag:msd2095msd@cluster0.ssfgi.mongodb.net/test
	clientOptions := options.Client().ApplyURI(Conn)
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	collection := client.Database("teverpan").Collection("kisiler")

	//Aranan = r.FormValue("advesoyad")
	Aranan = r.FormValue("id")
	var aramaman int
	aramaman, _ = strconv.Atoi(Aranan)
	filter := bson.D{{"id", aramaman}} //teverpan kisiler için
	//	filter := bson.D{{"nameSurname", bson.D{{"$regex", Aranan}}}} //teverpan kisiler için
	e = collection.FindOne(context.TODO(), filter).Decode(&Kisi)

	if e != nil {
		KayitYok(w, r)
		return
	}

	t, e := template.ParseFiles("./pages/kisi.html")
	t.Execute(w, &Kisi)

	Ch = make(chan string)
	go func() {
		Ch <- Aranan
	}()
}

//---------------------------------------------------

func Guncelle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}
	t, e := template.ParseFiles("./pages/guncelle.html")
	Check(e)
	t.Execute(w, &Kisi)
}

type Ads struct {
	ResimAd string //string
	Resima  string // int
	Resimb  string // int
	Resimc  string // int
}

func Galeri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	go Kanaldan(Ch)
	go func() {
		Aranan = <-Ch
	}()
	//var a,b,c string
	Kds := Ads{
		ResimAd: strconv.Itoa(Kisi.ID), //Kisi.NameSurname,
		Resima:  "a",
		Resimb:  "b",
		Resimc:  "c",
	}
	t, e := template.ParseFiles("./pages/images.html")
	Check(e)
	t.Execute(w, &Kds)
}

//kayitguncelle.go
func KayitGuncelle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e) //b.CheckError(e)
	Collection := client.Database(db).Collection(coll)
	Aranan := Kisi.NameSurname //r.FormValue("advesoyadi")
	filter := bson.D{{"nameSurname", Aranan}}
	update := bson.D{
		{"$set", bson.D{
			{"id", Kisi.ID},
			{"nameSurname", Kisi.NameSurname},
			{"soygecmis", bson.D{
				{"calisanid", Kisi.ID},
				{"anne", Kisi.Soygecmis.Anne},
				{"baba", Kisi.Soygecmis.Baba},
			}},

			{"ozgecmis", bson.D{
				{"calisanid", Kisi.ID},
				{"krHastalik", Kisi.Ozgecmis.KrHastalik},
				{"tetanoz", Kisi.Ozgecmis.Tetanoz},
			}},

			{"anamnez", bson.D{
				{"calisanid", Kisi.ID},
				{"sigara", Kisi.Anamnez.Sigara},
				{"alkol", Kisi.Anamnez.Alkol},
				{"ameliyat", Kisi.Anamnez.Ameliyat},
			}},
		}}}
	//update
	_, e = Collection.UpdateOne(context.TODO(), filter, update)
	Check(e)
	//	tpl.ExecuteTemplate(w, "anasayfa.html", nil)
	Anasayfa(w, r)
}

//liste.go
var tusa []md.Calisan
var turna []md.Calisan
var tusar md.Calisan
var tpla *template.Template

func ArananBul(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	tusa = turna //tusa yı boşaltmak için boş Calisan yapısına eşitliyoruz
	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	ctx := context.Background()
	collection := client.Database("teverpan").Collection("kisiler")
	Aramalar := r.FormValue("aramalar")
	sAranan := r.FormValue("saranan")
	Aranan := r.FormValue("aranan")
	//	filter := bson.D{{"soygecmis.baba", Aranan}}
	//	filters := bson.D{{"soygecmis." + sAranan, Aranan}}
	filters := bson.D{{Aramalar + "." + sAranan, bson.D{{"$regex", Aranan}}}}

	cursor, e := collection.Find(ctx, filters)

	for cursor.Next(ctx) {
		if err := cursor.Decode(&tusar); err != nil {
			//	if err := cursor.Decode(&Kisi); err != nil {
			log.Fatal("cursor. Decode ERROR:", err)
		}
		tusa = append(tusa, tusar)
	}

	if len(tusa) == 0 {
		KayitYok(w, r)
	} else {
		for k := range r.Form {
			delete(r.Form, k)
		}
		t, e := template.ParseFiles("./pages/bulunanlar.html")
		Check(e)
		t.Execute(w, &tusa)
	}

}

func CalisanListesi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}

	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)
	collection := client.Database("teverpan").Collection("kisiler")

	//kayıtları bul, seçenekleri Find yöntemine ilet
	findOptions := options.Find()

	//listelenecek kayıt sayısı limitini belirleme
	findOptions.SetLimit(20)

	//Kodu çözülen belgeleri saklayabileceğiniz bir dizi tanımlayın
	var results []md.Calisan

	cursor, e := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	Check(e)
	for cursor.Next(context.TODO()) {
		//	e := cursor.Decode(&Kisi)
		e := cursor.Decode(&tusar)
		Check(e)
		results = append(results, tusar)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//işlem sonu cursoru kapat
	cursor.Close(context.TODO())

	if len(results) == 0 {
		KayitYok(w, r)
	} else {
		for k := range r.Form {
			delete(r.Form, k)
		}
		t, e := template.ParseFiles("./pages/bulunanlar.html")
		Check(e)
		t.Execute(w, &results)
	}
}

// ekrandanveral.go
func EkrandanVeriAl(w http.ResponseWriter, r *http.Request) md.Calisan {
	var Kisim md.Calisan
	glID, _ := strconv.Atoi(r.FormValue("id")) //strconv.Itoa(ID)
	glAdsoyad := r.FormValue("adsoyad")
	glAnne := r.FormValue("anne")
	glBaba := r.FormValue("baba")
	glKrHastalik := r.FormValue("krhastalik")
	glTetanoz := r.FormValue("tetanoz")
	glSigara := r.FormValue("sigara")
	glAlkol := r.FormValue("alkol")
	glAmeliyat := r.FormValue("ameliyat")

	Kisim = md.Calisan{
		ID:          glID,
		NameSurname: glAdsoyad,
		Soygecmis: md.Soygecmis{
			Calisanid: glID,
			Anne:      glAnne,
			Baba:      glBaba,
		},
		Ozgecmis: md.Ozgecmis{
			CalisanID:  glID,
			KrHastalik: glKrHastalik,
			Tetanoz:    glTetanoz,
		},
		Anamnez: md.Anamnez{
			CalisanId: glID,
			Sigara:    glSigara,
			Alkol:     glAlkol,
			Ameliyat:  glAmeliyat,
		},
	}
	fmt.Println("ekrandanverial: ", Kisim)
	return Kisim
}

//helper.go
func Gecersizkullanici(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("./pages/gecersizkullanici.html")
	t.Execute(w, nil)
	Check(e)
}

func Kayit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Ilkkontrol(w, r)
		return
	}
	t, e := template.ParseFiles("./pages/kayit.html")
	t.Execute(w, nil)
	Check(e)

}

func Ilkkontrol(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, e := template.ParseFiles("./pages/welcome.html")
		t.Execute(w, nil)
		Check(e)
		return
	}
}

func Kanala(aranan string, ch chan string) {
	ch <- aranan
	fersa := <-ch
	time.Sleep(time.Second * 2)
	fersa = fersa
}

func Kanaldan(ch chan string) {
	Aranan := <-ch
	time.Sleep(time.Second * 2)
	Aranan = Aranan
	return
}

func Yenikayitno() int {
	clientOptions := options.Client().ApplyURI(Conn)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	Check(e)

	collection := client.Database("teverpan").Collection("kisiler")
	findOptions := options.Find()

	/* filtreli kayıt sayısı
	filter := bson.D{{"xxxxx", bson.D{{"$lt", 6}}}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	*/

	// filtresiz toplam kayıt sayısı
	count, err := collection.CountDocuments(context.TODO(), findOptions)
	if err != nil {
		log.Fatal(err)
	}
	a := int(count) + 1

	/* json kayıt sayısı
	kisilerVeri, err := ioutil.ReadFile("json/kisiler.json")
	if err != nil {
		log.Println("json okuma hatası", err)
	}
	err = json.Unmarshal(kisilerVeri, &Kisiler)
	a := len(Kisiler) + 1
	*/
	return a

}

func KayitYok(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		t, e := template.ParseFiles("./pages/kayityok.html")
		Check(e)
		t.Execute(w, nil)
	}

}
func KayitVar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		t, e := template.ParseFiles("./pages/kayitvar.html")
		Check(e)
		t.Execute(w, nil)
	}

}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
