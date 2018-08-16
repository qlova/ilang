package software 

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/global"
import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable{
	compiler.English: "software",
	compiler.Dutch: "software",
	compiler.Maori: "taupānga",
	compiler.Chinese: "软件",
	compiler.French: "logiciel",
	compiler.Spanish: "software",
	compiler.Japanese: "ソフトウェア",
	compiler.Samoan: "polokalama",
	compiler.German: "software",
	
	compiler.Afrikaans: "sagteware",
	compiler.Albanian: "program",
	compiler.Amharic: "ሶፍትዌር",
	compiler.Arabic: "البرمجيات",
	compiler.Armenian: "ծրագրայինապահովում",
	compiler.Azerbaijani: "proqramtəminatı",
	
	compiler.Basque: "software",
	compiler.Belarusian: "праграмнаезабеспячэнне",
	compiler.Bengali: "সফটওয়্যার",
	compiler.Bosnian: "softver",
	compiler.Bulgarian: "софтуер",
	compiler.Burmese: "ဆော့ဖျဝဲ",
	
	compiler.Catalan: "programari",
	compiler.Cebuano: "software",
	compiler.Chichewa: "software",
	compiler.Corsican: "software",
	compiler.Croatian: "softver",
	compiler.Czech: "software",
	
	compiler.Danish: "software",
	
	compiler.Esperanto: "programaro",
	compiler.Estonian: "tarkvara",
	
	compiler.Filipino: "software",
	compiler.Finnish: "ohjelmisto",
	compiler.Frisian: "software",
	
	compiler.Galician: "software",
	compiler.Georgian: "პროგრამულიუზრუნველყოფა",
	compiler.Greek: "λογισμικό",
	compiler.Gujarati: "સોફ્ટવેર",
	
	compiler.HaitianCreole: "lojisyèl",
	compiler.Hausa: "software",
	compiler.Hawaiian: "polokalamu",
	compiler.Hebrew: "תוֹכנָה",
	compiler.Hindi: "सॉफ्टवेयर",
	compiler.Hmong: "software",
	compiler.Hungarian: "szoftver",
	
	compiler.Icelandic: "hugbúnaður",
	compiler.Igbo: "ngwanrọ",
	compiler.Indonesian: "perangkatlunak",
	compiler.Irish: "bogearraí",
	compiler.Italian: "software",
	
	compiler.Javanese: "pirantilunak",
	
	compiler.Kannada: "ಸಾಫ್ಟ್ವೇರ್",
	compiler.Kazakh: "бағдарламалыққамтамасызету",
	compiler.Khmer: "កម្មវិធី",
	compiler.Klingon: "ngoq",
	compiler.Korean: "소프트웨어",
	compiler.Kurdish: "nivîsbar",
	compiler.Kyrgyz: "программалыккамсыздоо",
	
	compiler.Lao: "ຊອບແວ",
	compiler.Latin: "software",
	compiler.Latvian: "programmatūra",
	compiler.Lithuanian: "programinėįranga",
	compiler.Luxembourgish: "software",
	
	compiler.Macedonian: "софтвер",
	compiler.Malagasy: "rindrambaiko",
	compiler.Malay: "perisian",
	compiler.Malayalam: "സോഫ്റ്റ്വെയർ",
	compiler.Maltese: "softwer",
	compiler.Marathi: "सॉफ्टवेअर",
	compiler.Mongolian: "програмхангамж",
	
	compiler.Nepali: "सफ्टवेयर",
	compiler.Norwegian: "programvare",
	compiler.Pashto: "ساوتري",
	compiler.Persian: "نرمافزار",
	compiler.Polish: "oprogramowanie",
	compiler.Portuguese: "programas",
	compiler.Punjabi: "ਸਾਫਟਵੇਅਰ",
	
	compiler.Romanian: "software",
	compiler.Russian: "программногообеспечения",
	
	compiler.ScotsGaelic: "batharbog",
	compiler.Serbian: "софтвер",
	compiler.Sesotho: "software",
	compiler.Shona: "software",
	compiler.Sindhi: "سافٽويئر",
	compiler.Sinhala: "මෘදුකාංග",
	compiler.Slovak: "softvér",
	compiler.Slovenian: "programskoopremo",
	compiler.Somali: "software",
	compiler.Sundanese: "software",
	compiler.Swahili: "programu",
	compiler.Swedish: "programvara",
	
	compiler.Tajik: "нармафзор",
	compiler.Tamil: "மென்பொருள்",
	compiler.Telugu: "సాఫ్ట్వేర్",
	compiler.Thai: "ซอฟต์แวร์",
	compiler.Turkish: "yazılım",
	
	compiler.Ukrainian: "програмнезабезпечення",
	compiler.Urdu: "سافٹویئر",
	compiler.Uzbek: "dasturiytaminot",
	
	compiler.Vietnamese: "phầnmềm",
	
	compiler.Welsh: "meddalwedd",
	
	compiler.Xhosa: "software",
	
	compiler.Yiddish: "ווייכווארג",
	compiler.Yoruba: "software",
	
	compiler.Zulu: "isofthiwe",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Exit()
	},
}

var Exit = compiler.Statement {
	Name: compiler.Translatable{
		compiler.English: "exit",
	},
	 
	OnScan: func(c *compiler.Compiler) {
		c.Exit()
	},
}
 
var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Main()
		c.Expecting(symbols.CodeBlockBegin)

		c.GainScope()
		c.SetFlag(Flag)
		
		
		global.Init(c)
	},
}

var End = compiler.Statement {
	Name: compiler.NoTranslation(symbols.CodeBlockEnd),
	 
	OnScan: func(c *compiler.Compiler) {
		c.LoseScope()
	},
}
