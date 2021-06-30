package main

import (
        "RAMS/src"
        "flag"
        "fmt"
        "log"
        "net/http"
        "os"
        "time"
)

//------//
var file_str string = ""
var app_key string = ""
var src_path string = ""
//------//

func now_time() string {
        date := time.Now()
        dat := date.Format(("02/Jan/2006 15:04:05"))
        return dat
}

func handler_file(w http.ResponseWriter, r *http.Request) {
        var nowlog string
        content,_ := src.OpenFile(file_str)

        if len(app_key) > 0 {
                var output_key string
                out,ok := r.URL.Query()["key"]
                if !ok || len(out) < 1 {
                        output_key = ""
                } else {
                        output_key = out[0]
                }
                if output_key != app_key {
                        var key_show string
                        if len(output_key) > 25 {
                                key_show = output_key
                        } else {
                                key_show = "key too long"
                        }
                        nowlog = now_time() +" - "+file_str+" - "+ r.RemoteAddr +"<- Wrong key: "+key_show
                        w.Write([]byte("Please, provide a valid key!"))
                } else {
                        w.Write(content)
                        nowlog = now_time() +" - "+file_str+" - "+ r.RemoteAddr
                }
        } else {
                nowlog = now_time() +" - "+file_str+" - "+ r.RemoteAddr
                w.Write(content)
        }

        fmt.Println(nowlog)
}

func files() http.Handler {
        return http.StripPrefix("/", http.FileServer(http.Dir(src_path)))
}

func main() {
        var file = flag.String("f","","File to delivery.")
        var port = flag.String("p","8000","Local Port")
        var host = flag.String("h","0.0.0.0","Local Host")
        var key = flag.String("k","","Key")
        var root = flag.String("root","","Root.")
        flag.Parse()

        app_key = *key
        hh := *host
        pp := *port
        h := hh+":"+pp

        println("SERVY")
        if len(app_key) > 0 {
                println("Key:",app_key)
        }

        println("Start Listener:",h)
        if len(*file) > 0 {
                file_str = *file
                println("Serving file:",file_str)
                http.HandleFunc("/", handler_file)
                log.Fatal(http.ListenAndServe(h, nil))
        } else {
                var path string
                if len(*root) > 0 {
                        path = *root
                        println("Source path:",path)
                } else {
                        path = os.Getenv("PWD")
                        println("Source path:",path)
                }
                src_path = path
                http.Handle("/",files())
                log.Fatal(http.ListenAndServe(h, nil))
        }

}
