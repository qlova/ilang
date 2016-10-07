
gui {
	<row>
		<button id="open">Open</button>
		<right>
			<text id="*">*</text>
			<button id="save">Save As</button>
		</right>
	</row>
	
	//Manual sizing of the texteditor widget.
	<div height="(100% - 15px)">
		<texteditor/>
	</div>
}

software {
	watch("open")
	watch("save")
	watch("ctrl+s")
	
	var filename = ""
	
	loop {
		! var event = grab("event")
		issues {
			break
		}
		
		if event = "open"
			var newfile = grab("file")
			if newfile != ""
				filename = newfile
				! var file = open(newfile)
				if error/0
					var data = ""
					loop {
					
						! data = data + reada@file('\n') + "\n"
						issues {
							break
						}
					}
					edit("texteditor.setValue()", data)
					close(file)
				end
			end
			
		elseif (event = "save") + (event = "ctrl+s")
		
			if (filename = "") + (event = "save")
				filename = grab("filename")
			end
			if filename != ""
		
				var data = grab("texteditor.getValue()")
			
				delete(filename)
				var file = open(filename)
				! output@file(data)
				issues {
					print("failed to write file!")
				}
				close(file)
			
				edit("*", "")
			end
		end
	}
}
