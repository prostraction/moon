package moon

import (
	"math"
	"time"
)

var months = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

var phasesEn = []string{"Waxing Crescent", "First quarter", "Waxing Gibbous", "Full Moon", "Waning Gibbous", "Third quarter", "Waning Crescent", "New Moon"}
var phasesRu = []string{"Растущий серп", "Первая четверть", "Растущая Луна", "Полнолуние", "Убывающая луна", "Последняя четверть", "Убывающий серп", "Новолуние"}
var phasesEs = []string{"Luna creciente", "Cuarto creciente", "Gibosa creciente", "Luna llena", "Gibosa menguante", "Cuarto menguante", "Luna menguante", "Luna nueva"}
var phasesDe = []string{"Zunehmende Sichel", "Erstes Viertel", "Zunehmender Mond", "Vollmond", "Abnehmender Mond", "Letztes Viertel", "Abnehmende Sichel", "Neumond"}
var phasesFr = []string{"Premier croissant", "Premier quartier", "Gibbeuse croissante", "Pleine lune", "Gibbeuse décroissante", "Dernier quartier", "Dernier croissant", "Nouvelle lune"}
var phasesJp = []string{"三日月", "上弦の月", "十三夜月", "満月", "十六夜月", "下弦の月", "有明の月", "新月"}

var phasesEmoji = []string{"🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘", "🌑"}

var signsEn = []string{"Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces", "Aries", "Taurus", "Gemini", "Cancer", "Leo"}
var signsRu = []string{"Дева", "Весы", "Скорпион", "Стрелец", "Козерог", "Водолей", "Рыбы", "Овен", "Телец", "Близнецы", "Рак", "Лев"}
var signsEs = []string{"Virgo", "Libra", "Escorpio", "Sagitario", "Capricornio", "Acuario", "Piscis", "Aries", "Tauro", "Géminis", "Cáncer", "Leo"}
var signsDe = []string{"Jungfrau", "Waage", "Skorpion", "Schütze", "Steinbock", "Wassermann", "Fische", "Widder", "Stier", "Zwillinge", "Krebs", "Löwe"}
var signsFr = []string{"Vierge", "Balance", "Scorpion", "Sagittaire", "Capricorne", "Verseau", "Poissons", "Bélier", "Taureau", "Gémeaux", "Cancer", "Lion"}
var signsJp = []string{"おとめ座", "てんびん座", "さそり座", "いて座", "やぎ座", "みずがめ座", "うお座", "おひつじ座", "おうし座", "ふたご座", "かに座", "しし座"}

var signsEmoji = []string{"♍", "♎", "♏", "♐", "♑", "♒", "♓", "♈", "♉", "♊", "♋", "♌"}

const Fhour = 24.
const Fminute = 24. * 60.
const Fseconds = 24. * 60. * 60.

const toRad = math.Pi / 180.
