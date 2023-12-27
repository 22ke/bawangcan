# bawangcan

使用场景：
大众点评的霸王餐批量报名工具

使用教程：
必须三个文件在同一个目录下面，即同一个文件(```main.exe,log.txt,config.lua```)
所以执行的话需要下载这三个文件
main.go源码文件不需要下载

1,修改配置文件：```config.yaml```

```
dper: "*****"
cityname: "shanghai"
#1:美食 2:丽人 3:结婚 4:亲子 5:家装 6:玩乐 7:酒旅 8:培训 9:生活 15:医美 0:全部
menu:
  - 1
#间隔时间，默认连续发，可设置每个报名请求直接间隔几秒
internaltime: 3
debug: false
```
字段解析:

## dper(见20231227更新内容)：
首先获取dper字段
https://www.dianping.com/
从电脑段进去大众点评首页并登录
然后在cookie中找到 deper这个字段，然后把值复制到config.lua的相应位置，双引号不要丢了

![获取dper](https://github.com/22ke/bawangcan2/blob/master/huoqudper.png)

### cityname：
所在城市的城市名拼音，shanghai，shenzhen。。。。。。双引号不要丢了
必须得符合要求，点评要求得账号居住地和报名的城市相同。

### menu：
报名的类别，填数字就可以了，如果填多个数字的话为
menu:
- 1
- 2
- 3
多个时参考yaml格式填写

## 最后双击main.exe就可以了

##20231227更新内容
**********dper字段可以不自己获取了，弹出的浏览器直接用大众点评APP扫码登陆即可
**********debug字段忽略即可，一直改client代理太麻烦了就做了个参数
**********新增字段internaltime，表示每两个报名之前的间隔时间（对美团的风控有没有用我也不知道。。大概率心理作用）
