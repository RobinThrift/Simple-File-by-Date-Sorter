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

    for Y, monthMap := range fileMap {
        var path = cwd + "/sorted/" + strconv.Itoa(Y)

        // fmt.Println("Sorting for " + strconv.Itoa(Y));
        os.Mkdir(path, os.FileMode(0777))

        for M, dayMap := range monthMap {
            fmt.Println("Sorting for " + M  + " " + strconv.Itoa(Y));
            path = cwd + "/sorted/" + strconv.Itoa(Y) + "/" + M

            os.Mkdir(path, os.FileMode(0777))

            for D, files := range dayMap {
                // fmt.Print(D);
                path += "/" + D

                os.Mkdir(path, os.FileMode(0777))

                for _, _f := range files {
                    copyFileContents(*srcDir + "/" + _f.Name(), path + "/" + _f.Name(), _f.ModTime())
                }
            }
        }
    }

    var input string
    fmt.Scanln(&input)

}


// src http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src string, dst string, date time.Time) (err error) {
    // fmt.Println(src);
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