package logger

import (
	"database/sql"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/go-sql-driver/mysql"
)

/*
	提供一个向mysql写入日志的功能，写入方式和表格 需要进一步扩展
	可以被启用为一个基础服务，也可以作为独立服务
	例子:
		e := logger.NewRemoteServer(10, "")
		if e != nil {
			seelog.Error(e)
		}
		logger.RemoteServerWorkerNum(10)
		logger.RemoteServerRun()
	TODO：
		1.更多的数据存储源，或许可以对接于日志分析工具
		2.收集日志的格式提供定制化
		3.存储方式根据定制化做区分
		4.配套的查询组件
*/
var remoteLog *RemoteLogServer

// RemoteLogPush 日志服务器，直接写入到日志服务的方法
func RemoteLogPush(uid uint64, account string, kind, source int, param string) {
	l := &LogData{
		Uid:    uid,
		OpenId: account,
		Stamp:  time.Now().Unix(),
		Kind:   kind,
		Source: source,
		Param:  param,
	}
	remoteLog.AddLog(l)
}

// RemoteLogAddLog 外部接口发送的日志数据方法
func RemoteLogAddLog(l *LogData) {
	remoteLog.AddLog(l)
}

// RemoteServerClose 关闭
func RemoteServerClose() {
	remoteLog.Close()
}

// RemoteServerRun  启动 独立协程处理
func RemoteServerRun() {
	for i := 0; i < remoteLog.workerNum; i++ {
		remoteLog.addWorker()
	}
	remoteLog.RotateTable()
	go remoteLog.ConsumeLog()
}

// RemoteServerWorkerNum 设置工作线程熟练
func RemoteServerWorkerNum(i int) {
	remoteLog.workerNum = i
}

// NewRemoteServer 初始，设定工作协程数量，日志存储地址，
func NewRemoteServer(host string) error {
	db, e := sql.Open("mysql", host)
	if e != nil {
		seelog.Error("Init Mysql Error ", e)
		return e
	}
	e = db.Ping()
	if e != nil {
		seelog.Error("Mysql Init Error ", e)
		return e
	}
	workerPool := &RemoteLogServer{
		db:          db,
		mysqlDSN:    host,
		listChannel: make(chan *LogData, 200),
		workerNum:   1,
	}

	remoteLog = workerPool
	return nil
}

type RemoteLogServer struct {
	mysqlDSN    string
	db          *sql.DB
	workerPool  []*RemoteLogWorker
	nowTable    string
	listChannel chan *LogData
	workerNum   int
}

func (server *RemoteLogServer) AddLog(l *LogData) {
	if len(server.listChannel) > 199 {
		return
	}
	server.listChannel <- l
}

//独立协程 处理所有的日志写入
func (server *RemoteLogServer) ConsumeLog() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case data := <-server.listChannel:
			w := server.getWorker()
			w.Push(data)
		case <-ticker.C:
			server.RotateTable()
			server.Recycle()
		}
	}
}

func (server *RemoteLogServer) RotateTable() {
	tableName := "log_" + time.Now().Format("20060102")
	if tableName != server.nowTable {
		c, _ := mysql.ParseDSN(server.mysqlDSN)
		// check table exists and create new table
		row, e := server.db.Query("SELECT table_name FROM information_schema.TABLES WHERE table_name ='" + tableName + "' and table_schema ='" + c.DBName + "';")
		if e != nil {
			seelog.Error("Query Table Name Error ", e)
			return
		}
		var t string
		for row.Next() {
			if e = row.Scan(&t); e != nil {
				seelog.Error("Scan Error ", e)
			}
		}
		row.Close()
		if t == "" {
			server.createTable(tableName)
		}
		server.nowTable = tableName
		server.pushToWorker(1)
	}
}

func (server *RemoteLogServer) createTable(tableName string) {
	tableSql := "CREATE TABLE `" + tableName + "`  (" +
		"`id` int(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id'," +
		"`uid` int(20) UNSIGNED NOT NULL COMMENT '玩家id'," +
		"`openId` varchar(255) NOT NULL COMMENT '玩家账号'," +
		"`stamp` bigint(255) NOT NULL COMMENT '日志时间'," +
		"`kind` int(10) NOT NULL COMMENT '日志类型'," +
		"`source` int(10) NOT NULL COMMENT '日志来源'," +
		"`param` text COMMENT '日志参数'," +
		"PRIMARY KEY (`id`)," +
		"INDEX `uid`(`kind`, `source`,`stamp`, `uid`) USING BTREE," +
		"INDEX `kind`(`kind`,`source`) USING BTREE" +
		") ENGINE = MyISAM CHARACTER SET = utf8;"
	_, e := server.db.Exec(tableSql)
	if e != nil {
		seelog.Error(e)
	}
}

func (server *RemoteLogServer) Close() {
	ww := &sync.WaitGroup{}
	for _, w := range server.workerPool {
		ww.Add(1)
		w.Close(ww)
	}
	ww.Wait()
}
func (server *RemoteLogServer) Recycle() {

}
func (server *RemoteLogServer) pushToWorker(signal int) {
	for _, w := range server.workerPool {
		w.signal <- signal
	}
}
func (server *RemoteLogServer) addWorker() *RemoteLogWorker {
	lw := server.newWorker()
	go lw.do()
	server.workerPool = append(server.workerPool, lw)
	return lw
}
func (server *RemoteLogServer) getWorker() *RemoteLogWorker {
	var lw *RemoteLogWorker
	var n = 5000
	for _, w := range server.workerPool {
		if w.close {
			continue
		}
		c := len(w.data)
		if c >= 5000 {
			continue
		}
		if c <= n {
			n = c
			lw = w
		}
	}
	if lw == nil {
		lw = server.addWorker()
	}
	return lw
}
func (server *RemoteLogServer) newWorker() *RemoteLogWorker {
	//TODO 是否需要设置5000的缓存
	//TODO 另一种思路：worker不设置缓存，同时只能处理一个任务，有mgr做分配，
	//TODO  可以维护多个worker 分隔数据和逻辑 减少冗余内存
	//TODO  根据数据的量级进行worker数量的动态伸缩
	return &RemoteLogWorker{
		data:   make(chan *LogData, 5000),
		signal: make(chan int, 5),
		mgr:    server,
	}
}

type RemoteLogWorker struct {
	mgr    *RemoteLogServer
	data   chan *LogData
	close  bool
	ww     *sync.WaitGroup
	stmt   *sql.Stmt
	signal chan int
}

func (worker *RemoteLogWorker) Close(ww *sync.WaitGroup) {
	worker.ww = ww
	worker.close = true
}

func (worker *RemoteLogWorker) Push(d *LogData) {
	worker.data <- d
}

func (worker *RemoteLogWorker) do() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case l := <-worker.data:
			worker.insert(l)
		case s := <-worker.signal:
			if s == 1 {
				if worker.stmt != nil {
					worker.stmt.Close()
					worker.stmt = nil
				} else {
					seelog.Error("Do Stmt Close Error ")
				}
			}
		case <-ticker.C:
			if worker.close && len(worker.data) == 0 {
				worker.ww.Done()
				worker.stmt.Close()
				return
			}
		}
	}
}

type LogData struct {
	Uid    uint64 `json:"uid"`
	OpenId string `json:"open_id"`
	Stamp  int64  `json:"stamp"`
	Kind   int    `json:"kind"`
	Source int    `json:"source"`
	Param  string `json:"param"`
}

func (worker *RemoteLogWorker) insert(l *LogData) error {
	if worker.stmt == nil {
		var e error
		worker.stmt, e = worker.mgr.db.Prepare(`INSERT INTO ` + worker.mgr.nowTable + `(uid,openId,stamp,kind,source,param) VALUES (?,?,?,?,?,?)`)
		if e != nil {
			seelog.Error("prepare Error ", e)
			return e
		}
	}
	_, err := worker.stmt.Exec(l.Uid, l.OpenId, l.Stamp, l.Kind, l.Source, l.Param)
	if err != nil {
		seelog.Error("Insert Error ", err)
		return err
	}
	return nil
}
