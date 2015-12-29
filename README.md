tailMail
========

日志文件监控及邮件通知小工具 V0.3


# 功能描述 #
监控日志类型文本文件的变化，如果有新增内容时，截取最新的10000个字符之内的内容，通过配置的邮件发送给相关人。

典型用途
+ 监控系统自动产生的异常日志、业务日志,并定时发送邮件。


# 使用说明 #

## 配置文件 ##
配置文件支持两种格式： toml 和 json
默认是 toml， 目前不支持参数修改，需要修改代码的参数来修改。

+ 配置文件 config.toml 或 config.json 文件。   这里目前发送邮件的配置密码是瞎写的，使用者需要修改成自己的账户密码。

+ 进度文件 progress.toml 或 progress.json。  
测试时进度文件是可以随时删除的， 删除掉意味着下次当新文件来处理。

json格式配置文件的校验请使用： http://jsonlint.com/

toml 文件的注释符是 #    
 # This is a full-line comment    
key = "value" # This is a comment at the end of a line
请参考： https://github.com/toml-lang/toml    

+ template.html 是发送邮件内容的模板文件，可以根据自己的情况进行定制修改。


## 命令参数 ##

-o tailMail 本身的日志输出是否要输出到文件，默认 false， 如果设置成true 则输出到当前执行目录下，每天一个文件。
命令例子：
tailMail -o=true 

-p 配置执行目录， 当我们部署在 crontab 时，由于 crontab 的执行目录是crontab的当前目录，不是我们期望的代码部署目录，这时候需要指定这个参数。
不指定这个参数，执行时，则直接取当前目录。

使用 crontab 配置时，建议 -o -p 都使用。


-i 产生一个参考的配置文件。
注意，这里的邮箱、密码都是瞎写的，需要改成自己需要的。

## crontab 配置例子 ##

命令参数  
/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/cmd/tailMail -o=true -p=/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/cmd

注意，这里路径都需要完整路径，这样才能避免 crontab 当前目录不一样。

*/5 * * * * root /Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/cmd/tailMail -o=true -p=/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/cmd

http://linuxtools-rst.readthedocs.org/zh_CN/latest/tool/crontab.html

## 监控文件模版使用 ##

在配置文件中，如果配置的 FileNameUseTemplate 为 false， 则不启用文件名模版， 如果是 true ，则启用。 比如下面的部分配置就是启用了。
"FileName": "/wangapp/tomcat-wxmember/logs/localhost.{{formatNow \"2006-01-02\"}}.log",
"FileNameUseTemplate": true,

这里的 formatNow 参数接受的值是 go 时间格式描述符。
go 的格式化时间是用的一个特殊的时间做格式化的，参考如下：
2006-01-02T15:04:05

# 变更历史 #

+ V0.1 在畅游时（2013年）完成第0.1版
+ V0.2 在微智全景时(2015-12-21)，由于邮箱使用的是SSL方式才能发送，重构出第0.2版。
+ V0.3 
    + 配置文件默认修改为 toml 格式 (2015-12-24)。
    + 调整截取算法，确保完整的截取一行（如果发生这个逻辑，丢弃第一行）（2015-12－29）。

