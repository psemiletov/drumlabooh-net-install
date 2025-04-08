package main

import (
    "fmt"
    "os"
    "net/http"
    "os/user"
    "io"
    "log"
    "runtime"
    "strings"
    "archive/zip"
    "path/filepath"
)

func read_url_as_string(url string) string {
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
    return string(b)
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

type ProgressReader struct {
    Reader    io.Reader
    Total     int64
    Downloaded int64
    Name      string
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
    n, err = pr.Reader.Read(p)
    pr.Downloaded += int64(n)
    
    // Проверка на случай, если Total <= 0
    var percent float64
    if pr.Total > 0 {
        percent = float64(pr.Downloaded) / float64(pr.Total) * 100
    }
    
    // Ограничиваем percent от 0 до 100
    if percent > 100 {
        percent = 100
    }
    if percent < 0 {
        percent = 0
    }
    
    barLength := int(percent / 2)
    if barLength < 0 {
        barLength = 0
    }
    
    fmt.Printf("\rDownloading %s: [%-50s] %.2f%%", 
        pr.Name,
        strings.Repeat("#", barLength),
        percent)
    if pr.Downloaded >= pr.Total && pr.Total > 0 {
        fmt.Println()
    }
    return
}

func downloadFile(filepath string, url string) error {
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    totalSize := resp.ContentLength
    progressReader := &ProgressReader{
        Reader:    resp.Body,
        Total:     totalSize,
        Name:      filepath[strings.LastIndex(filepath, "/")+1:],
    }
    
    _, err = io.Copy(out, progressReader)
    if err != nil {
        return err
    }
    return nil
}

func Unzip(src, dest, ver string) error {
    dest = filepath.Clean(dest) + string(os.PathSeparator)
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()

    os.MkdirAll(dest, os.ModePerm)
    
    totalFiles := len(r.File)
    currentFile := 0

    extractAndWriteFile := func(f *zip.File) error {
        currentFile++
        path := filepath.Join(dest, f.Name)
        if !strings.HasPrefix(path, dest) {
            return fmt.Errorf("%s: illegal file path", path)
        }
        
        path = strings.Replace(path, "drum_sklad-"+ver, "drum_sklad", 1)
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        // Проверка на случай, если totalFiles == 0
        var percent float64
        if totalFiles > 0 {
            percent = float64(currentFile) / float64(totalFiles) * 100
        }
        
        // Ограничиваем percent от 0 до 100
        if percent > 100 {
            percent = 100
        }
        if percent < 0 {
            percent = 0
        }
        
        barLength := int(percent / 2)
        if barLength < 0 {
            barLength = 0
        }

        fmt.Printf("\rUnpacking %s: [%-50s] %.2f%% (%d/%d)", 
            src[strings.LastIndex(src, "/")+1:],
            strings.Repeat("#", barLength),
            percent,
            currentFile,
            totalFiles)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, os.ModePerm)
        } else {
            os.MkdirAll(filepath.Dir(path), os.ModePerm)
            outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer outFile.Close()
            _, err = io.Copy(outFile, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        if err := extractAndWriteFile(f); err != nil {
            return err
        }
    }
    fmt.Println()
    return nil
}

func isRoot() bool {
    currentUser, err := user.Current()
    if err != nil {
        log.Fatalf("[isRoot] Unable to get current user: %s", err)
    }
    return currentUser.Username == "root"
}

func main() {
    VER_LOCAL := "4"
    
    if isRoot() {
        fmt.Println("Please run as non-root")
        return
    }
        
    var arg string
    if len(os.Args) == 2 {
        arg = os.Args[1]
    }
    
    flag_test := false
    if arg == "test" {
        flag_test = true
        fmt.Println("TEST")
    }
        
    fmt.Println("Drumlabooh Net Installer " + VER_LOCAL)
    
    home_dir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("home_dir:" + home_dir)
    
    tempdir, err := os.MkdirTemp("", "laboohtempdir")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Temp dir name:", tempdir)
    
    labooh_ver := read_url_as_string("https://raw.githubusercontent.com/psemiletov/drumlabooh/refs/heads/main/version.txt")
    kits_ver := read_url_as_string("https://raw.githubusercontent.com/psemiletov/drum_sklad/refs/heads/main/version.txt")

    fmt.Println("Install/update Drumlabooh v.", labooh_ver)
    fmt.Println("Install/update kits v.", kits_ver)
    
    lv2_url := "https://github.com/psemiletov/drumlabooh/releases/download/" + labooh_ver + "/drumlabooh.lv2.zip"
    vst_url := "https://github.com/psemiletov/drumlabooh/releases/download/" + labooh_ver + "/drumlabooh.vst3.zip"
    drumkits_url := "https://github.com/psemiletov/drum_sklad/archive/refs/tags/" + kits_ver + ".zip"
    
    source_path_to_lv2_zip := tempdir + "/labooh_lv2.zip"
    source_path_to_vst_zip := tempdir + "/labooh_vst.zip"
    source_path_to_drum_sklad := tempdir + "/drum_sklad.zip"
    
    dest_lv2_path := home_dir + "/.lv2"
    dest_vst_path := home_dir + "/.vst3"
    dest_drumsklad_path := home_dir
    
    if (flag_test){
        dest_lv2_path := dest_lv2_path + "TEST" 
        dest_vst_path := dest_vst_path + "TEST"
        dest_drumsklad_path := "/STEST"
    }
    
    
    fmt.Println("Downloading LV2 to " + source_path_to_lv2_zip)
    if !flag_test {
        downloadFile(source_path_to_lv2_zip, lv2_url)
    }
    
    fmt.Println("Downloading VST3i to " + source_path_to_vst_zip)
    if !flag_test {
        downloadFile(source_path_to_vst_zip, vst_url)
    }

    fmt.Println("Downloading kits to " + source_path_to_drum_sklad)
    if !flag_test {
        downloadFile(source_path_to_drum_sklad, drumkits_url)
    }

    fmt.Println("Unpacking LV2 from " + source_path_to_lv2_zip)
    fmt.Println("Unpacking LV2 to " + dest_lv2_path)
    if !flag_test {
        Unzip(source_path_to_lv2_zip, dest_lv2_path, kits_ver)
    }
    
    fmt.Println("Unpacking VST3i from " + source_path_to_vst_zip)
    fmt.Println("Unpacking VST3i to " + dest_vst_path)
    if !flag_test {
        Unzip(source_path_to_vst_zip, dest_vst_path, kits_ver)
    }

    fmt.Println("Unpacking drumkits to subdir drum_sklad at " + dest_drumsklad_path + "/drum_sklad")
    if !flag_test {
        Unzip(source_path_to_drum_sklad, dest_drumsklad_path, kits_ver)
    }
   
    fmt.Println("Removing temp dir: " + tempdir)
    os.RemoveAll(tempdir)
}