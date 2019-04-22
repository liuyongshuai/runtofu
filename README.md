# search-broker #

search-broker

depend:
1. libevent
2. thrift.0.10.0 //go context
3. gtest
4. go 1.8.3

## 环境部署

### 1. 安装 gcc 4.8.1
可以参考这个地址[gcc安装说明](https://www.cnblogs.com/codemood/archive/2013/06/01/3113200.html)
默认安装到/usr/local/bin下，可以软连接g++/gcc/c++/cc 到/usr/bin下

### 2. 安装 thrift 0.10.0

#### 2.1 安装bison 

安装2.5.1 及以上版本 [参考安装方式](https://www.cnblogs.com/keitsi/p/5368184.html)
指定安装路径: /usr
#### 2.2 安装boost

要求boost 1.5.3 以上版本,
指定安装路径/usr/local [参考安装方式](https://www.cnblogs.com/keitsi/p/5368184.html)
````
$ tar zxvf boost_1_60_0.tar.gz && cd boost_1_60_0
$ ./bootstrap.sh
$ sudo ./b2 threading=multi address-model=64 variant=release stage install —prefix=/usr/local
````
### 2.3 安装libevent 
安装到/usr/local 目录下 [参考安装方式](https://www.cnblogs.com/keitsi/p/5368184.html)

### 2.4 安装thrift 
````
wget http://archive.apache.org/dist/thrift/0.10.0/thrift-0.10.0.tar.gz
tar zxvf thrift-0.10.0.tar.gz
cd thrift-0.10.0
./configure --prefix=/usr/local --with-boost=/usr/local --with-libevent=/usr/local
````
### 2.5 安装thrift 可能遇到的问题
#### 2.5.1 version `GLIBCXX_3.4.15' not found

解决办法：
````
sudo cp ~/gcc-build-4.8.1/x86_64-unknown-linux-gnu/libstdc++-v3/src/.libs/libstdc++.so.6.0.18 /usr/lib64/
rm -rf /usr/lib64/libstdc++.so.6
sudo ln -s /usr/lib64/libstdc++.so.6.0.18 /usr/lib64/libstdc++.so.6
````
#### 2.5.2 g++: error: /usr/lib64/libboost_unit_test_framework.a: No such file or directory

解决办法：
````
 cp /usr/local/lib/libboost_unit_test_framework.a /usr/lib64/
````
 
其他问题自行google

### 3. 安装google test
````
wget https://github.com/google/googletest/archive/release-1.8.0.tar.gz

tar zxvf release-1.8.0.tar.gz && cd googletest-release-1.8.0/googletest/
````
修改一下CMakeLists.txt 中, BUILD_SHARED_LIBS 的值为 ON
````
$cmake .
$make && sudo make install
````
检查 /usr/local/lib/libgtest.so 是否存在

### 4. 确认git 版本

目前统一使用的是git 1.8.3版本

### 5. 安装go （可选，1.8 或者1.9）



