package importation

import "github.com/qlova/uct/compiler"
import "os"
import "os/user"
import "path/filepath"

const Debug = false

var Name = compiler.Translatable{
	compiler.English: "import",
}

var FileDepth int
var Dirs []string

var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		pkg := c.Scan()
		for c.Scan() == "." {
			pkg += "/"+c.Scan()
		}

		dir, _ := os.Getwd()
		if Debug {
			println("importing, ", pkg, " in ", dir)
		}

		var filename = ""

		retry:
		file, err := os.Open(pkg+".i")
		if err != nil {
			if file, err = os.Open(pkg+"/"+filepath.Base(pkg)+".i"); err != nil {
				
				//Search through parent folders?
				dir, _ := os.Getwd()
				if FileDepth > 0 {
					FileDepth--
					os.Chdir(Dirs[len(Dirs)-1])
					Dirs = Dirs[:len(Dirs)-1]
					goto retry
				}
				
				//Search in ~/.ilang.
				if FileDepth == 0 {
					FileDepth--
					usr, err := user.Current()
					if err == nil {
						os.Chdir(usr.HomeDir+"/.ilang/imports/")
						Dirs = append(Dirs, usr.HomeDir+"/.ilang/imports/")
						goto retry
					}
				}
				
				c.RaiseError(compiler.Translatable{
					compiler.English: "Cannot import "+pkg+", does not exist! "+dir,
				})
				
			} else {
					filename = pkg+"/"+pkg+".i"
					
					dir, _ := os.Getwd()
					
				Dirs= append(Dirs, dir)
					
				err := os.Chdir("./"+pkg)
				if err != nil {
					c.RaiseError(compiler.Translatable{
						compiler.English: err.Error(),
					})
				}
				
					dir, _ = os.Getwd()
					if Debug {
					println("moved to ", dir)
					}
				
				FileDepth++
			}
		} else {
			filename = pkg+".i"
		}
		
		c.AddInput(file)
		c.Scanners[len(c.Scanners)-1].Filename = filename
	},
}
