
# 任务
1. 理解map和slice内存扩容的区别
	- slice 内存扩容时申请一个内存容量为原来两倍的新内存，将原来的数据复制到新的内存中去
	- 当 map 增长到容量上限的时候，如果再增加新的 key-value 对，map 的大小会自动加 1。扩容方式有两种：
		i. 翻倍扩容（增量扩容）：因为元素太多了（装载因子过高），新申请的桶数量 = 旧桶数量 × 2。
		ii. 等量扩容（等大扩容）：内存总量不变！ 如果 Map 里元素不多，但是存在大量因为删除操作留下的“空洞”（溢出桶太多），导致查询变慢。这时候 Go 会申请和原来一模一样大的新桶数组。目的不是为了扩大内存，而是为了把数据重新排列紧凑，消除碎片。
		- 不会将旧桶里的数据立刻搬到新桶里，而是在以后的对这个Map进行增、删、改、查操作时，顺便将旧桶里的1到2个桶的数据挪到新桶里去

		- 为什么 Map 不学 Slice 一次性拷贝？
		- 假设你的 Map 里有 1000 万个键值对，如果像 Slice 那样一次性拷贝并重新计算哈希（Hash），主程序可能会卡顿好几秒，这在 Web 服务或高并发场景下是不可接受的。所以 Go 采用了“化整为零”的策略，把扩容的代价分摊到了成百上千次后续的操作中。

2. 利用闭包实现一个计数器
	
	```go fold
	package main

  

	import "fmt"

  

	func main() {

	    count1 := Counter()

	    var sum int

  

	    for i := 0; i < 5; i++ {

	        sum = count1()

	    }

  

	    fmt.Printf("The result is %d.", sum)

	}

  

	func Counter() func() int {

	    count := 0

	    return func() int {

	        count += 1

	        return count

	    }

	}
	```

	- 错误
		- 作用域越界：你在 for 循环里面用 := 定义了 sum。在 Go 语言里，:= 定义的变量只在当前代码块（也就是这个 for 循环）内有效。出了循环，sum 就销毁了。所以下面 fmt.Println 根本不认识 sum 是谁（未定义）。
		- 打印格式错误：fmt.Println 不支持格式化占位符（像 %d 这种）。如果你要格式化输出，得用 fmt.Printf；如果是 Println，直接用逗号隔开就行。


3. 基于下面这个结构体实现四个方法：
    

- Average：返回平均分；`Scores` 为空时返回 `0`。
    
- AddScore：追加一个成绩，**调用方能看见这个改动**。
    
- Top：返回最高的 `n` 个成绩，从高到低。`n` 超过长度时返回全部，`n <= 0` 时返回 `nil`。**不能修改** **`s.Scores`** **本来的顺序，也不能让返回的切片存在影响原数据的可能。**
    
- Merge：合并两个人的成绩，返回一个新的 `Student`（`Name` 用接收者的）。**合并结果和两个原对象必须互不影响。**
    

	```Go
	type Student struct {
	    Name   string
	    Scores []int // 记录若干项成绩
	}
	```

```go
package main

  

import (

    "fmt"

)

  

type Student struct {

    Name   string

    Scores []int // 记录若干项成绩

}

  

func main() {

    s1 := Student{Name: "A", Scores: []int{80, 90, 100}}

    s2 := Student{Name: "B", Scores: []int{70}}

  

    fmt.Println("test Average", s1.Average())

    s1.AddScore(60)

    fmt.Println("test AddScore", s1.Scores)

    fmt.Println("test Top", s1.Top(2))

    fmt.Println("test Merge", s1.Merge(s2))

  

}

  

func (s Student) Average() float64 {

    if len(s.Scores) == 0 {

        return 0.0

    }

    sum := 0

    for _, score := range s.Scores {

        sum += score

    }

  

    return float64(sum) / float64(len(s.Scores))

}

  

func (s *Student) AddScore(score int) {

    s.Scores = append(s.Scores, score)

}

  

func (s Student) Top(n int) []int {

    if n <= 0 {

        return nil

    }

    length := len(s.Scores)

  

    temp := make([]int, length)

    copy(temp, s.Scores)

  

    for i := 0; i < length-1; i++ {

        for j := 0; j < length-i-1; j++ {

            if temp[j] < temp[j+1] {

                temp[j], temp[j+1] = temp[j+1], temp[j]

            }

        }

    }

  

    if n > length {

        n = length

    }

  

    return temp[:n]

}

  

func (s Student) Merge(other Student) Student {

    NewScores := make([]int, 0, len(s.Scores)+len(other.Scores))

    NewScores = append(NewScores, s.Scores...)

    NewScores = append(NewScores, other.Scores...)

  

    return Student{

        Name:   s.Name,

        Scores: NewScores,

    }

}
```

4. 思考以下程序在并发出现panic的原因，并解决
    

```Go
package main

type message struct {
    Topic     string
    Partition int32
    Offset    int64
}

type FeedEventDM struct {
    Type    string
    UserID  int
    Title   string
    Content string
}

type MSG struct {
    ms        message
    feedEvent FeedEventDM
}

const ConsumeNum = 5

func main() {
    var consumeMSG []MSG
    var lastConsumeTime time.Time // 记录上次消费的时间
    msgs := make(chan MSG)

    //这里源源不断的生产信息
    go func() {
       for i := 0; ; i++ {
          msgs <- MSG{
             ms: message{
                Topic:     "消费主题",
                Partition: 0,
                Offset:    0,
             },
             feedEvent: FeedEventDM{
                Type:    "grade",
                UserID:  i,
                Title:   "成绩提醒",
                Content: "您的成绩是xxx",
             },
          }
          //每次发送信息会停止0.01秒以模拟真实的场景
          time.Sleep(100 * time.Millisecond)
       }
    }()

    //不断接受消息进行消费
    for msg := range msgs {
       // 添加新的值到events中
       consumeMSG = append(consumeMSG, msg)
       // 如果数量达到额定值就批量消费
       if len(consumeMSG) >= ConsumeNum {
          //进行异步消费
          go func() {
             m := consumeMSG[:ConsumeNum]
             fn(m)
          }()
          // 更新上次消费时间
          lastConsumeTime = time.Now()
          // 清除插入的数据
          consumeMSG = consumeMSG[ConsumeNum:]
       } else if !lastConsumeTime.IsZero() && time.Since(lastConsumeTime) > 5*time.Minute {
          // 如果距离上次消费已经超过5分钟且有未处理的消息
          if len(consumeMSG) > 0 {
             //进行异步消费 
             go func() {
                m := consumeMSG[:ConsumeNum]
                fn(m)
             }()
             // 更新上次消费时间
             lastConsumeTime = time.Now()
             // 清空插入的数据
             consumeMSG = consumeMSG[ConsumeNum:]
          }
       }
    }
}

func fn(m []MSG) {
    fmt.Printf("本次消费了%d条消息\n", len(m))
}
```

```go
if len(consumeMSG) >= ConsumeNum {
          //进行异步消费
          go func() {
             m := consumeMSG[:ConsumeNum]
             fn(m)
          }()
          // 更新上次消费时间
          lastConsumeTime = time.Now()
          // 清除插入的数据
          consumeMSG = consumeMSG[ConsumeNum:]
       } 
```
- 主线程打开协程后没有暂停，继续执行下面的语句，

5. 使用for循环生成20个goroutine,每个goroutine随机休眠0~1000ms，并向一个channel传入随机数和goroutine编号(从1-20)，等待这些goroutine都生成完后，想办法给这些随机数**按照编号进行排序**(输出排序前和排序后的结果,要求不使用额外的空间存储这20个数据)
```go fold
package main

  

import (

    "fmt"

    "math/rand"

    "sync"

    "time"

)

  

type Thing struct {

    ID  int

    NUM int

}

  

func main() {

    ch := make(chan Thing, 20)

    var wg sync.WaitGroup

  

    for i := 1; i <= 20; i++ {

        wg.Add(1)

  

        go func(id int) {

            defer wg.Done()

  

            sleeptime := rand.Intn(1001)

            time.Sleep(time.Duration(sleeptime) * time.Millisecond)

  

            ch <- Thing{

                ID:  id,

                NUM: rand.Intn(101),

            }

        }(i)

    }

    wg.Wait()

    close(ch)

  

    var arr []Thing

    for thing := range ch {

        arr = append(arr, thing)

    }

    fmt.Println("排序前", arr)

    Sort(arr)

}

  

func Sort(arr []Thing) {

    for i := 0; i < 20; i++ {

        for arr[i].ID != i+1 {

            arr[i], arr[arr[i].ID-1] = arr[arr[i].ID-1], arr[i]

        }

    }

    fmt.Println("排序后", arr)

}
```

[[random_3.go]]

6. 拓展题：只使用channel进行并发控制，交叉打印下面两个字符串（要求一个打印完，另一个会继续打印） "ABCDEFGHIJKLMNOPQRSTUVWXYZ" "0123..." 得到："AB01CD23EF34..."（不会的话可以降低难度，使用其他的并发控制工具）

[[ChannleConcurrence_2.go]]
# 相关知识点
## sync 包
1. sync 包里的核心结构体

-  sync.WaitGroup：你刚才用的，等待一组协程结束。底层是一个结构体。
-  sync.Mutex：互斥锁。最常用的并发控制工具。
-  sync.RWMutex：读写锁。允许多个读，但只允许一个写。
-  sync.Once：保证某段代码在并发环境下只执行一次（比如单例模式初始化）。
-  sync.Pool：临时对象池，用来复用内存，减轻 GC（垃圾回收）压力。
-  sync.Map：并发安全的 Map（普通 Map 不支持并发写）。

(注：sync 包里除了结构体，还有几个极少用的底层函数，比如 sync.OnceFunc，但 99% 的日常开发都在和这些结构体打交道。)

2. 传递
- 所有sync包里的结构体，绝对不能通过值传递，必须传指针


## time 包
time 包主要解决四个问题：获取当前时间（Time）、计算时间差（Duration）、格式化输出、以及并发定时控制（Timer/Ticker）。

---

一、 核心一：时间点（time.Time）

time.Time 是一个结构体，代表绝对的时间（年月日时分秒）。

1. 获取当前时间与时间戳

```go
now := time.Now()            // 获取当前时间
fmt.Println(now)             // 2023-11-14 15:30:00 +0800 CST

fmt.Println(now.Unix())      // 秒级时间戳（10位数字，极其常用）
fmt.Println(now.UnixNano())  // 纳秒级时间戳（19位数字）
```

2. ⚠️ 最反直觉的设计：时间格式化
Go 语言没有 YYYY-MM-DD 这种写法。它规定了一个“魔法时间”作为模板：2006-01-02 15:04:05（记忆口诀：1 2 3 4 5 6，即 1月 2日 3点(15) 4分 5秒 06年）。

```go
// 将时间转为字符串（按你想展示的格式，把魔法时间替换成对应的数字）
str := now.Format("2006-01-02 15:04:05")

// 将字符串解析为时间
t, err := time.Parse("2006-01-02 15:04:05", "2023-11-14 15:30:00")
```

---

二、 核心二：时间段（time.Duration）

time.Duration 表示“过了多久”，底层是 int64，单位是纳秒。Go 提供了一组直观的常量：

```go
time.Nanosecond  // 1 纳秒
time.Microsecond // 1 微秒
time.Millisecond // 1 毫秒
time.Second      // 1 秒
time.Minute      // 1 分钟
time.Hour        // 1 小时
```

时间的加减与比较：

```go
// 1. 加减时间（Add 接收 Duration 参数）
t1 := now.Add(2 * time.Hour)     // 2小时后
t2 := now.Add(-30 * time.Minute) // 30分钟前

// 2. 计算时间差（Sub 返回 Duration 类型）
diff := t1.Sub(t2) 
fmt.Println(diff.Seconds())      // 获取相差了多少秒

// 3. 比较时间（⚠️ 绝对不要用 ==，因为含有单调时钟信息）
now.Before(t1) // now 在 t1 之前吗？
now.After(t1)  // now 在 t1 之后吗？
now.Equal(t1)  // 两个时间是否相等？
```

---

三、 核心三：并发与定时器

这是 time 包在并发编程中大放异彩的地方。

1. time.Sleep（睡眠）
让当前协程（Goroutine）阻塞指定的时间。

```go
time.Sleep(2 * time.Second) // 程序在这停 2 秒
```

2. time.After（超时控制，极常用）
返回一个 channel，在指定时间后往 channel 发送当前时间。经常配合 select 用来做防阻塞的超时控制。

```go
select {
case res := <-ch:
    fmt.Println("任务正常完成:", res)
case <-time.After(3 * time.Second):
    fmt.Println("超时了，任务没完成")
}
```

3. time.NewTicker（周期定时器）
每隔一段时间触发一次。常用于心跳检测或定时任务。

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop() // ⚠️ 必须写，否则内存泄漏

for t := range ticker.C {
    fmt.Println("每秒触发一次:", t)
}
```

---

四、 避坑指南（面试常考）

1. time.After 在 for 循环里的内存泄漏：
   如果你在死循环里不断调用 time.After，在定时器到期前，底层的计时器是不会被 GC 回收的，会导致内存暴涨。正确做法是使用 time.NewTimer 配合 Reset()。
2. 不要拷贝 time.Time：
   time.Time 包含内部状态（比如单调时钟），直接传值可能导致精度丢失或比较异常。
3. 时区问题：
   服务器默认时区可能是 UTC，如果你直接 Format，会发现时间差 8 小时。开发中通常统一存 UTC 时间，展示时再转成本地时间：
   ```go
   loc, _ := time.LoadLocation("Asia/Shanghai")
   now.In(loc).Format("2006-01-02 15:04:05")
   ```

## math/rand 包
math/rand 是 Go 语言标准库中用于生成伪随机数的包。它适用于模拟、游戏、随机抽样等场景，但绝不能用于对安全性要求高的场景（如生成密码或密钥，这类需求应使用 crypto/rand）。

🎲 核心函数：生成随机数

最常用的是 math/rand 提供的顶层函数，它们使用一个全局的随机数生成器：

· rand.Intn(n)：返回一个 [0, n) 范围内的随机整数。例如 rand.Intn(100) 生成 0 到 99 的整数。
· rand.Float64()：返回一个 [0.0, 1.0) 范围内的随机浮点数。
· rand.Shuffle(n, swap)：原地打乱一个切片。用法如 rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })。

此外还有 rand.Int()（非负随机整数）、rand.Int31()、rand.Int63() 等。

🌱 种子与版本演进（关键）

“种子”决定了随机数序列的起点，相同的种子会产生完全相同的序列。

· Go 1.20 之前：全局生成器的默认种子固定为 1。如果不手动调用 rand.Seed(time.Now().UnixNano())，每次运行程序都会得到同一串“随机数”。
· Go 1.20 及之后：全局生成器在程序启动时会自动随机播种，每次运行的结果都不同。因此，rand.Seed 已被标记为弃用，通常不再需要调用它。

🔒 并发安全性：两种使用方式

· 全局函数（如 rand.Intn）：是并发安全的，可以在多个 goroutine 中直接调用。在 Go 1.20+ 中，不调用 Seed 时，全局函数甚至使用了无锁的生成器，性能更好。
· 自定义实例（rand.New(rand.NewSource(seed))）：不是并发安全的。如果需要在多个 goroutine 中使用自定义实例，必须自己加锁（如 sync.Mutex）。
  · 适用场景：当你需要可重复的随机序列（例如在测试中），或需要独立的随机数流互不干扰时，就应该使用这种方式。

✨ 重要更新：math/rand/v2 (Go 1.22+)

Go 1.22 引入了 math/rand/v2，带来了更现代的改进：

· 新算法：用 PCG 和 ChaCha8 替换了旧生成器，性能更好。
· 接口变更：移除了全局的 Seed 函数，并更改了 Source 接口。
· 包路径：使用时导入 math/rand/v2，函数名有微调（如 rand.IntN 取代了 rand.Intn）。

💎 总结与建议

· 日常使用：直接调用 math/rand 的全局函数（如 rand.Intn）即可，无需自己设置种子，代码简洁且并发安全。
· 需要可重复结果：使用 rand.New(rand.NewSource(seed)) 创建独立实例，但要注意其非并发安全的特性。
· 安全敏感场景：务必改用 crypto/rand。
## make 函数
- 三个参数，分别是类型、长度、容量

## 深拷贝（不影响原数据）
```go
temp := make([]int, length, capacity)
copy(temp, s.Scores)
```
