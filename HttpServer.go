package gothic

import (
	"reflect"
	"runtime"
	"fmt"
	"net/http"
	"time"
	"strings"
)

const DEFAULT_SERVER_NAME = "GothicServer"
const DefaultListenPort = 9991

var gothicHttpApiServer = NewDefaultGothicHttpServer()

/**
 * @desc HttpServer
 * @author zhaojiangwei
 * @date 2018-12-17 18:50
 */
func NewGothicHttpServer(addr string, port int, readTimeout, writeTimeout, idleTimeout int, pprof bool) *GothicHttpApiServer {
	ret := &GothicHttpApiServer{
		HttpAddr:     addr,
		HttpPort:     port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		hanlder:      &httpApiHandler{routMap: make(map[string]map[string]reflect.Type), enablePprof: pprof},
	}
	return ret
}

func NewDefaultGothicHttpServer() *GothicHttpApiServer{
	ret := &GothicHttpApiServer{
		HttpAddr:     "",
		HttpPort:     DefaultListenPort,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
		hanlder:      &httpApiHandler{routMap: make(map[string]map[string]reflect.Type), enablePprof: false},
	}

	return ret
}

//http服务监听,路由
type GothicHttpApiServer struct {
	HttpAddr     string
	HttpPort     int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	hanlder      *httpApiHandler
}

func (this *GothicHttpApiServer) AddController(c interface{}) {
	this.hanlder.addController(c)
}

func (this *GothicHttpApiServer) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := fmt.Sprintf("%s:%d", this.HttpAddr, this.HttpPort)
	s := &http.Server{
		Addr:              addr,
		Handler:           this.hanlder,
		ReadHeaderTimeout: time.Duration(this.ReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(this.WriteTimeout) * time.Millisecond,
		IdleTimeout:       time.Duration(this.IdleTimeout) * time.Millisecond,
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timeNow, DEFAULT_SERVER_NAME + " Listen At: ", addr)
	fmt.Println("------------------------------------------------------------------------")
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

//controller中以此结尾的方法会参与路由
const METHOD_EXPORT_TAG = "Action"

type httpApiHandler struct {
	routMap     map[string]map[string]reflect.Type //key:controller: {key:method value:reflect.type}
	enablePprof bool
	Application *GothicApplication
}

func (this *httpApiHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			switch errInfo := err.(type) {
			case error:
				http.Error(rw, errInfo.Error(), http.StatusOK)
			default:
				Logger.Format(EntryFields{"msg": "parse request error", "error": err}).Error()
				http.Error(rw, fmt.Sprintln(err), http.StatusInternalServerError)
			}
		}
	}()

	Logger.Format(EntryFields{"msg": "New Connection In", "url": r.URL.Path}).Debug()

	serverHeader := Config.GetString("application.server_header")
	if serverHeader == ""{
		serverHeader = DEFAULT_SERVER_NAME
	}
	rw.Header().Set("Server", serverHeader)
	r.ParseMultipartForm(defaultMaxMultipartMemory)

	router := r.URL.Path

	//不处理favicon.ico路由
	if !Config.GetBool("application.enable_favicon") {
		if router == "/favicon.ico" {
			return
		}
	}

	routers := strings.Split(router, "/")
	if len(routers) < 3 || routers[1] == "" || routers[2] == "" {
		http.NotFound(rw, r)
		return
	}

	cname := strings.Title(routers[1])
	mname := strings.Title(routers[2]) + METHOD_EXPORT_TAG
	canhandler := false

	var contollerType reflect.Type
	if cname != "" && mname != "" {
		if methodMap, ok := this.routMap[cname]; ok {
			if contollerType, ok = methodMap[mname]; ok {
				canhandler = true
			}
		}
	}

	if !canhandler {
		http.NotFound(rw, r)
		return
	}

	threadContext := &ThreadContext{
		Application: Application,
		Controller: cname,
		Action: mname,
		Params: make(map[string]interface{}),
	}
	if err := InvokeThreadHook(cname, mname, BeforeInitController, threadContext); err != nil{
		Logger.Format(EntryFields{"msg": "Hook Exception", "point": BeforeInitController, "url": router, "error": fmt.Sprint(err)}).Error()
		panic(err)
	}

	vc := reflect.New(contollerType)
	var in []reflect.Value
	var method reflect.Value

	defer func() {
		if err := recover(); err != nil {
			switch errInfo := err.(type) {
			case error:
				http.Error(rw, errInfo.Error(), http.StatusOK)
			default:
				Logger.Format(EntryFields{"msg": "Handle Request Exception", "url": router, "error": fmt.Sprint(err)}).Error()
				http.Error(rw, fmt.Sprintln("Internal Server Error"), http.StatusInternalServerError)
			}
		}
	}()

	in = make([]reflect.Value, 3)
	in[0] = reflect.ValueOf(threadContext)
	in[1] = reflect.ValueOf(rw)
	in[2] = reflect.ValueOf(r)
	method = vc.MethodByName("Init")
	method.Call(in)

	Logger.Format(EntryFields{"msg": "Action Begin", "controller": cname, "action": mname}).Debug()

	in = make([]reflect.Value, 0)
	method = vc.MethodByName(mname)
	if err := InvokeThreadHook(cname, mname, BeforeInvokeAction, threadContext); err != nil{
		Logger.Format(EntryFields{"msg": "Hook Exception", "point": BeforeInvokeAction, "url": router, "error": fmt.Sprint(err)}).Error()
		panic(err)
	}
	method.Call(in)
	if err := InvokeThreadHook(cname, mname, AfterInvokeAction, threadContext); err != nil{
		Logger.Format(EntryFields{"msg": "Hook Exception", "point": AfterInvokeAction, "url": router, "error": fmt.Sprint(err)}).Error()
		panic(err)
	}

	//post request
	method = vc.MethodByName("Destruct")
	method.Call(in)
	if err := InvokeThreadHook(cname, mname, AfterSendResponse, threadContext); err != nil{
		Logger.Format(EntryFields{"msg": "Hook Exception", "point": AfterSendResponse, "url": router, "error": fmt.Sprint(err)}).Error()
	}
}

func (this *httpApiHandler) addController(c interface{}) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	firstParam := strings.TrimSuffix(ct.Name(), "Controller")
	if _, ok := this.routMap[firstParam]; ok {
		return
	} else {
		this.routMap[firstParam] = make(map[string]reflect.Type)
	}
	var mname string
	for i := 0; i < rt.NumMethod(); i++ {
		mname = rt.Method(i).Name
		if strings.HasSuffix(mname, METHOD_EXPORT_TAG) {
			this.routMap[firstParam][rt.Method(i).Name] = ct
		}
	}
}