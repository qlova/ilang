/* 
 * This package contains help for the itool.
 * Ideally, the itool should be able to teach the user how to program in 'i'
 */
package main

import _ "github.com/qlova/ilang/src/modules/all"
import "github.com/qlova/ilang/src"
import "os"
import "strings"
import "sort"
import "fmt"


func beginner() {
	
}

func experienced() {
	
}	

func help() {
	if len(os.Args) == 1 {
		//Enter interactive mode.
		println("Welcome to itool!")
		println("(You can view this again by running it help)")
		println()
		println("Everything you need in order to learn 'i' is available in the help section.")
		println("What best describes you?")
		println("(1) Little to no experience in programming, wanting to learn.")
		println("(2) Experienced in programming, want to learn how 'i' works.")
		var choice int
		fmt.Scanln(&choice)
		
		if choice == 1 {
			println("Great! this tutorial will run you through the basics of programming.\n[Press enter to continue]")
			fmt.Scanln(&choice)
			beginner()
		}
		if choice == 2 {
			println("Great! this tutorial will bring you up to speed on the subtleties of 'i'.\n[Press enter to continue]")
			fmt.Scanln(&choice)
			experienced()
		}
		
		return
	}
	
	var topic = os.Args[1]
	
	//Should just be stored in a map.
	switch topic {
		case "builtin":
			ic := ilang.NewCompiler(nil)
			var list []string
			for name := range ic.DefinedFunctions {
					if !strings.Contains(name, "_") && !strings.Contains(name, ".") {
						list = append(list, name)
					}
			}
			sort.Strings(list)
			for _, value := range list {
					println(value)
			}
		
		case "tokens":
			var list []string
			for name := range ilang.EnglishTokens {
					if !strings.Contains(name, "_") && !strings.Contains(name, ".") {
						list = append(list, name)
					}
			}
			sort.Strings(list)
			for _, value := range list {
					println(value)
			}
	}
}
