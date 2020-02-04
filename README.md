# 木犀课栈

## 课程导入

每学期选课手册公布时，运维人员执行脚本，手动导入选课手册的课程和云课堂的课程

#### 环境变量

> 一般而言无需改变环境变量

```shell
export MUXIKSTACK_DB_ADDR=127.0.0.1:3306
export MUXIKSTACK_DB_USERNAME=muxi
export MUXIKSTACK_DB_PASSWORD=muxi
```

#### 导入选课手册

导入选课手册要先将Excel文件移动至`script`目录下，然后在`script`目录下执行

```shell
go run add_using.go -file sample.xlsx
```

>ps: 脚本文件要先去除注释

#### 导入云课堂课程

同样，进入`script`目录，在去除注释后运行go文件

## 关于成绩

成绩从教务处获取，当用户进行如下操作时会触发成绩爬取：

+ 用户加入成绩共享计划
+ 查看个人课程

另外，设置环境变量`MUXIKSTACK_GRADE_CRAWL`决定在*查看个人课程*时是否爬取成绩

```shell
export MUXIKSTACK_GRADE_CRAWL=true  # 爬取，临近期末时设置
export MUXIKSTACK_GRADE_CRAWL=false # 不爬取，下学期开学时设置
```
