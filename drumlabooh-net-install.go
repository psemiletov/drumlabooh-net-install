package main


import (
    "fmt"
    "os"
    "net/http"
    //"strconv"
    "io"
    "log"
    "runtime"
    "strings"
    "archive/zip"
    "path/filepath"
)



func read_url_as_string (url string) string {
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return (string(b))
	
}


func userHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    } else if runtime.GOOS == "linux" {
        home := os.Getenv("XDG_CONFIG_HOME")
        if home != "" {
            return home
        }
    }
    return os.Getenv("HOME")
}


func Unzip (src, dest, ver string) error {
    
  //  DRUMSKLAD_DIR := "drum_sklad-1.0.0"
    
    dest = filepath.Clean(dest) + string(os.PathSeparator)

    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        path := filepath.Join(dest, f.Name)
        // Check for ZipSlip: https://snyk.io/research/zip-slip-vulnerability
        if !strings.HasPrefix(path, dest) {
            return fmt.Errorf("%s: illegal file path", path)
        }
        
//        if (f.Name == "drum_sklad-1.0.0") {
 //           path = strings.Replace (path, f.Name, "drum_sklad", 1)
   //     }

        path = strings.Replace (path, "drum_sklad-" + ver, "drum_sklad", 1)
 
        rc, err := f.Open()
        if err != nil {
            return err
        }
        
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        if f.FileInfo().IsDir() {
                os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}


func downloadFile(filepath string, url string) (err error){

  // Create the file
  out, err := os.Create(filepath)
  if err != nil  {
    return err
  }
  defer out.Close()

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    return err
  }

  return nil
}


func main() {
    
    VER_LOCAL :="1" 
    
     
    
    var arg string
    
    if len(os.Args) == 2 {
        arg = os.Args[1]
        //fmt.Println ("AAAAAAAAAAAA")
        
    }   
    
    flag_test := false
    
    if (arg == "test"){
        flag_test = true
        fmt.Println ("TEST")
    }    
	//fmt.Println (exe_path)
    
    ver_remote := read_url_as_string ("https://raw.githubusercontent.com/psemiletov/drumlabooh-net-install/refs/heads/main/version.txt")
    
 /*   version_remote, err := strconv.Atoi(ver_remote)
    if err != nil {
        // ... handle error
        panic(err)
    } 
   */ 
 
 
    
    if (ver_remote == VER_LOCAL){
        fmt.Println ("Installer version is up-to-date")
    } else {
            fmt.Println ("Installer is updating... Wait")
        
            exe_path, err := os.Executable()

            fmt.Println (exe_path)

    
   	        if err != nil {
	  	       log.Fatalln(err)
	        }
	
        
           fmt.Println ("Done")
           fmt.Println ("Please restart Installer")
        
           return
    }
        
    
 //   source_path_to_binary := ""
 //   binary_url := ""
   /*
    if (ver_remote > VER_LOCAL){
        
        //update binary
        fmt.Println ("We need to update the binary") 
        //Dowloading 
         // downloadFile (source_path_to_lv2_zip, lv2_url)
        
        //exit
        return
    }    
     */      
           
    
    /*
    req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/psemiletov/drumlabooh/refs/heads/main/version.txt", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
    */
    
    fmt.Println ("Drumlabooh Net Installer " + VER_LOCAL)
    
    
    home_dir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal( err )
    }
    
    fmt.Println ("home_dir:" + home_dir)
    
    
    tempdir, err := os.MkdirTemp("", "laboohtempdir")

    if err != nil {
        log.Fatal( err )
    }
    
    
    fmt.Println("Temp dir name:", tempdir)
    
    labooh_ver := read_url_as_string ("https://raw.githubusercontent.com/psemiletov/drumlabooh/refs/heads/main/version.txt")
    kits_ver := read_url_as_string ("https://raw.githubusercontent.com/psemiletov/drum_sklad/refs/heads/main/version.txt")

    fmt.Println ("Install/update Drumlabooh v.", labooh_ver)
    
    
    
    
//    lv2_url := "https://github.com/psemiletov/drumlabooh/releases/download/6.0.0/drumlabooh.lv2.zip"
//    vst_url := "https://github.com/psemiletov/drumlabooh/releases/download/6.0.0/drumlabooh.vst3.zip"
    
    //lv2_url := "https://github.com/psemiletov/drumlabooh/releases/download/" + string(b) + "/drumlabooh.lv2.zip"
    //vst_url := "https://github.com/psemiletov/drumlabooh/releases/download/" + string(b) + "/drumlabooh.vst3.zip"
    
    lv2_url := "https://github.com/psemiletov/drumlabooh/releases/download/" + labooh_ver + "/drumlabooh.lv2.zip"
    vst_url := "https://github.com/psemiletov/drumlabooh/releases/download/" + labooh_ver + "/drumlabooh.vst3.zip"
    
    
    drumkits_url := "https://github.com/psemiletov/drum_sklad/archive/refs/tags/"+ kits_ver + ".zip"
    
    source_path_to_lv2_zip := tempdir + "/labooh_lv2.zip"
    source_path_to_vst_zip := tempdir + "/labooh_vst.zip"
    source_path_to_drum_sklad := tempdir + "/drum_sklad.zip"
    
    
    dest_lv2_path := home_dir + "/.lv2"
    dest_vst_path := home_dir + "/.vst3"
//    dest_drumsklad_path := home_dir + "/drum_skladT"
    dest_drumsklad_path := home_dir;
  
    
    
    
    fmt.Println ("Downloading LV2 to " + source_path_to_lv2_zip)
    if (! flag_test){
        downloadFile (source_path_to_lv2_zip, lv2_url)
    } 
    
    fmt.Println ("Downloading VST3i to " + source_path_to_vst_zip)
    if (! flag_test){
        downloadFile (source_path_to_vst_zip, vst_url)
    } 

     fmt.Println ("Downloading kits to " + source_path_to_drum_sklad)

    if (! flag_test){
       downloadFile (source_path_to_drum_sklad, drumkits_url)
    }   

    
    //downloadFile("5.0.0.zip", "https://github.com/psemiletov/drumlabooh/archive/refs/tags/5.0.0.zip")
    //Unzip("5.0.0.zip", ".lv2/")

    fmt.Println ("Unpacking LV2 to " + dest_lv2_path)

    if (! flag_test){
         Unzip(source_path_to_lv2_zip, dest_lv2_path, kits_ver)
    } 
    
    fmt.Println ("Unpacking VST3i to " + dest_vst_path)
    
    if (! flag_test){
        Unzip(source_path_to_vst_zip, dest_vst_path, kits_ver)
    } 

    fmt.Println ("Unpacking drumkits to subdir drum_sklad at " + dest_drumsklad_path + "/drum_sklad")

    if (! flag_test){
       Unzip (source_path_to_drum_sklad, dest_drumsklad_path, kits_ver)
    }
   
    //remove all archives
    fmt.Println ("Removing temp dir: " + tempdir)
    
    os.RemoveAll (tempdir)
    
}
