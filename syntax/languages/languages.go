package languages

import "github.com/qlova/uct/compiler"

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		switch c.Token() {
			case "English":
				
				c.Language = compiler.English
				
				return true
			
			case "Māori", "Maori":
				
				c.Language = compiler.Maori
				
				return true
				
			case "Dutch":
				
				c.Language = compiler.Dutch
				
				return true
				
			case "French":
				
				c.Language = compiler.French
				
				return true
			
			case "Spanish":
				
				c.Language = compiler.Spanish
				
				return true
				
			case "Chinese":
				
				c.Language = compiler.Chinese
				
				return true
				
			case "Japanese":
				
				c.Language = compiler.Japanese
				
				return true
				
				
			case "Samoan":
				
				c.Language = compiler.Samoan
				
				return true
			
			case "German":
				
				c.Language = compiler.German
				
				return true
				
			case "Afrikaans": c.Language = compiler.Afrikaans; return true
			case "Albanian": c.Language = compiler.Albanian; return true
			case "Amharic": c.Language = compiler.Amharic; return true
			case "Arabic": c.Language = compiler.Arabic; return true
			case "Armenian": c.Language = compiler.Armenian; return true
			case "Azerbaijani": c.Language = compiler.Azerbaijani; return true
				
			case "Basque": c.Language = compiler.Basque; return true
			case "Belarusian": c.Language = compiler.Belarusian; return true
			case "Bengali": c.Language = compiler.Bengali; return true
			case "Bosnian": c.Language = compiler.Bosnian; return true
			case "Bulgarian": c.Language = compiler.Bulgarian; return true
			case "Burmese": c.Language = compiler.Burmese; return true
				
			case "Catalan": c.Language = compiler.Catalan; return true
			case "Cebuano": c.Language = compiler.Cebuano; return true
			case "Chichewa": c.Language = compiler.Chichewa; return true
			case "Corsican": c.Language = compiler.Corsican; return true
			case "Croatian": c.Language = compiler.Croatian; return true
			case "Czech": c.Language = compiler.Czech; return true
				
			case "Danish": c.Language = compiler.Danish; return true
				
			case "Esperanto": c.Language = compiler.Esperanto; return true
			case "Estonian": c.Language = compiler.Estonian; return true
				
			case "Filipino": c.Language = compiler.Filipino; return true
			case "Finnish": c.Language = compiler.Finnish; return true
			case "Frisian": c.Language = compiler.Frisian; return true
				
			case "Galician": c.Language = compiler.Galician; return true
			case "Georgian": c.Language = compiler.Georgian; return true
			case "Greek": c.Language = compiler.Greek; return true
			case "Gujarati": c.Language = compiler.Gujarati; return true
			
			case "HaitianCreole": c.Language = compiler.HaitianCreole; return true
			case "Hausa": c.Language = compiler.Hausa; return true
			case "Hawaiian": c.Language = compiler.Hawaiian; return true
			case "Hebrew": c.Language = compiler.Hebrew; return true
			case "Hindi": c.Language = compiler.Hindi; return true
			case "Hmong": c.Language = compiler.Hmong; return true
			case "Hungarian": c.Language = compiler.Hungarian; return true
				
			case "Icelandic": c.Language = compiler.Icelandic; return true
			case "Igbo": c.Language = compiler.Igbo; return true
			case "Indonesian": c.Language = compiler.Indonesian; return true
			case "Irish": c.Language = compiler.Irish; return true
			case "Italian": c.Language = compiler.Italian; return true
				
			case "Javanese": c.Language = compiler.Javanese; return true
				
			case "Kannada": c.Language = compiler.Kannada; return true
			case "Kazakh": c.Language = compiler.Kazakh; return true
			case "Khmer": c.Language = compiler.Khmer; return true
			case "Klingon": c.Language = compiler.Klingon; return true
			case "Korean": c.Language = compiler.Korean; return true
			case "Kurdish": c.Language = compiler.Kurdish; return true
			case "Kyrgyz": c.Language = compiler.Kyrgyz; return true
				
			case "Lao": c.Language = compiler.Lao; return true
			case "Latin": c.Language = compiler.Latin; return true
			case "Latvian": c.Language = compiler.Latvian; return true
			case "Lithuanian": c.Language = compiler.Lithuanian; return true
			case "Luxembourgish": c.Language = compiler.Luxembourgish; return true
				
			case "Macedonian": c.Language = compiler.Macedonian; return true
			case "Malagasy": c.Language = compiler.Malagasy; return true
			case "Malay": c.Language = compiler.Malay; return true
			case "Malayalam": c.Language = compiler.Malayalam; return true
			case "Maltese": c.Language = compiler.Maltese; return true
			case "Marathi": c.Language = compiler.Marathi; return true
			case "Mongolian": c.Language = compiler.Mongolian; return true
				
			case "Nepali": c.Language = compiler.Nepali; return true
			case "Norwegian": c.Language = compiler.Norwegian; return true
				
			case "Pashto": c.Language = compiler.Pashto; return true
			case "Persian": c.Language = compiler.Persian; return true
			case "Polish": c.Language = compiler.Polish; return true
			case "Portuguese": c.Language = compiler.Portuguese; return true
			case "Punjabi": c.Language = compiler.Punjabi; return true
				
			case "Romanian": c.Language = compiler.Romanian; return true
			case "Russian": c.Language = compiler.Russian; return true
			
			case "ScotsGaelic": c.Language = compiler.ScotsGaelic; return true
			case "Serbian":  c.Language = compiler.Serbian; return true
			case "Sesotho": c.Language = compiler.Sesotho; return true
			case "Shona": c.Language = compiler.Shona; return true
			case "Sindhi": c.Language = compiler.Sindhi; return true
			case "Sinhala": c.Language = compiler.Sinhala; return true
			case "Slovak": c.Language = compiler.Slovak; return true
			case "Slovenian": c.Language = compiler.Slovenian; return true
			case "Somali": c.Language = compiler.Somali; return true
			case "Sundanese": c.Language = compiler.Sundanese; return true
			case "Swahili": c.Language = compiler.Swahili; return true
			case "Swedish": c.Language = compiler.Swedish; return true
				
			case "Tajik": c.Language = compiler.Tajik; return true
			case "Tamil": c.Language = compiler.Tamil; return true
			case "Telugu": c.Language = compiler.Telugu; return true
			case "Thai": c.Language = compiler.Thai; return true
			case "Turkish": c.Language = compiler.Turkish; return true
				
			case "Ukrainian": c.Language = compiler.Ukrainian; return true
			case "Urdu": c.Language = compiler.Urdu; return true
			case "Uzbek": c.Language = compiler.Uzbek; return true
				
			case "Vietnamese": c.Language = compiler.Vietnamese; return true
				
			case "Welsh": c.Language = compiler.Welsh; return true
			case "Xhosa": c.Language = compiler.Xhosa; return true
				
			case "Yiddish": c.Language = compiler.Yiddish; return true
			case "Yoruba": c.Language = compiler.Yoruba; return true
			
			case "Zulu": c.Language = compiler.Zulu; return true
				
			default:
				return false
		}
	},
}
