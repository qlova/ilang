package targets

import "path"
import "os"
import "os/exec"
import "strings"
import "errors"

type Android struct {}

func (Android) Compile(mainFile string) error { 
	
	//Copy in MainActivity.
	var SDK = Home()+"/.ilang/targets/android"
	var base = path.Base(mainFile[:len(mainFile)-2])
	var safebase = strings.Replace(base, " ", "", -1)
	
	os.Rename("Stack.java", "Stack.android")
	CopyFile(base+".android", "MainActivity.java", "package nz.co.qlova.ilang."+safebase+";\n")
	CopyFile("Stack.android", "Stack.java", "package nz.co.qlova.ilang."+safebase+";\n")
	os.Mkdir("bin", 0755)
	
	//Generate AndroidManifest
	manifest, err := os.Create("AndroidManifest.xml")
	if err != nil {
		return err
	}
	
	
	
	_, err = manifest.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="nz.co.qlova.ilang.`+safebase+`">
    
    <uses-permission android:name="android.permission.INTERNET" />
    <uses-permission android:name="android.permission.READ_PHONE_STATE" />
    <uses-permission android:name="android.permission.ACCESS_WIFI_STATE" />

    <application
        android:allowBackup="true"
        android:icon="@drawable/ic_launcher"
        android:label="`+base+`"
		android:theme="@android:style/Theme.NoTitleBar">
		
        <activity android:name=".MainActivity">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />

                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>

</manifest>
	`))
	if err != nil {
		return err
	}
	
	manifest.Close()
	
	cmd := exec.Command("javac", "-classpath", SDK+"/sdk.jar", "-sourcepath", "src:gen", "-d", "bin", 
							"-target", "1.7", "-source" ,"1.7", "MainActivity.java", "Stack.java")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	
	cmd = exec.Command(SDK+"/dx", "--dex", "--output=classes.dex", "bin")
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	//Move resources into the correct places.
	Unzip(SDK+"/res.zip", "./res")
	
	cmd = exec.Command(SDK+"/aapt", "package", "-f", "-M", "AndroidManifest.xml" ,"-S", "res", "-I", SDK+"/sdk.jar", "-F", base+".apk.unaligned")
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	cmd = exec.Command(SDK+"/aapt", "add", base+".apk.unaligned", "classes.dex")	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
		
	cmd = exec.Command("jarsigner", "-keystore", SDK+"/debug.keystore", "-storepass", "android", base+".apk.unaligned", "androiddebugkey")	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	cmd = exec.Command(SDK+"/zipalign", "-f", "4", base+".apk.unaligned", base+"-debug.apk")	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	return nil
}
func (Android) Run(mainFile string) error {
	var base = path.Base(mainFile[:len(mainFile)-2])
	var safebase = strings.Replace(base, " ", "", -1)
	var packagename = "nz.co.qlova.ilang."+safebase
	
	cmd := exec.Command("adb", "install", "-r", base+"-debug.apk")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	cmd = exec.Command("adb", "shell", "am", "start", "-n", packagename+"/"+packagename+".MainActivity")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	return nil
}
func (Android) Export(mainFile string) error { 
	
	return errors.New("Exporting not enabled for android!")
}

func init() {
	RegisterTarget("android", Android{})
}
