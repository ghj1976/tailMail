tailMail
========

日志文件监控及邮件通知小工具


# 功能描述 #
监控日志类型文本文件的变化，如果有新增内容时，截取最新的10000个字符之内的内容，通过配置的邮件发送给相关人。

典型用途
+ 监控系统自动产生的异常日志、业务日志


# 使用说明 #
+ 注意需要修改配置文件 config.json 文件。 这里目前发送邮件的配置密码是瞎写的，使用者需要修改成自己的账户密码。
+ 测试时 progress.json 进度文件是可以随时删除的， 删除掉意味着下次当新文件来处理。
+ template.html 是发送邮件内容的模板文件，可以根据自己的情况进行定制修改。
