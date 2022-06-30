# MoreFind
一款用于快速导出URL、Domain和IP的小工具


## 快速安装
方式一: 通过Go包管理安装
```bash
go install  github.com/mstxq17/MoreFind@latest
```
方式二: 直接安装二进制文件
```bash
wget --no-check-certificate  https://ghproxy.com/https://github.com/mstxq17/MoreFind/releases/download/v1.2.0/MoreFind_1.2.0_`uname -s`_`uname -m`.tar.gz
tar -xzvf MoreFind_1.2.0_`uname -s`_`uname -m`.tar.gz
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
```bash
MoreFind is a very fast script for searching URL、Domain and Ip from specified stream

Usage:
  morefind [flags]
  morefind [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the semantic version number of MoreFind

Flags:
  -d, --domain          search domain from stdin or file
  -e, --exclude         exclude internal/private segment of ip when searching ip
  -f, --file string     search the info in specified file
  -h, --help            help for morefind
  -i, --ip              search ip from stdin or file
  -l, --len string      search specify the length of string, "-l 35" == "-l 0-35" 
  -o, --output string   output the result to specified file
  -s, --show            show the length of each line and summaries
  -u, --url             search url from stdin or file

Use "morefind [command] --help" for more information about a command.

```



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
# 默认会搜索全部ipv4地址
MoreFind -i
# 加上--exclude 或者 -e 会排除属于内网的ip
MoreFind -i -e
```

![image-20220613101715993](README.assets/image-20220613101715993.png)

5)输出统计信息

```bash
MoreFind -s
```

6)筛选指定长度字符串

```bash
MoreFind -l 35 
MoreFind -l 0-35
```

7)支持导出结果

```bash
MoreFind -u -d -i -o result.txt
```



8)联动使用

```bash
echo -e 'baidu.com ccccxxxx 1.com'|MoreFind -d |MoreFind -l 5  
```



## TODO

- [x] 输出结果自动去重复

- [x] 搜索ip的时候支持排除私有IP地址

- [x] 读取文件流，输出统计信息，显示每行长度

- [x] 可指定每行长度筛选出符合条件的字符串

- [ ] 完善脚本异常处理部分

- [ ] 加入部分URL智能去重代码

- [ ] 完善Log的输出和处理

- [ ] 实现自动强制更新

- [ ] 优化代码逻辑、结构和提高执行速度