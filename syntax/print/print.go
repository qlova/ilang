package print

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/types/list"
import "github.com/qlova/ilang/types/array"
import "github.com/qlova/ilang/types/thing"
import "github.com/qlova/ilang/types/concept"
import "github.com/qlova/ilang/types/function"

import "fmt"
import "io/ioutil"

var Name = compiler.Translatable{
	compiler.English: "print",
	compiler.Dutch: "afdrukken",
	compiler.Maori: "perehitia",
	compiler.French: "imprime",
	compiler.Spanish: "imprimir",
	compiler.Chinese: "打印",
	compiler.Japanese: "プリント",
	compiler.Samoan: "lolomi",
	compiler.German: "drucke",
	
	compiler.Afrikaans: "druk",
	compiler.Albanian: "printoni",
	compiler.Amharic: "አትም",
	compiler.Arabic: "اطبع",
	compiler.Armenian: "Տպել",
	compiler.Azerbaijani: "çap",
	
	compiler.Basque: "inprimatu",
	compiler.Belarusian: "друк",
	compiler.Bengali: "মুদ্রণ",
	compiler.Bosnian: "štampanje",
	compiler.Bulgarian: "печат",
	compiler.Burmese: "ပုံနှိပ်",
	
	compiler.Catalan: "imprimiu",
	compiler.Cebuano: "imprinta",
	compiler.Chichewa: "sindikizani",
	compiler.Corsican: "stampa",
	compiler.Croatian: "ispisati",
	compiler.Czech: "tisk",
	
	compiler.Danish: "print",
	
	compiler.Esperanto: "presi",
	compiler.Estonian: "trüki",
	
	compiler.Filipino: "print",
	compiler.Finnish: "tulostettava",
	compiler.Frisian: "ôfdrukke",
	
	compiler.Galician: "imprimir",
	compiler.Georgian: "ბეჭდვითი",
	compiler.Greek: "Τυπώνω",
	compiler.Gujarati: "પ્રિન્ટ",
	
	compiler.HaitianCreole: "enprime",
	compiler.Hausa: "buga",
	compiler.Hawaiian: "papa",
	compiler.Hebrew: "הדפס",
	compiler.Hindi: "प्रिंट",
	compiler.Hmong: "sau",
	compiler.Hungarian: "nyomtatott",
	
	compiler.Icelandic: "prenta",
	compiler.Igbo: "ebipụta",
	compiler.Indonesian: "mencetak",
	compiler.Irish: "priontáil",
	compiler.Italian: "stampare",
	
	compiler.Javanese: "print",
	
	compiler.Kannada: "ಮುದ್ರಣ",
	compiler.Kazakh: "басыпшығару",
	compiler.Khmer: "បោះពុម្ព",
	compiler.Klingon: "SevIr",
	compiler.Korean: "인쇄",
	compiler.Kurdish: "çap",
	compiler.Kyrgyz: "басып",
	
	compiler.Lao: "ພິມ",
	compiler.Latin: "print",
	compiler.Latvian: "drukas",
	compiler.Lithuanian: "spausdinimo",
	compiler.Luxembourgish: "drucken",
	
	compiler.Macedonian: "печатење",
	compiler.Malagasy: "print",
	compiler.Malay: "cetak",
	compiler.Malayalam: "അച്ചടി",
	compiler.Maltese: "stampa",
	compiler.Marathi: "प्रिंट",
	compiler.Mongolian: "хэвлэх",
	
	compiler.Nepali: "प्रिन्ट",
	compiler.Norwegian: "skrive",
	
	compiler.Pashto: "چاپ",
	compiler.Persian: "چاپ",
	compiler.Polish: "wydrukuj",
	compiler.Portuguese: "impressão",
	compiler.Punjabi: "ਛਾਪੋ",
	
	compiler.Romanian: "tipăriți",
	compiler.Russian: "Распечатать",
	
	compiler.ScotsGaelic: "clò",
	compiler.Serbian: "штампани",
	compiler.Sesotho: "hatisa",
	compiler.Shona: "print",
	compiler.Sindhi: "پرنٽ",
	compiler.Sinhala: "මුද්රිත",
	compiler.Slovak: "tlačový",
	compiler.Slovenian: "tiskalni",
	compiler.Somali: "daabacan",
	compiler.Sundanese: "nyitak",
	compiler.Swahili: "magazeti",
	compiler.Swedish: "skriva",
	
	compiler.Tajik: "чоп",
	compiler.Tamil: "அச்சு",
	compiler.Telugu: "ప్రింట్",
	compiler.Thai: "พิมพ์",
	compiler.Turkish: "baskı",
	
	compiler.Ukrainian: "друкований",
	compiler.Urdu: "پرنٹکریں",
	compiler.Uzbek: "bosma",
	
	compiler.Vietnamese: "in",
	
	compiler.Welsh: "argraffu",
	
	compiler.Xhosa: "lokuprinta",
	
	compiler.Yiddish: "דרוק",
	compiler.Yoruba: "tẹjade",
	
	compiler.Zulu: "phrinta",
}

func PrintType(c *compiler.Compiler, t compiler.Type) {
	
	if list.Is(t) || array.Is(t) {
		
		if thing.Is(list.SubType(t)) {
			c.Unimplemented()
		}
		
		c.CopyPipe()
		c.List()
		c.Int('[')
		c.Put()
		c.Send()
		
		c.Int(0)
		c.Loop()
			c.Copy()
			c.Size()
			c.Same() 
			c.If()
				c.Done()
			c.No()
		
			c.Copy()
			c.Get()
			
			if list.SubType(t).Base != compiler.INT && !list.SubType(t).Equals(text.Type) {
				print(list.SubType(t).String())
				c.Unimplemented()
			}
			
			//Add beginning quote.
			if  list.SubType(t).Equals(text.Type) {
				
				c.CopyPipe()
				c.List()
				c.Int('"')
				c.Put()
				c.Send()
				
				c.HeapList()
			}
			
			PrintType(c, list.SubType(t))
			
			if  list.SubType(t).Equals(text.Type) {
				c.CopyPipe()
				c.List()
				c.Int('"')
				c.Put()
				c.Send()
			}
			
			c.Copy()
			c.Size()
			c.Int(1)
			c.Sub()
			c.More()
			c.If()
				c.CopyPipe()
				c.List()
				c.Int(',')
				c.Put()
				c.Send()
			c.No()
			
			c.Int(1)
			c.Add()
		c.Redo()
		
		c.CopyPipe()
		c.List()
		c.Int(']')
		c.Put()
		c.Send()
		
		return
	}
	c.Cast(t, text.Type)
	
	c.CopyPipe()
	c.Send()
}


var Expression = compiler.Expression {
	Name: Name,
}

var Tmp int

func init() {
	Expression.OnScan = func(c *compiler.Compiler) compiler.Type {
		
		switch c.Peek() {
			
			//Testcase.
			case symbols.IndexBegin:
				c.Scan()
				
				var args []string
				var types  []*compiler.Type
				for {
					var v = c.Scan()
					if c.GetType(v) != nil {
						types = append(types, c.GetType(v))
						args = append(args, "")
					} else {
						types = append(types, nil)
						args = append(args, v)
					}
					if c.ScanIf(symbols.IndexEnd) {
						break
					}
					c.Expecting(symbols.ArgumentSeperator)
				}
				Tmp++
				
				
				//TODO cache this.
				
				var fargs []compiler.Type
				
				
				var out = c.Output 
				c.Output = ioutil.Discard
				
				var tmp = text.Tmp
				
				//Sift out any text literals. Swappy madness.
				for i := 0; i < len(args); i++ {
					if types[i] == nil {
							
							//DO some hacker level business here.
							var cache compiler.Cache
							cache.Write([]byte(args[i]))
							cache.Write([]byte{')'})
							
							c.LoadCache(cache, "print.go", 0)
							
							PrintType(c, c.ScanExpression())
							c.Expecting(")")
					}
				}
				text.Tmp = tmp
				
				c.SwapOutput()
				
				c.Code("print_"+fmt.Sprint(Tmp))
				
					c.List()
					c.Open()
				
					for i := 0; i < len(args); i++ {						
						if types[i] == nil {
							
							//DO some hacker level business here.
							var cache compiler.Cache
							cache.Write([]byte(args[i]))
							cache.Write([]byte{')'})
							
							c.LoadCache(cache, "print.go", 0)
							
							PrintType(c, c.ScanExpression())
							c.Expecting(")")
							
							
						} else {
							fargs = append(fargs, *types[i])
							PrintType(c, *types[i])
						}
					}
					
					c.List()
					c.Int('\n')
					c.Put()
					c.Send()
				
				c.Back()
				c.SwapOutput()
				
				c.Output = out
				
				c.Wrap("print_"+fmt.Sprint(Tmp))
				
				return function.Type.With(function.Data{
					Arguments: fargs,
				})
			
			default:
				return concept.Type.With(concept.Data{
					Statement: Statement,
					Expression: Expression,
				})
		}
		
		return compiler.Type{Fake: true}
	}
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.FunctionCallBegin)
		
		c.List()
		c.Open()
		
		if c.ScanIf(symbols.FunctionCallEnd) {
			//Print Newline
			c.List()
			c.Int('\n')
			c.Put()
			c.Send()
			return
		}
		
		for {
			PrintType(c, c.ScanExpression())
			
			switch c.Scan() {
				case symbols.ArgumentSeperator:
					continue
				case symbols.FunctionCallEnd:
					//Print Newline
					c.List()
					c.Int('\n')
					c.Put()
					c.Send()
					
					return
				default:
					c.Expected(symbols.ArgumentSeperator, symbols.FunctionCallEnd)
			}
		}
	},
}

