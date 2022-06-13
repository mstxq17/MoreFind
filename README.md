# MoreFind
一款用于快速导出URL、Domain和IP的小工具


## 快速安装
方式一: 通过Go包管理安装
```bash
go install  github.com/mstxq17/MoreFind@latest
```
方式二: 直接安装二进制文件
```bash
wget --no-check-certificate  https://ghproxy.com/https://github.com/mstxq17/MoreFind/releases/download/v1.0.2/MoreFind_1.0.2_`uname -s`_`uname -m`.tar.gz
tar -xzvf MoreFind_1.0.2_`uname -s`_`uname -m`.tar.gz
sudo mv ./MoreFind /usr/bin/MoreFind && chmod +x /usr/bin/MoreFind
```

方式三: 本地编译
```bash
git clone https://github.com/mstxq17/MoreFind.git
chmod +x ./build.sh && ./build.sh
```

## 用法说明
1)帮助信息
```bash
MoreFind -h
```
![img.png](README.assets/img.png)



2)导出URL

```bash
MoreFind -u
```

![image-20220613101518150](README.assets/image-20220613101518150.png)



3)导出域名

```bash
MoreFind -d
```

![image-20220613101624590](README.assets/image-20220613101624590.png)



4)导出ip

```bash
MoreFind -i
```

![image-20220613101715993](README.assets/image-20220613101715993.png)



5)支持导出结果

```bash
MoreFind -u -d -i -o result.txt
```



## TODO

1)优化代码逻辑和结构

2)输出结果自动去重复

3)完善脚本异常处理部分

4)加入部分URL智能去重代码

5)完善Log输出处理

6)实现自动强制更新