package sigr

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

var onStopService = struct {
	m       sync.Mutex
	handers map[string]func()
	state   int32
}{
	m:       sync.Mutex{},
	handers: map[string]func(){},
	state:   0,
}

var autoIncId uint64 = 0
var verboseLog = false
var quitDirectly = true

func SetQuitDirectyl(setting bool) {
	quitDirectly = setting
}

func VerboseLog() {
	verboseLog = true
}

func NoLog() {
	verboseLog = false
}

func handlerNameExists(name string) bool {
	onStopService.m.Lock()
	defer onStopService.m.Unlock()
	_, ok := onStopService.handers[name]
	return ok
}

func RegisterOnStopFuncAutoName(f func()) (name string) {
	atomic.AddUint64(&autoIncId, 1)
	name = fmt.Sprintf("__auto_%d", autoIncId)
	for handlerNameExists(name) {
		atomic.AddUint64(&autoIncId, 1)
		name = fmt.Sprintf("__auto_%d", autoIncId)
	}
	RegisterOnStopFunc(name, f)
	return
}

func RegisterOnStopFunc(name string, f func()) {
	// register a new function on signal int(interrupt) and term(terminate)
	onStopService.m.Lock()
	defer onStopService.m.Unlock()
	if atomic.CompareAndSwapInt32(&onStopService.state, 0, 1) {
		if verboseLog {
			log.Println("OnStopService Initialized..")
		}
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go func() {
			var sig os.Signal = <-sigs
			if verboseLog {
				log.Println(fmt.Sprintf("sig: [%v] Processing...", sig))
			}
			for k, v := range onStopService.handers {
				if verboseLog {
					log.Println(fmt.Sprintf("Processing task [%v]", k))
				}
				v()
			}
			signal.Stop(sigs)
			if verboseLog {
				if quitDirectly {
					log.Println(fmt.Sprintf("sig: [%v] Processed, quitting directly", sig))
				} else {
					log.Println(fmt.Sprintf("sig: [%v] Processed, please try again to terminate the process", sig))
				}
			}
			atomic.CompareAndSwapInt32(&onStopService.state, 1, 0)
			if quitDirectly {
				if p, e := os.FindProcess(syscall.Getpid()); e == nil {
					p.Signal(sig)
				}
			}
		}()
	}
	onStopService.handers[name] = f
}

func UnregisterOnStopFunc(name string) {
	onStopService.m.Lock()
	defer onStopService.m.Unlock()
	delete(onStopService.handers, name)
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
