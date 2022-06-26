# 分布式缓存 
——from《分布式缓存 原理、架构及Go语言实现》

## Dependencies 
- MacOS 12.4
- Golang 1.18

## How to Run 
- brew install boost
- brew install rocksdb
- brew install snappy
- for include 'boost' 'snappy' 'rocksdb' in C++
  - export CPLUS_INCLUDE_PATH=/opt/homebrew/Cellar/boost/1.78.0_1/include:/opt/homebrew/Cellar/rocksdb/7.0.3/include:/opt/homebrew/Cellar/snappy/1.1.9/include

- for '-lsnappy' in rocksdb_performance/makefile
  - my snappy path is /opt/homebrew/Cellar/snappy/1.1.9/
  - sudo ln -s /opt/homebrew/Cellar/snappy/1.1.9/lib/libsnappy.dylib  /usr/local/lib 
  - sudo ln -s /opt/homebrew/Cellar/snappy/1.1.9/include/snappy  /usr/local/include/snappy
- for '-lrocksdb' in rocksdb_performance/makefile
  - ...
  - ...
- for '-lboost_program_options' in rocksdb_performance/makefile 
  - ...
  - ...




