#bawangcan

使用场景：
大众点评的霸王餐批量报名工具

使用教程：
必须三个文件在同一个目录下面，即同一个文件(```main.exe,log.txt,config.lua```)

1,修改配置文件：```config.lua```

```
config.dper = "7b5"
config.cityname = "shenzhen"
--1:美食2:丽人3:结婚4:亲子5:家装6:玩乐7:酒旅8:培训9:生活15:医美0:全部
config.menu = "0"
```
字段解析:

###dper：
首先获取dper字段
https://www.dianping.com/
从电脑段进去大众点评首页并登录
然后在cookie中找到 deper这个字段，然后把值复制到config.lua的相应位置，双引号不要丢了

![获取dper](https://github.com/22ke/bawangcan2/blob/master/huoqudper.png)

###cityname：
所在城市的城市名拼音，shanghai，shenzhen。。。。。。双引号不要丢了
必须得符合要求，点评要求得账号居住地和报名的城市相同。

$$$menu：
报名的类别，填数字就可以了，如果填多个数字的话为
config.menu = "1,2,3"
数字间用逗号隔开，注意逗号是英文的逗号

最后双击main.exe就可以了
