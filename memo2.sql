-- MySQL dump 10.13  Distrib 8.4.3, for Linux (x86_64)
--
-- Host: localhost    Database: memo2
-- ------------------------------------------------------
-- Server version	8.4.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `books`
--

DROP TABLE IF EXISTS `books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `books` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `desc` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `books`
--

LOCK TABLES `books` WRITE;
/*!40000 ALTER TABLE `books` DISABLE KEYS */;
INSERT INTO `books` VALUES (1,'Golang基础','- GMP\r\n- GC\r\n- channel, select\r\n- sync.Mutex\r\n- slice/map\r\n'),(2,'MySQL',''),(3,'Redis',''),(5,'计算机网络',''),(6,'操作系统',''),(7,'项目经历',''),(8,'目前的需要巩固的知识',''),(13,'zmj',''),(14,'看书','');
/*!40000 ALTER TABLE `books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questions`
--

DROP TABLE IF EXISTS `questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questions` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `book_id` smallint unsigned DEFAULT NULL,
  `name` longtext NOT NULL,
  `answer` longtext NOT NULL,
  `audio` longtext,
  `image` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_questions_book_id` (`book_id`),
  CONSTRAINT `fk_questions_book` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questions`
--

LOCK TABLES `questions` WRITE;
/*!40000 ALTER TABLE `questions` DISABLE KEYS */;
INSERT INTO `questions` VALUES (1,1,'new和make','# 区别\r\n\r\nnew用于`值类型`的初始化，比如结构体，数组和基础值类型。返回的是内存的指针。当你需要获得一个值类型的指针，并将其指向零值的时候，可以使用new。\r\n\r\nmake主要用于分配和初始化`引用类型`的数据结构，如slice, map, channel。\r\n\r\n# 堆栈\r\n\r\n分配在堆还是栈上，取决于编译器的逃逸分析。如果变量离开了当前作用域，或者申请的内存过大，则会分配到堆上。如果没有离开当前作用于，且申请的内存较小，则会分配到栈上。','static/b/1/q/1/audio.mp3',''),(2,1,'OOM','# OOM\r\n\r\nOOM全称是Out of Memory，是操作系统对那些内存占用过大的进程进行的一种强制杀死的操作。当Go程序分配的内存过大的时候，操作系统就会触发OOM（不过这取决于Linux系统的systemctl配置，也可以配置为不触发，但是如果配置为不触发的话可能会导致系统卡死，甚至连ssh远程登录都登录不进去，所以我一般都是开启这个oom配置的）\r\n\r\n# 怎么处理\r\n\r\n如果程序进程被OOM杀死了，怎么处理这种情况？\r\n首先呢，需要确认程序是否是因为OOM导致的异常退出，这个我们可以使用linux的dmesg工具来查看操作系统的日志。当我们确定了oom导致的异常退出之后，\r\n我们通过pprof工具采集的样本对程序进行性能分析和错误排查。因为一般我们线上的程序会开启pprof，周期性的采集程序运行情况,然后导出为日志文件。这就是我们对OOM的处理方式。\r\n[link](https://gitlab.com/-/snippets/3756025#pprof)','static/b/1/q/2/audio.mp3',''),(3,1,'PProf排查','[pprof代码](https://gitlab.com/-/snippets/3756025#pprof)\r\n\r\nGo语言标准库提供了丰富的性能分析工具，我们平时经常用到的是两种。\r\n\r\n第一种方式是在线上的Go程序开启pprof采样，并且周期性地打印内存堆栈信息，和CPU运行情况，比如说每隔5分钟统计一下，不会影响性能。这些数据会导出到本地日志文件夹里面，一旦程序出现内存泄漏，内存占用异常增大，或者意外崩溃了，我们就可以通过这些统计好的数据来进行性能剖析和问题排查。\r\n\r\n另一种方式是，如果你想要更全面地分析并发性能问题，goroutine的调度，系统调度和GC暂停等runtime信息，可以在线上程序中开启runtime/trace。但是pprof和trace这两种工具一般不建议同时开启，具体用哪个这取决于你要分析什么样的问题。','static/b/1/q/3/audio.mp3',''),(4,1,'错误处理','Go的错误处理分为两种\r\n\r\n一种是通过error返回值的方式来表示是否出错了，如果返回错误的层级比较多，我们可以对错误进行层层包装，比如说利用fmt.Errorf()函数来对错误进行二次包装，然后把错误继续传递给上层函数。像我们平时会在代码里面把一些固定的常见的错误定义到整个包的全局变量里面，这样别的地方在调用的时候，可以通过errors.Is()函数来判断返回的这个error是否属于某个常见的公共错误，这就相当于对错误进行拆包对比。\r\n\r\n当然我们还可以自己定义一个struct类型来实现error方法，对于这种自定义错误类型，就需要利用errors.As()函数来进行类型断言。\r\n\r\n另一种方式就是通过panic抛出异常，然后在最外层的defer语句中调用recover()函数来进行异常捕捉。但是我们平常应该在程序中应该尽量避免使用这种方式来传递错误，这种一般是在遇到程序无法恢复的异常情况才会使用的。','static/b/1/q/4/audio.mp3',''),(5,1,'锁','# Mutex锁\r\n\r\nMutex是Go语言最常用的互斥锁，当多个goroutine去争夺同一个公共资源的时候使用。\r\n\r\nMutex底层有两个int32字段组成，一个是state字段，存储了Mutex锁的状态，另一个是semaphore字段，用于存储锁使用的信号量ID。state字段呢，一共32个比特位，不同的比特位存储了不同的信息。比如最右边的三个比特位分别存储的是：是否处于饥饿状态，是否处于唤醒状态，是否处于上锁状态。其他剩下的比特位一起用来存储等待队列的goroutine数量。\r\n\r\n上锁的过程是这样的：先通过CAS的方式进行自旋上锁，如果成功则进入上锁状态，如果失败了则放入等待队列，阻塞挂起。\r\n解锁的过程是这样的：如果等待队列里面还有goroutine，则从等待队列里面随机唤醒一个goroutine，让它重新与新来的goroutine随机竞争，通过自旋来上锁。如果等待队列没有goroutine了，则直接解锁。\r\n\r\n因为正常模式下，等待队列是随机唤醒的，所以有可能有些goroutine一直处在在等待中，没有被唤醒。时间长了之后，就会触发Mutex的饥饿模式。在饥饿模式下，新来的goroutine不自旋，直接休眠。那些从等待队列里面唤醒的goroutine也不需要自旋了，直接获取锁。饥饿模式一直持续到等待队列为空为止。\r\n\r\n# RWMutex\r\n\r\nMutex一般用于读写差不多的情况，然而对于读多写少的情况，建议使用RWMutex。\r\n','static/b/1/q/5/audio.mp3',''),(6,1,'channel select','# 介绍\r\n\r\nchannel是Go语言中实现CSP并发模型的最重要的数据结构，他是一个管道数据类型，你可以往里面发数据，然后在另一个goroutine里面接收数据。一般用于并发编程中，goroutine之间相互通信的场景。因为Go语言不建议通过共享内存来通信。\r\n\r\nchannel的底层实现，主要包含一个Ring Buffer环形缓冲区（用于存储管道缓冲区的数据），还有接收端和发送端的阻塞队列，然后还有一个Mutex互斥锁，再加一个closed字段用于记录管道是否已经被关闭了。\r\n\r\n发送过程包含三种情况：\r\n- 直接发送，不走缓冲区，数据直接拷贝过去\r\n- 放入缓冲区\r\n- 阻塞等待\r\n\r\n接收过程包含三种情况：\r\n- 直接接收，不走缓冲区\r\n- 读缓冲区，发送端如果人在等待则唤醒它\r\n- 阻塞等待\r\n\r\nselect分两种情况\r\n- 有default条件的为非阻塞select，主要是看多种条件里面有没有能立即执行的，如果有则立即执行，如果没有则直接跳过\r\n- 没有default条件的为阻塞select，Go程序会从那些有数据的channel里面随机选择一个条件执行。如果所有channel全都读不出数据，则阻塞在那里，一直等到能读出数据为止。\r\n\r\n# 坑\r\n\r\nchannel和select在使用的时候，比较常见的坑是阻塞问题和close关闭问题。比如：发送端忘记关闭channel导致接收端一直在阻塞等待数据。\r\n','static/b/1/q/6/audio.mp3',''),(7,1,'GMP','- G代表goroutine，是Go语言运行调度的最小单位，相比操作系统的线程要小很多，非常轻量。\r\n- P代表操作者，调度者。维护了一个g的本地队列，负责往M里面运送g\r\n- M代表操作系统的线程，是goroutine的实际执行者，\r\n\r\n如果一个P的队列已经执行完了，会从其他P队列里面偷走一半来执行\r\n# 切换时机\r\n\r\n- 主动挂起time.Sleep()\r\n- 系统调用syscall()\r\n- 函数切换\r\n- 定时切换，g运行超过10ms\r\n- 强制暂停，doGC()','',''),(8,1,'GC','# 三色标级法\r\n三色标级法，为了实现并发执行。极大减少了暂停时间。\r\n\r\n- 白色：已处理，没用\r\n- 灰色：没有遍历完，待处理\r\n- 黑色：已处理，有用\r\n\r\n## 写屏障\r\n\r\n- 插入屏障：新对象标记为灰色。确保新对象能正确参与GC标记。\r\n- 删除屏障：删除时只标记为灰色，先不删除。等GC遍历完再删。\r\n','',''),(9,2,'备份与恢复','- 逻辑备份\r\n  - mysql_dump\r\n- 物理备份\r\n  - 文件方式: xtrabackup\r\n  - 快照方式: btrfs','',''),(10,2,'事务的隔离级别','- 读未提交\r\n- 读已提交\r\n- 可重复读\r\n- 可串行化','',''),(11,2,'日志','任何更新操作都会产生以下3种log，但是读操作不会有\r\n\r\n- binlog是 Server 层生成的日志，主要用于数据备份和主从复制；\r\n- redolog是 Innodb 存储引擎层生成的日志，实现了事务中的持久性，主要用于掉电等故障恢复\r\n- undolog是 Innodb 存储引擎层生成的日志，实现了事务中的原子性，主要用于事务回滚和 MVCC。','',''),(12,2,'B+树和B树','\r\n- B树在所有节点都存储数据，所以MongoDB推荐反范式数据库设计，减少访问叶子节点。\r\n- B+树只在叶子节点存储数据，并且每个叶子节点前后相连，方便进行范围查询','',''),(13,2,'锁','# 全局锁\r\n\r\n用于全库逻辑备份\r\n```sql\r\nflush tables with read lock\r\n```\r\n\r\n# 表级锁\r\n\r\n锁一张表\r\n```sql\r\nlock table users write;\r\n```\r\n\r\n# 行级锁\r\n\r\n- Record Lock记录锁，`select for update`\r\n- Gap Lock间隙锁，防止幻读\r\n- Next-Key Lock临键锁=Record lock+Gap lock','',''),(14,2,'Buffer Pool','# 与查询缓存的区别\r\n\r\n- 查询缓存针对的是SQL语句的查询结果\r\n- Buffer Pool缓存的是数据页','',''),(15,3,'持久化','# RDB\r\n\r\nDump整个内存的`全量快照`。[默认开启]\r\n\r\n执行过程中，Redis也能继续工作，因为有`写时复制`机制\r\n# AOF\r\n\r\nAppend-Only File，每一条写操作执行成功后，才记录日志到aof_buf缓存。不会阻塞写操作。\r\n\r\n写入硬盘时机\r\n- Always\r\n- Everysec\r\n- No，交由操作系统写入硬盘','',''),(16,3,'数据类型','---\r\n\r\n- String，分布式锁\r\n- List，消息队列，\r\n- Hash，\r\n- Set，交集/并集/合集\r\n- ZSet（Sorted Set），排行榜\r\n- Pub/Sub，订阅\r\n- Stream，消息队列\r\n- GEO，地理位置计算\r\n- HyperLogLog，网页UV统计\r\n- BitMap，签到统计\r\n','',''),(17,3,'过期删除策略','- Timer删除，每个expire设置一个Timer\r\n- 惰性删除，下次读取的时候删【Redis默认】\r\n- 定期删除，每隔一段时间`随机`取一些数据删【Redis默认】','',''),(18,3,'内存淘汰策略','当内存占满的时候，淘汰掉不用的key\r\n\r\n- LRU\r\n- LFU\r\n- Radom\r\n- volatile-ttl，优先淘汰更早过期的','',''),(19,3,'哨兵模式','当master节点挂了的时候，slave节点顶上去','',''),(20,3,'缓存雪崩，击穿，穿透','- 雪崩，大量key同时过期\r\n- 击穿，热key突然过期\r\n- 穿透，404 Not Found','',''),(21,5,'TCP','TCP特性\r\n- 面向连接\r\n- 可靠，重传，ACK，拥塞控制\r\n- 字节流，消息分包，有序\r\n\r\n（四元组）才能唯一确定一个TCP连接，而UDP只需要二元组\r\n\r\n- ![img](https://cdn.xiaolincoding.com//mysql/other/format,png-20230309230534096.png)\r\n\r\n# 三次握手\r\n\r\n- SYN+client_seqnum->\r\n- <-SYN+ACK ack_num+1 & server_seqnum+1\r\n- ACK+data ack_num+1->\r\n\r\n# 四次挥手\r\n\r\n- FIN->\r\n- <-ACK\r\n- <-FIN\r\n- ACK->\r\n\r\nTime wait','',''),(22,5,'OSI网络模型','- 应用层\r\n- 传输层。TCP/UDP\r\n- 网络层。路由器，IP/ICMP\r\n- 数据链路层。交换机，MAC/ARP\r\n- 物理层','',''),(23,3,'缓存','- Cache Aside, Redis+Application\r\n- Write Through, SQL Proxy\r\n- Write Back, CPU脏','',''),(24,3,'分布式锁','# 加锁\r\n\r\nSetNX\r\n```go\r\ncli.SetNX(context.Background(), \"hello\", \"clientID-1\", time.Second*30)\r\n```\r\n\r\n# 解锁\r\n\r\nLua\r\n```go\r\nresult:=cli.Eval(\r\n    context.Background(), \r\n    \"if redis.call(\'GET\',KEYS[1]) == ARGV[1] then return redis.call(\'DEL\',KEYS[1]) else return 0 end\", \r\n    []string{\"hello\"},\r\n    \"clientID-1\",\r\n)\r\n```\r\n','',''),(25,6,'虚拟内存','为了隔离进程资源\r\n\r\n# 多级页表\r\n\r\n分页4K，Page Table存储在内存中，MMU负责转换虚拟地址->物理地址\r\n\r\n# TLB\r\n\r\nTranslation Lookaside Buffer。是CPU芯片中的MMU单元所集成的Page Cache。','',''),(26,6,'进程间通信','- Socket\r\n- Pipeline\r\n- 共享内存\r\n- 消息队列\r\n- 信号\r\n- 信号量\r\n','',''),(27,6,'epoll','IO多路复用，注册TCP handler函数，','',''),(28,7,'go-zero','','',''),(29,7,'ClickHouse大数据','','',''),(30,7,'思维导图','','',''),(31,7,'自我介绍','# 关键词\r\n\r\n英语，电商，框架设计go-zero，大数据','',''),(32,8,'工作经验少','','',''),(33,8,'SQL优化','','',''),(34,8,'分布式锁','# RedLock','',''),(35,8,'Redis集群','','',''),(36,8,'微服务','# 限流\r\n\r\n# ','',''),(37,8,'Go语言','','',''),(38,8,'Redis缓存同步策略','','',''),(39,5,'HTTPS','# 四次握手\r\n\r\n前2步有几个目的\r\n- 选择合适的加密套件\r\n- 交换random\r\n- 验证服务端的证书\r\n\r\n后2步只有一个目的\r\n- exchange pre-master随机数\r\n\r\n最终根据3个随机数，生成一个master key会话密钥，用于后续的AES对称加密。\r\n\r\n![img](https://cdn.xiaolincoding.com/gh/xiaolincoder/ImageHost4@main/%E7%BD%91%E7%BB%9C/https/tls%E6%8F%A1%E6%89%8B.png)\r\n\r\n# 加密算法\r\n\r\n## AES\r\n对称加密\r\n## RSA\r\n非对称加密\r\n## ECDSA\r\n\r\n类型: 数字签名算法，基于椭圆曲线的非对称加密\r\n\r\n工作原理: 使用椭圆曲线加密算法（ECC）来生成和验证数字签名。生成签名时使用私钥，验证签名时使用公钥。\r\n\r\n常见曲线: secp256k1（如比特币使用）、P-256、P-384等。\r\n\r\n应用场景: 数字签名、身份验证、区块链技术（如比特币的交易签名）。\r\n\r\n优点: 相比RSA，在提供相同安全性的情况下密钥长度更短、计算更高效，适合资源有限的设备（如移动设备）。\r\n\r\n缺点: 理解和实现复杂度较高，对随机数生成的质量要求较高。\r\n\r\n','',''),(40,5,'WireShark','','',''),(41,5,'VPN协议','# WireGuard\r\n\r\n新型的VPN协议，采用最先进的加密技术，比IPSec和OpenVPN更简单，而且性能更好\r\n\r\n协议类型: 一种较新的、简洁的 VPN 协议，工作在网络层（第3层）。\r\n\r\n使用场景: 由于其简洁性和速度，逐渐在个人和企业中获得广泛使用。\r\n\r\n安全性: 使用现代加密协议（如 ChaCha20、Poly1305），代码库较小，降低了漏洞风险。\r\n\r\n性能: 由于设计简洁、代码优化，与 IPSec 和 OpenVPN 相比速度更快。\r\n\r\n兼容性: 支持主要操作系统，包括 Linux（集成在 Linux 内核中）、Windows、macOS、iOS 和 Android。\r\n\r\n# IPSec\r\n\r\n协议类型: 工作在网络层（第3层），通过对每个 IP 数据包进行身份验证和加密来保护 IP 通信。\r\n\r\n使用场景: 常用于站点到站点 VPN 和远程访问 VPN。\r\n\r\n安全性: 提供强大的安全性，包括加密（如 AES、DES）、身份验证（如 HMAC）以及密钥交换（如 IKE、IKEv2）。\r\n\r\n性能: 由于其复杂性，尤其是在使用强加密时，性能可能较慢，特别是对于旧硬件。\r\n\r\n兼容性: 广泛支持各种平台和设备，是企业环境中的常见选择。\r\n\r\n# OpenVPN\r\n\r\n协议类型: 开源协议，使用 SSL/TLS 进行加密，可以在 UDP 和 TCP 上运行。\r\n\r\n使用场景: 常用于个人或小型企业的远程访问 VPN。\r\n\r\n安全性: 配置得当时非常安全，支持多种加密算法，包括 AES、RSA 等。\r\n\r\n性能: 性能因配置而异，但由于其复杂性，通常比 WireGuard 慢。\r\n\r\n兼容性: 非常灵活和可定制，支持 Windows、macOS、Linux、iOS 和 Android 等多个平台。\r\n\r\n# 总结\r\nWireGuard 因为性能高、部署简单，通常是现代VPN解决方案的首选。而OpenVPN更适合需要灵活配置的场景，IPSec则在企业级站点对站点连接中被广泛使用，强调强大的安全性和兼容性。\r\n','',''),(55,13,'项目展示','','',''),(56,13,'pic','','',''),(57,14,'刻意练习','','',''),(58,14,'李嘉诚传','','','');
/*!40000 ALTER TABLE `questions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-04-12 23:24:23
