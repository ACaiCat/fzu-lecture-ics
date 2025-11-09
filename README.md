# fzu-lecture-ics
一个用来获取福大讲座日历ics接口(/v1/lecture/calender)的hertz服务端demo，仅供学习使用。

## 坑
- 被声明但未初始化的map是nil，会爆炸
- 做时间格式化要指定时区
- 服务器的timezone不全，没装go sdk会发生panic
- 小米日历的URL导入只能使用https链接，否则一定导入失败，甚至请求都不发
- ~~把calendar拼成candle~~