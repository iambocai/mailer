falcon-mailer
====================

falcon-sender 发送邮件到smtp服务器的接口实现，当然你也可以把它当做一个简单的httpmail接口使用也没问题:)

##安装
--------------------

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/iambocai
cd $GOPATH/src/github.com/iambocai
git clone https://github.com/iambocai/mailer.git
cd mailer
go get ./...
chmod 755 control
./control build
# vi cfg.json modify configuration
./control start
# test if works, e.g.: 
# curl http://127.0.0.1:1925/api/mail -XPOST -d 'tos=aaa@abc.com,bbb@mno.com&subject=Hello&content=word'
# if anything goes wrong, check log for reasons
```

##配置
-----------------------

- whitelist: 允许发起http调用的客户端IP地址或CIDR格式的范围列表，逗号分隔
- spliter: tos参数中，收件人地址之间的分隔符号
- maxbytes: 发送带附件的邮件时，允许附件最大多大？
- 其它应该都认识吧。。

##API
-----------------------

####发送普通邮件

**URL**: /api/mail

**请求方式**：POST

**Content_type**: application/x-www-form-urlencoded

|参数|必需|含义|备注|
|---|----|---|---|
|tos|是|	邮件收件人列表|多个用英文逗号,分割，分割符号可以在配置中修改|
|subject|是|邮件主题|请使用utf8编码|
|content|是|邮件正文|请使用utf8编码|
|format|否|邮件正文的格式|默认为text，可设置为text或html|
|from|否|发件人|不能含有非ASCII字符，不建议设置|
|server|否|自定义SMTP服务器，格式如：smtp.exmail.qq.com:25|如要自定义则server，user，passwd必须同时设定，否则使用配置中的信息|
|user|否|登陆邮箱使用的用户名，如abc@qq.com|如要自定义则server，user，passwd必须同时设定，否则使用配置中的信息|
|passwd|否|登陆邮箱使用的密码，如123456|如要自定义则server，user，passwd必须同时设定，否则使用配置中的信息|



**请求响应：**

成功：{"status":0,"msg":"ok"}

失败：错误消息

**代码示例**

```shell 
:~$   curl http://127.0.0.1:1925/api/mail -XPOST -d 'tos=iambocai@163.com,test@163.com&subject=Hello&content=word'
:~$   {"status":0,"msg":"ok"}
:~$   curl http://127.0.0.1:1925/api/mail -XPOST -d 'tos=iambocai@163.com,test@163.com&subject=Hello&content=<h1>你好</h1><br/><strong>world</strong>&format=html'
:~$   {"status":0,"msg":"ok"}
```


####发送带附件的邮件

**URL**: /api/attachmail

**请求方式**：POST

**Content_type**: multipart/form-data

|参数|必需|含义|备注|
|---|----|---|---|
|tos|是|	邮件收件人列表|多个用英文逗号,分割，分割符号可以在配置中修改|
|subject|是|邮件主题|请使用utf8编码|
|content|是|邮件正文|请使用utf8编码|
|attachNum|是|附件文件数|必须是大于0的整数N|
|attach1..attachN|是|附件文件|请注意所有文件大小之和不要超过配置中设置的文件大小限制|
|format|否|邮件正文的格式|默认为text，可设置为text或html|
|from|否|发件人|不能含有非ASCII字符，不建议设置|
|server|否|自定义SMTP服务器，格式如：smtp.exmail.qq.com:25|如要自定义则server，user，passwd必须同时设定，否则使用配置中的信息|
|user|否|登陆邮箱使用的用户名，如abc@qq.com|如要自定义则server，user，passwd必须同时设定，否则使用配置中的信息|
|passwd|否|登陆邮箱使用的密码，如123456|如要自定义则server，user，passwd必须同时设定，否则使用配置中的信息|


**请求响应：**

成功：{"status":0,"msg":"ok"}

失败：错误消息

**代码示例**

```python

#!/usr/bin/env python
# encoding: UTF-8
 
from poster.encode import multipart_encode
from poster.streaminghttp import register_openers
import urllib2
# 在 urllib2 上注册 http 流处理句柄
register_openers()
# 开始对文件 "DSC0001.jpg" 的 multiart/form-data 编码
# "image1" 是参数的名字，一般通过 HTML 中的 <input> 标签的 name 参数设置
# headers 包含必须的 Content-Type 和 Content-Length
# datagen 是一个生成器对象，返回编码过后的参数
datagen, headers = multipart_encode({
    "tos":"iambocai@163.com",
    "content":"<h1>hello, world</h1>",
    "subject":"hello",
    "format":"html",
    "attachNum":2,
    "attach1": open("file1.txt", "rb"),
    "attach2": open("file2.txt", "rb")
    })                                                                                                                                                           
# 创建请求对象
request = urllib2.Request("http://127.0.0.1:1925/api/attachmail", datagen, headers)
# 实际执行请求并取得返回
print urllib2.urlopen(request).read()

````



##FAQ
-----------------------

1. 我想用QQ/163/yahoo/gmail邮箱发信，smtp怎么配置？

   参见[这个文档](http://wenku.baidu.com/link?url=huuX0lkfuJh_44awQwciIVq-S-twddd429dqAzY3jCpsij6eU5aUNrM6rA0ZaUdlsuehRDk_slt2rilJdeOF-8E12tx33aBiK-gbooF1F7G)

   *注意：一定不要漏掉smtp地址中的端口号，比如对于163邮箱，smtp地址应为 smtp.163.com:25*

2. 我想用我公司自己的邮箱服务器发信，smtp怎么配置？

	同上，如不清楚smtp服务器地址和端口，请联系你司邮箱管理员获取。
	
3. 发信失败，请求返回非success

- remote not in whitelist: 发起请求的客户端不在白名单中，请检查配置
- body is blank: 请求方式必须是POST
- param error: 参数不完整，必须有tos，subject，content三个参数，且不为空
- send mail error:  与smtp交互失败，请查看日志获得具体原因

4. 请求返回success，但收件方未收到信

	市面上常见的邮箱产品都有严格的反垃圾邮件规则设置，请检查您的邮件是否被标记为垃圾邮件，以及您的账号是否已被封禁。
	
5. 支持TLS/SSL么？

    暂时不支持，如有需要您可以参照[这篇文章](http://www.oschina.net/code/snippet_166520_34694)
