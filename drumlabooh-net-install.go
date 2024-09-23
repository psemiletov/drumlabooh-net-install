package main


import (
    "fmt"
    "os"
    "net/http"
    "io"
    "log"
    "runtime"
    "strings"
    "archive/zip"
    "path/filepath"
)



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


func Unzip(src, dest string) error {
    
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

        path = strings.Replace (path, "drum_sklad-1.0.0", "drum_sklad", 1)
 

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
    
    fmt.Println ("Drumlabooh Net Installer 6.0.0")
    
    
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
    
    lv2_url := "https://github.com/psemiletov/drumlabooh/releases/download/6.0.0/drumlabooh.lv2.zip"
    vst_url := "https://github.com/psemiletov/drumlabooh/releases/download/6.0.0/drumlabooh.vst3.zip"
    drumkits_url := "https://github.com/psemiletov/drum_sklad/archive/refs/tags/1.0.0.zip"
    
    
    source_path_to_lv2_zip := tempdir + "/labooh_lv2.zip"
    source_path_to_vst_zip := tempdir + "/labooh_vst.zip"
    source_path_to_drum_sklad := tempdir + "/drum_sklad.zip"
    
    
    dest_lv2_path := home_dir + "/.lv2"
    dest_vst_path := home_dir + "/.vst3"
//    dest_drumsklad_path := home_dir + "/drum_skladT"
    dest_drumsklad_path := home_dir;
  
    
    
    
    fmt.Println ("Downloading LV2 to " + source_path_to_lv2_zip)
    downloadFile (source_path_to_lv2_zip, lv2_url)
    
    fmt.Println ("Downloading VST3i to " + source_path_to_vst_zip)
    downloadFile (source_path_to_vst_zip, vst_url)
        
    fmt.Println ("Downloading kits to " + source_path_to_drum_sklad)
    downloadFile (source_path_to_drum_sklad, drumkits_url)
    

    
    //downloadFile("5.0.0.zip", "https://github.com/psemiletov/drumlabooh/archive/refs/tags/5.0.0.zip")
    //Unzip("5.0.0.zip", ".lv2/")

    fmt.Println ("Unpacking LV2 to " + dest_lv2_path)
    Unzip(source_path_to_lv2_zip, dest_lv2_path)
    
    fmt.Println ("Unpacking VST3i to " + dest_vst_path)
    Unzip(source_path_to_vst_zip, dest_vst_path)

    fmt.Println ("Unpacking drumkits to to " + dest_drumsklad_path)
    Unzip (source_path_to_drum_sklad, dest_drumsklad_path)
   
    //remove all archives
    fmt.Println ("Removing temp dir: " + tempdir)
    
    os.RemoveAll (tempdir)
    
}