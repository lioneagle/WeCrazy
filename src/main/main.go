package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gin-gonic/gin"
	"github.com/lioneagle/goutil/src/algorithm/timewheel"
)

var DB = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

	return r
}

type A struct {
	x int
	b *B
}

type B struct {
	y int
	z float32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var old time.Time
var count int = 0

func tm1Expire(val interface{}) {
	count++
	//if count%1000 == 0 {
	now := time.Now()
	fmt.Println("tm1 expire: ", now, ", delta =", now.Sub(old).Seconds()*1000.0)
	old = now
	/*if count > 1100 {
		panic("finish")
	}*/
	//}
}

func TimeWheelTimer() {
	//tw := timewheel.NewTimeWheel(3, []int{1000, 60, 60}, int64(time.Millisecond)*1,
	tw := timewheel.NewTimeWheel(3, []int{10000, 600, 600}, int64(time.Millisecond)*1,
		int64(time.Now().UnixNano()), 1000)
	//fmt.Println("start: ", time.Now())
	//tw.AddCycle(int64(time.Second*13/10), nil, tm1Expire)
	tw.AddCycle(int64(time.Millisecond*10), nil, tm1Expire)

	go func(tw *timewheel.TimeWheel) {
		old = time.Now()

		for {
			//fmt.Println("now =", time.Now())
			tw.Step(int64(time.Now().UnixNano()))
			//time.Sleep(time.Millisecond * 1)
		}
	}(tw)
}

func GoTimer() {
	ticker := time.NewTicker(time.Millisecond * 10)
	go func(ticker *time.Ticker) {
		old = time.Now()
		for {
			//t1 := <-ticker.C
			<-ticker.C
			//fmt.Println("t1 =", t1)
			//fmt.Println("t1 =", t1, "now =", time.Now())
			count++
			//if count%100 == 0 {
			now := time.Now()
			fmt.Println("now =", now, ", delta =", now.Sub(old).Seconds()*1000.0)
			old = now
			//}
			//tw.Step(int64(time.Now().UnixNano()))
			//time.Sleep(time.Millisecond * 1)
		}
	}(ticker)
}

func main() {
	TimeWheelTimer()

	/*go func() {
		for {
			fmt.Println("please input a string")
			var str string
			fmt.Scanln(&str)
			fmt.Println("input =", str)
		}
	}()*/

	//GoTimer()

	for {
		time.Sleep(time.Second * 1)
	}

	return

	db, err := sql.Open("mysql", "root:s74760302H1@/test?charset=utf8")
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT userinfo SET username=?,department=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("test", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)

	//更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	res, err = stmt.Exec("lioneagle", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("affect =", affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}
	//删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)
	res, err = stmt.Exec(id)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)

	return
	//a := new(A)
	a := &A{}
	fmt.Println("a =", a)

	return
	fmt.Println("WeCrazy")
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run("localhost:8080")
}
