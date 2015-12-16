falcon-mailer
====================

falcon-sender 发送邮件到smtp服务器的接口实现，当然你也可以把它当做一个简单的httpmail接口使用也没问题:)

安装
--------------------

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/iambocai
cd $GOPATH/src/github.com/iambocai
git clone https://github.com/iambocai/mailer.git
cd mailer
go get ./...
./control build
# vi cfg.json modify configuration
./control start
# test if works, e.g.: 
# curl http://127.0.0.1:1925/api/mail -XPOST -d 'tos=aaa@abc.com,bbb@mno.com&subject=Hello&content=word'
# if anything goes wrong, check log for reasons
```

配置
-----------------------

- whitelist: 允许发起http调用的客户端IP地址或CIDR格式的范围列表，逗号分隔
- spliter: tos参数中，收件人地址之间的分隔符号
- 其它应该都认识吧。。

FAQ
-----------------------

1. 我想用QQ/163/yahoo/gmail邮箱发信，smtp怎么配置？

   参见[这个文档](http://wenku.baidu.com/link?url=huuX0lkfuJh_44awQwciIVq-S-twddd429dqAzY3jCpsij6eU5aUNrM6rA0ZaUdlsuehRDk_slt2rilJdeOF-8E12tx33aBiK-gbooF1F7G)

   *注意：一定不要漏掉smtp地址中的端口号，比如对于163邮箱，smtp地址应为 smtp.163.com:25*

2. 我想用我公司自己的邮箱服务器发信，smtp怎么配置？

	同上，如不清楚smtp服务器地址和端口，请联系你司邮箱管理员获取。
	
3. 发信失败，请求返回非success

- remote not in whitelist: 发起请求的客户端不在白名单中，请检查配置
- body is blank: 请求方式必须是POST
- connot decode body: 参数不完整，必须有tos，subject，content三个参数，且不为空
- send mail error:  与smtp交互失败，请查看日志获得具体原因

4. 请求返回success，但收件方未收到信

	市面上常见的邮箱产品都有严格的反垃圾邮件规则设置，请检查您的邮件是否被标记为垃圾邮件，以及您的账号是否已被封禁。
	
5. 支持TLS/SSL么？

    暂时不支持，如有需要您可以参照[这篇文章](http://www.oschina.net/code/snippet_166520_34694)