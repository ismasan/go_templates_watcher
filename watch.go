package main
import(
  "github.com/howeyc/fsnotify"
  "path/filepath"
  "io/ioutil"
  "strings"
  "net/http"
  "log"
)

const(
  API_HOST = "http://localhost:4567"
)

var HTTP = &http.Client{}

// Helpers
func checkError(err error) {
  if err != nil {
    log.Fatal("[error]", err)
  }
}

// HTTP STUFF
func put(path, contentType string, fileData []byte) error {
  
  req, err := http.NewRequest("PUT", API_HOST + path, nil)
  // ...
  req.Header.Add("Content-type", contentType)
  resp, err := HTTP.Do(req)
  
  log.Println("RESP", resp)
  return err
}

// FILE STUFF
func getFileName(filePath string) string {
  segments := strings.Split(filePath, string(filepath.Separator))
  l := len(segments)
  return segments[l - 1]
}

func Create (fileName string) {
  log.Println("[create]", fileName)
}

func Delete (fileName string) {
  log.Println("[delete]", fileName)
}

func Modify (filePath string) {
  log.Println("[modify]", filePath)
  fileData, err := ioutil.ReadFile(filePath)
  checkError(err)
  fileName := getFileName(filePath)
  err = put("/templates/" + fileName, "application/json", fileData)
}

func Rename (fileName string) {
  log.Println("[rename]", fileName)
}

func main() {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
      log.Fatal(err)
  }
  
  absPath, err := filepath.Abs("./templates")
  
  if err != nil {
    log.Fatal("Error finding ./templates dir")
  }
  
  err = watcher.Watch(absPath)
  if err != nil {
      log.Fatal(err)
  }
  log.Println("WATCHING")
  
  for {
      select {
      case ev := <-watcher.Event:
          if ev.IsCreate() {
            go Create(ev.Name)
          } else if ev.IsDelete() {
            go Delete(ev.Name)
          } else if ev.IsModify() {
            go Modify(ev.Name)
          } else if ev.IsRename() {
            go Rename(ev.Name)
          } else {
            panic("Unknown operation." + ev.Name)
          }
      case err := <-watcher.Error:
          log.Println("error:", err)
      }
  }
  panic("NEVER HERE")
      
  /* ... do stuff ... */
  defer func() {
    log.Println("BYEBYE")
    watcher.Close()  
  }()
}
