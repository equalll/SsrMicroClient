package main

import (
    "fmt"
    "encoding/base64"
    "net/http"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
    "os/exec"
    "bytes"
    "regexp"
//    "time"
    "runtime"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
  //  "sync"
//    "log"
    //"strconv"

    "./subscription"
    "./init"
    "./node"
)

var ssr_config_path string
type ssr_config struct {
    python_path,config_path,log_file,pid_file,fast_open,workers string
    connect_verbose_info,ssr_path,server,server_port,protocol,method string
    obfs,password,obfsparam,protoparam,local_port,local_address,remarks,config_url,deamon string
    acl string
}


//判断目录是否存在返回布尔类型
func path_exists(path string)bool{
    _,err := os.Stat(path)
    if err!=nil{
        if os.IsExist(err){
            return true
        }else{
            return false
        }
    }else{
        return true
    }
}

//对base64进行长度补全(4的倍数)
func base64d(str string)string{
    for i:=0;i<=len(str)%4;i++{
        str+="="
    }
    de_str,_ := base64.URLEncoding.DecodeString(str)
    return string(de_str)
}

func ssr_config_init()ssr_config{
    if runtime.GOOS=="linux"{
        ssr_config_path = os.Getenv("HOME")+"/.config/SSRSub/ssr_config.conf"
    }
    
    return ssr_config{"","","","","","","","","","","","","","","","","","","","","",""}
}

//读取配置文件
func read_config()ssr_config{
    ssr_config := ssr_config_init()
    //var log_file,pid_file,fast_open,workers,connect_verbose_info,ssr_path,python_path,config_path,config_url string
    config_temp,err := ioutil.ReadFile(ssr_config_path)
    if err != nil {
        fmt.Println(err)
    }
    lines := strings.Split(string(config_temp),"\n")
    re3,_ := regexp.Compile("#.*$")
    for _,config_temp2 := range lines{
        config_temp2 = re3.ReplaceAllString(config_temp2,"")
        config_temp2 := strings.Split(config_temp2," ")
        if config_temp2[0] == "python_path"{
            ssr_config.python_path = config_temp2[1]+" "
        } else if config_temp2[0] == "ssr_path"{
            ssr_config.ssr_path = config_temp2[1]+" "
        }else if config_temp2[0] == "config_path"{
            ssr_config.config_path = config_temp2[1]
        }else if config_temp2[0] == "config_url"{
            ssr_config.config_url = config_temp2[1]
        }else if config_temp2[0] == "ssr_config"{
            ssr_server_config := strings.Split(config_temp2[1],",")
            ssr_config.server = "-s "+ssr_server_config[0]+" "
            ssr_config.server_port = "-p " +ssr_server_config[1]+" "
            ssr_config.protocol = "-O "+ssr_server_config[2]+" "
            ssr_config.method = "-m "+ssr_server_config[3]+" "
            ssr_config.obfs = "-o "+ssr_server_config[4]+" "
            ssr_config.password = "-k "+ssr_server_config[5]+" "
            ssr_config.obfsparam = "-g "+ssr_server_config[6]+" "
            ssr_config.protoparam = "-G "+ssr_server_config[7]+" "
            ssr_config.remarks = ssr_server_config[8]+" "
        }else if config_temp2[0] == "connect-verbose-info"{
            ssr_config.connect_verbose_info = "--connect-verbose-info "
        }else if config_temp2[0] == "workers"{
            ssr_config.workers = "--workers "+config_temp2[1]+" "
        }else if config_temp2[0] == "fast-open"{
            ssr_config.fast_open = "--fast-open "
        }else if config_temp2[0] == "pid-file"{
            ssr_config.pid_file = "--pid-file "+config_temp2[1]+" "
        }else if config_temp2[0] == "log-file"{
            ssr_config.log_file = "--log-file "+config_temp2[1]+" "
        }else if config_temp2[0] == "local_address"{
            ssr_config.local_address = "-b "+config_temp2[1]+" "
        }else if config_temp2[0] == "local_port"{
            ssr_config.local_port = "-l "+config_temp2[1]+" "
        }else if config_temp2[0] == "acl"{
            ssr_config.acl = "--acl "+config_temp2[1]+" "
        }else if config_temp2[0] == "deamon"{
            ssr_config.deamon = "-d start"
        }
    }
    //fmt.Println(python_path+ssr_path+log_file+pid_file+fast_open+workers+connect_verbose_info+config_url+config_path)
    return ssr_config
    //for _,config_temp2 := range config_temp{
    //    fmt.Println(string(config_temp2))
    //}
}

//读取配置文件
func read_config_db(db_path string)ssr_config{
    ssr_config := ssr_config_init()
    //var log_file,pid_file,fast_open,workers,connect_verbose_info,ssr_path,python_path,config_path,config_url string
    config_temp,err := ioutil.ReadFile(ssr_config_path)
    if err != nil {
        fmt.Println(err)
    }
    lines := strings.Split(string(config_temp),"\n")
    re3,_ := regexp.Compile("#.*$")
    for _,config_temp2 := range lines{
        config_temp2 = re3.ReplaceAllString(config_temp2,"")
        config_temp2 := strings.Split(config_temp2," ")
        if config_temp2[0] == "python_path"{
            ssr_config.python_path = config_temp2[1]+" "
        } else if config_temp2[0] == "ssr_path"{
            ssr_config.ssr_path = config_temp2[1]+" "
        }else if config_temp2[0] == "config_path"{
            ssr_config.config_path = config_temp2[1]
        }else if config_temp2[0] == "connect-verbose-info"{
            ssr_config.connect_verbose_info = "--connect-verbose-info "
        }else if config_temp2[0] == "workers"{
            ssr_config.workers = "--workers "+config_temp2[1]+" "
        }else if config_temp2[0] == "fast-open"{
            ssr_config.fast_open = "--fast-open "
        }else if config_temp2[0] == "pid-file"{
            ssr_config.pid_file = "--pid-file "+config_temp2[1]+" "
        }else if config_temp2[0] == "log-file"{
            ssr_config.log_file = "--log-file "+config_temp2[1]+" "
        }else if config_temp2[0] == "local_address"{
            ssr_config.local_address = "-b "+config_temp2[1]+" "
        }else if config_temp2[0] == "local_port"{
            ssr_config.local_port = "-l "+config_temp2[1]+" "
        }else if config_temp2[0] == "acl"{
            ssr_config.acl = "--acl "+config_temp2[1]+" "
        }else if config_temp2[0] == "deamon"{
            ssr_config.deamon = "-d start"
        }
    }

    db,err := sql.Open("sqlite3",db_path)
    if err!=nil{
        fmt.Println(err)
    }
    defer db.Close()
    var server,server_port,protocol,method,obfs,password,obfsparam,protoparam string
    rows,err := db.Query("SELECT server,server_port,protocol,method,obfs,password,obfsparam,protoparam FROM SSR_present_node")
    for rows.Next(){rows.Scan(&server,&server_port,&protocol,&method,&obfs,&password,&obfsparam,&protoparam)}
    ssr_config.server = "-s "+server+" "
    ssr_config.server_port = "-p " +server_port+" "
    ssr_config.protocol = "-O "+protocol+" "
    ssr_config.method = "-m "+method+" "
    ssr_config.obfs = "-o "+obfs+" "
    ssr_config.password = "-k "+password+" "
    ssr_config.obfsparam = "-g "+obfsparam+" "
    ssr_config.protoparam = "-G "+protoparam+" "
    fmt.Println(ssr_config)
    return ssr_config
}

//读取订阅链接文件(后面改成数据库)
func read_ssr_config()string{
    config_temp,err := ioutil.ReadFile(read_config().config_path)
    if err != nil {
        fmt.Println(err)
    }
    return string(config_temp)
}

//更新订阅
func update_config(){
    res,_ := http.Get(read_config().config_url)
    body,err := ioutil.ReadAll(res.Body)
    if err!=nil{
        fmt.Println(err)
        fmt.Println("可能出错原因,请检查能否成功访问订阅连接.")
        return
    }
    ioutil.WriteFile(read_config().config_path,[]byte(body),0644)
}



//方便进行分割对字符串进行替换
func str_replace(str string)[]string{
    var config[] string
    scanner := bufio.NewScanner(strings.NewReader(strings.Replace(base64d(str),"ssr://","",-1)))
    for scanner.Scan() {
    str_temp := strings.Replace(base64d(scanner.Text()),"/?obfsparam=",":",-1)
    str_temp = strings.Replace(str_temp,"&protoparam=",":",-1)
    str_temp = strings.Replace(str_temp,"&remarks=",":",-1)
    str_temp = strings.Replace(str_temp,"&group=",":",-1)
    config = append(config,str_temp)
    }
    return config
}

func list_list(config_array []string){
    for num,config_temp := range config_array{
        config_temp2 := strings.Split(config_temp,":")
        fmt.Println(num+1,base64d(config_temp2[len(config_temp2)-2]))
    }
}


//更换节点
func ssr__server_config(){
    config_middle_temp := str_replace(string(read_ssr_config()))
    list_list(config_middle_temp)
    var config_split []string
    select_temp := menu_select()-1
    if (select_temp>0&&select_temp<len(config_middle_temp)){
        config_split = strings.Split(config_middle_temp[select_temp],":")
    }else{
        fmt.Println("\nenter error,please enter correct number.")
        ssr__server_config()
        return
    }
    var server string
    if len(config_split) == 17 {
        server = config_split[0]+":"+config_split[1]+":"+config_split[2]+":"+config_split[3]+":"+config_split[4]+":"+config_split[5]+":"+config_split[6]+":"+config_split[7]
    } else if len(config_split) == 10 {
        server = config_split[0]
    }
    server_port := config_split[len(config_split)-9]
    protocol := config_split[len(config_split)-8]
    method := config_split[len(config_split)-7]
    obfs := config_split[len(config_split)-6]
    password := base64d(config_split[len(config_split)-5])
    obfsparam := base64d(config_split[len(config_split)-4])
    protoparam := base64d(config_split[len(config_split)-3])
    remarks := base64d(config_split[len(config_split)-2])
    //return server,server_port,protocol,method,obfs,password,obfsparam,protoparam,remarks
    //return ssr_config{server:server,server_port:server_port,protocol:protocol,method:method,obfs:obfs,password:password,obfsparam:obfsparam,protoparam:protoparam,remarks:remarks}
    config_temp,err := ioutil.ReadFile(ssr_config_path)
    if err != nil {
        fmt.Println(err)
    }
    lines := strings.Split(string(config_temp),"\n")
    //scanner := bufio.NewScanner(strings.NewReader(strings.Replace(string(config_temp)," ","",-1)))
    ///scanner := bufio.NewScanner(strings.NewReader(string(config_temp)))
    //for scanner.Scan(){
    for num,line := range lines{
        if strings.Contains(line, "ssr_config"){
            lines[num] = "ssr_config "+server+","+server_port+","+protocol+","+method+","+obfs+","+password+","+obfsparam+","+protoparam+","+strings.Replace(remarks," ","",-1)
        }
        output := strings.Join(lines, "\n")
        ioutil.WriteFile(ssr_config_path,[]byte(output),0644)
    }
}


func ssr_start(){
    ssr_config := read_config()
    /*fmt.Println(ssr_config.python_path,ssr_config.ssr_path,ssr_config.
        log_file,ssr_config.pid_file,ssr_config.fast_open,ssr_config.
        workers,ssr_config.connect_verbose_info,ssr_config.server,ssr_config.
        server_port,ssr_config.protocol,ssr_config.method,ssr_config.
        obfs,ssr_config.password,ssr_config.obfsparam,ssr_config.protoparam)*/

    /*cmd := exec.Command(ssr_config.python_path,ssr_config.ssr_path,ssr_config.fast_open,ssr_config.connect_verbose_info,strings.
        Split(ssr_config.local_port," ")[0],strings.Split(ssr_config.local_port," ")[1],strings.
        Split(ssr_config.log_file," ")[0],strings.Split(ssr_config.log_file," ")[1],strings.
        Split(ssr_config.pid_file," ")[0],strings.Split(ssr_config.pid_file," ")[1],strings.
        Split(ssr_config.workers," ")[0],strings.Split(ssr_config.workers," ")[1],strings.
        Split(ssr_config.server," ")[0],strings.Split(ssr_config.server," ")[1],strings.
        Split(ssr_config.server_port," ")[0],strings.Split(ssr_config.server_port," ")[1],strings.
        Split(ssr_config.protocol," ")[0],strings.Split(ssr_config.protocol," ")[1],strings.
        Split(ssr_config.method," ")[0],strings.Split(ssr_config.method," ")[1],strings.
        Split(ssr_config.obfs," ")[0],strings.Split(ssr_config.obfs," ")[1],strings.
        Split(ssr_config.password," ")[0],strings.Split(ssr_config.password," ")[1],strings.
        Split(ssr_config.obfsparam," ")[0],strings.Split(ssr_config.obfsparam," ")[1],strings.
        Split(ssr_config.protoparam," ")[0],strings.Split(ssr_config.protoparam," ")[1],strings.
        Split(ssr_config.deamon," ")[0],strings.Split(ssr_config.deamon," ")[1])
    */

    cmd_temp := ssr_config.python_path+ssr_config.ssr_path+ssr_config.local_address+ssr_config.
    local_port+ssr_config.log_file+ssr_config.pid_file+ssr_config.fast_open+ssr_config.
    workers+ssr_config.connect_verbose_info+ssr_config.server+ssr_config.
    server_port+ssr_config.protocol+ssr_config.method+ssr_config.
    obfs+ssr_config.password+ssr_config.obfsparam+ssr_config.protoparam+ssr_config.acl+ssr_config.deamon

    fmt.Println(cmd_temp)

    var cmd *exec.Cmd
    if runtime.GOOS == "linux"{
        cmd = exec.Command("/bin/sh", "-c",cmd_temp)
    }
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
        return
    }
    fmt.Println("Result: " + out.String())
    //fmt.Println(ssr_config.python_path,ssr_config.config_path,ssr_config.log_file,ssr_config.pid_file,ssr_config.fast_open,ssr_config.workers,ssr_config.connect_verbose_info,ssr_config.ssr_path,ssr_config.server,ssr_config.server_port,ssr_config.protocol,ssr_config.method,ssr_config.obfs,ssr_config.password,ssr_config.obfsparam,ssr_config.protoparam) 
}

func ssr_start_db(db_path string){
    ssr_config := read_config_db(db_path)
    cmd_temp := ssr_config.python_path+ssr_config.ssr_path+ssr_config.local_address+ssr_config.
    local_port+ssr_config.log_file+ssr_config.pid_file+ssr_config.fast_open+ssr_config.
    workers+ssr_config.connect_verbose_info+ssr_config.server+ssr_config.
    server_port+ssr_config.protocol+ssr_config.method+ssr_config.
    obfs+ssr_config.password+ssr_config.obfsparam+ssr_config.protoparam+ssr_config.acl+ssr_config.deamon

    fmt.Println(cmd_temp)

    var cmd *exec.Cmd
    if runtime.GOOS == "linux"{
        cmd = exec.Command("/bin/sh", "-c",cmd_temp)
    }
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
        return
    }
    fmt.Println("Result: " + out.String())
    //fmt.Println(ssr_config.python_path,ssr_config.config_path,ssr_config.log_file,ssr_config.pid_file,ssr_config.fast_open,ssr_config.workers,ssr_config.connect_verbose_info,ssr_config.ssr_path,ssr_config.server,ssr_config.server_port,ssr_config.protocol,ssr_config.method,ssr_config.obfs,ssr_config.password,ssr_config.obfsparam,ssr_config.protoparam) 
}


func ssr_stop(){
    cmd_temp := "cat "+strings.Split(read_config().pid_file," ")[1]+" | xargs kill"
    var cmd *exec.Cmd
    if runtime.GOOS == "linux"{
        cmd = exec.Command("/bin/sh", "-c",cmd_temp)
    }
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
        return
    }
    fmt.Println("Result: " + out.String())
}

func menu_select()int{
    var n int
    fmt.Scanln(&n)
    return n
}

func menu(){
    fmt.Println(runtime.GOOS+" "+runtime.GOARCH)
    fmt.Println("当前使用节点: "+read_config().remarks)
    fmt.Println("1.开启ssr\n2.update config\n3.更换节点\n4.get获取延迟\n5.结束ssr后台")
    select_temp := menu_select()

    if select_temp==1{
        ssr_start()
    }else if select_temp==3{
        ssr__server_config()
        menu()
        return
    }else if select_temp==2{
        update_config()
        menu()
        return
    }else if select_temp==4{
        menu()
        return
    }else if select_temp==5{
        ssr_stop()
        menu()
        return
    }else{
        fmt.Println("\nenter error,please enter correct number.")
        menu()
        return
    }
}


func menu_db(path,db_path string){
    //获取当前配置文件路径和可执行文件路径
    ssr_init.Menu_init(path)
    //获取当前节点
    node.Get_now_node(db_path)

    fmt.Print("1.开启ssr\n2.更换节点\n3.更新所有订阅\n4.添加订阅链接\n5.删除订阅链接\n6.获取延迟\n7.结束ssr后台\n>>>")



    var select_temp string
    fmt.Scanln(&select_temp)

    switch select_temp{
    case "1":
        ssr_start_db(db_path)
    case "2":
        node.Ssr_server_node_change(db_path)
    case "3":
        subscription.Add_config_db(db_path)
    case "4":
        fmt.Print("请输入要添加的订阅链接(一条):")
        var link_temp string
        fmt.Scanln(&link_temp)
        subscription.Subscription_link_add(link_temp,db_path)
    case "5":
        subscription.Subscription_link_delete(db_path)
    case "6":
    case "7":
    }

}

func test(db_path string){

    node.Ssr_server_node_change(db_path)
}
func main(){
    config_path := os.Getenv("HOME")+"/.config/SSRSub"
    path := os.Getenv("HOME")+"/.config/SSRSub/SSR_config.db"
    //ssr__server_config_db()
    //menu()
    //node.Ssr_server_node_change(path)
    menu_db(config_path,path)
    //subscription.Subscription_link_add("https://anydoor.me/link/zeGO2rxhKagZIqIc?mu=1",path)
    //fmt.Println(get_subscription_link())
    //subscription.Subscription_link_delete(path)


    //subscription.Add_config_db(path)


    //ssr_init.Init(config_path,path)
    /*
    node.List_list_db(path)
    subscription.Delete_config_db(path)
    fmt.Println(subscription.Get_subscription_link(path))
    */
    //ssr_start_db(path)
}