# 木犀课栈

## 课程导入

+ 云课堂课程导入 => history_course
+ 选课手册课程导入 => using_course

导入的脚本已转移至[manage_script](https://github.com/MuxiKeStack/manage_script)

## 关于成绩

成绩从教务处获取，当用户进行如下操作时会触发成绩爬取：

+ 用户加入成绩共享计划
+ 查看个人课程

另外，设置环境变量`MUXIKSTACK_GRADE_CRAWL`决定在**查看个人课程**时是否爬取成绩

```shell
export MUXIKSTACK_GRADE_CRAWL=on  # 爬取，临近期末时设置
export MUXIKSTACK_GRADE_CRAWL=off # 不爬取，下学期开学时设置
```
