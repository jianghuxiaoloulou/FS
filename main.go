package main

import (
	"WowjoyProject/FileServer/models"
	"WowjoyProject/FileServer/pkg/loggin"
	"WowjoyProject/FileServer/pkg/object"
	"WowjoyProject/FileServer/pkg/setting"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DataChan chan models.ObjectData
)

//任务
type Job interface {
	// do something...
	Do()
}

//worker 工人
type Worker struct {
	JobQueue chan Job  //任务队列
	Quit     chan bool //停止当前任务
}

//新建一个 worker 通道实例  新建一个工人
func NewWorker() Worker {
	return Worker{
		JobQueue: make(chan Job), //初始化工作队列为null
		Quit:     make(chan bool),
	}
}

/*
整个过程中 每个Worker(工人)都会被运行在一个协程中，
在整个WorkerPool(领导)中就会有num个可空闲的Worker(工人)，
当来一条数据的时候，领导就会小组中取一个空闲的Worker(工人)去执行该Job，
当工作池中没有可用的worker(工人)时，就会阻塞等待一个空闲的worker(工人)。
每读到一个通道参数 运行一个 worker
*/

func (w Worker) Run(wq chan chan Job) {
	//这是一个独立的协程 循环读取通道内的数据，
	//保证 每读到一个通道参数就 去做这件事，没读到就阻塞
	go func() {
		for {
			wq <- w.JobQueue //注册工作通道  到 线程池
			select {
			case job := <-w.JobQueue: //读到参数
				job.Do()
			case <-w.Quit: //终止当前任务
				return
			}
		}
	}()
}

//workerpool 领导
type WorkerPool struct {
	workerlen   int      //线程池中  worker(工人) 的数量
	JobQueue    chan Job //线程池的  job 通道
	WorkerQueue chan chan Job
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,                      //开始建立 workerlen 个worker(工人)协程
		JobQueue:    make(chan Job),                 //工作队列 通道
		WorkerQueue: make(chan chan Job, workerlen), //最大通道参数设为 最大协程数 workerlen 工人的数量最大值
	}
}

//运行线程池
func (wp *WorkerPool) Run() {
	//初始化时会按照传入的num，启动num个后台协程，然后循环读取Job通道里面的数据，
	//读到一个数据时，再获取一个可用的Worker，并将Job对象传递到该Worker的chan通道
	loggin.Debug("初始化worker")
	for i := 0; i < wp.workerlen; i++ {
		//新建 workerlen worker(工人) 协程(并发执行)，每个协程可处理一个请求
		worker := NewWorker() //运行一个协程 将线程池 通道的参数  传递到 worker协程的通道中 进而处理这个请求
		worker.Run(wp.WorkerQueue)
	}

	// 循环获取可用的worker,往worker中写job
	go func() { //这是一个单独的协程 只负责保证 不断获取可用的worker
		for {
			select {
			case job := <-wp.JobQueue: //读取任务
				//尝试获取一个可用的worker作业通道。
				//这将阻塞，直到一个worker空闲
				worker := <-wp.WorkerQueue
				worker <- job //将任务 分配给该工人
			}
		}
	}()
}

//----------------------------------------------
type Dosomething struct {
	key models.ObjectData
}

func (d *Dosomething) Do() {
	loggin.Info("正在处理的数据是：", d.key)
	obj := object.NewObject(d.key)
	switch d.key.Type {
	case models.UPLOAD:
		// 数据上传
		obj.UploadObject()
	case models.DOWNLOAD:
		// 数据下载
		obj.DownObject()
	case models.DELETE:
		// 数据删除
		obj.DelObject()
	}
}

func main() {
	// 数据通道
	DataChan = make(chan models.ObjectData)
	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)设置最大线程数
	wokerPool := NewWorkerPool(setting.MaxThreads)
	wokerPool.Run() //有任务就去做，没有就阻塞，任务做不过来也阻塞	// 初始化数据库
	// 处理任务：
	go func() { //这是一个独立的协程 保证可以接受到每个用户的请求
		for {
			select {
			case data := <-DataChan:
				sc := &Dosomething{key: data}
				wokerPool.JobQueue <- sc //往线程池 的通道中 写参数   每个参数相当于一个请求  来了100万个请求
			}
		}
	}()
	// router := routers.InitRouter()

	// s := &http.Server{
	// 	Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
	// 	Handler:        router,
	// 	ReadTimeout:    setting.ReadTimeout,
	// 	WriteTimeout:   setting.WriteTimeout,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// // http 服务
	// go func() {
	// 	if err := s.ListenAndServe(); err != nil {
	// 		loggin.Info("Listen:%s\n", err)
	// 	}
	// }()
	// 主程序逻辑：自动上传文件
	for {
		// models.AutoUploadData(DataChan)
		// 等待10秒下次查询数据
		// time.Sleep(time.Second * 10)
		// 自动下载任务：
		models.AutoDownData(DataChan)
		// models.TestAutoDownData(DataChan)

		// time.Sleep(60 * time.Second)

	}
}

//因为要多次检查错误，所以建立一个函数。
func check(err error) {
	if err != nil {
		loggin.Error(err)
	}
}

func CheckPath(path string) {
	dir, _ := filepath.Split(path)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
}
