package main


import (
    "fmt"
    "os"
    "flag"
    "io/ioutil"
    "io"
    "strconv"
    "time"
)


func main() {

    cwd, _ := os.Getwd()


    fmt.Println("Simple image sorter started:")

    var srcDir = flag.String("src", cwd + "/DCIM", "Set the source directory from which to pull the images from.")

    flag.Parse()

    var imgList, err = ioutil.ReadDir(*srcDir);

    if err != nil {
        return
    }

    numFiles := len(imgList)

    fileMap := make(map[int]map[string]map[string][]os.FileInfo)

    for _, _f := range imgList {
        t := _f.ModTime()
        y := t.Year()
        m := t.Month().String()
        d := strconv.Itoa(t.Day()) +  " " + t.Weekday().String()

        if _, ok := fileMap[y]; !ok {
            fileMap[y] = make(map[string]map[string][]os.FileInfo)
        }

        if _, ok := fileMap[y][m]; !ok {
            fileMap[y][m] = make(map[string][]os.FileInfo)
        }

        if val, ok := fileMap[y][m][d]; ok {
            fileMap[y][m][d] = append(val, _f)
        } else {
            fileMap[y][m][d] = []os.FileInfo{_f}
        }
    }


    // now create the sorted dir
    os.Mkdir(cwd + "/sorted", os.FileMode(0777))

    copied := 0
    chnl := make(chan int, 1)

    for Y, monthMap := range fileMap {
        var path = cwd + "/sorted/" + strconv.Itoa(Y)

        // fmt.Println("Sorting for " + strconv.Itoa(Y));
        os.Mkdir(path, os.FileMode(0777))

        for M, dayMap := range monthMap {
            // fmt.Println("Sorting for " + M  + " " + strconv.Itoa(Y));
            path = cwd + "/sorted/" + strconv.Itoa(Y) + "/" + M

            os.Mkdir(path, os.FileMode(0777))

            for D, files := range dayMap {
                // fmt.Print(D);
                path = cwd + "/sorted/" + strconv.Itoa(Y) + "/" + M + "/" + D

                os.Mkdir(path, os.FileMode(0777))

                for _, _f := range files {
                    go func(src string, dst string, date time.Time) {
                        copyFileContents(src, dst, date)
                        defer func() {chnl <- 1}()
                    }(*srcDir + "/" + _f.Name(), path + "/" + _f.Name(), _f.ModTime())
                }
            }
        }
    }

    fmt.Printf("\nCopying files\t %d/%d (%d%%)", 0, numFiles, int(0.0/numFiles * 100))

    for {
        select {
            case <-chnl:
                copied++
                updateStatus(numFiles, copied)
            default:
        }

        if copied == numFiles {
            break;
        }
    }
    
    fmt.Println("")

}


func updateStatus(total int, num int) {
    fmt.Printf("\r Copying files\t %d/%d (%d%%)", num, total, int(float64(num)/float64(total) * 100))
}

// src http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src string, dst string, date time.Time) (err error) {
    in, err := os.Open(src)
    if err != nil { return err }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil { return err }
    defer out.Close()
    _, err = io.Copy(out, in)
    cerr := out.Close()
    if err != nil { return err }

    cerr = os.Chtimes(dst, date, date)

    return cerr
}