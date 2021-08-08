# chat_notify

Hi, this is telegram bot announcement program.

# How does it work?

1. Edit robot configuration file
    * Under the ```config``` directory
    * The file name ends with .toml, eg: xxx.toml

```toml
[Telegram]
# Telegram Bot Key Input Here
BotKey = "1779011813:AAHKeFFD5R9mAg8pRfSKXfFkg05Mbo0NS3Q"

# Bot Program DB Connect Info
[MysqlDB]
Host = "127.0.0.1"
Port = 3306
Username = "root"
Password = "root"
DBName = "strangerbot"
# The following configuration can be left unchanged
Charset = "utf8"
MaxOpenConns = 1000
MaxIdleConns = 1000
ConnMaxLifetime = 10
```

2. Write Announcement Content
    * Under the ```notify``` directory
    * The file name is arbitrary, it is recommended to use the date record.
    * File does not need to extend the name


3. Run the program

* cfg: No need to enter the file extension
* notify: File needs no extension

```
./chat_notify -cfg=robot1 -notify=20210804
```

console will output:

```
2021/08/08 21:02:03 config file:  robot1 notify file:  20210804
2021/08/08 21:02:03 your send notify content: 
2021/08/08 21:02:03 Hi

This is Test message.
2021/08/08 21:02:03 chat count:  1
2021/08/08 21:02:05 start run sender...
2021/08/08 21:02:05 send offset:  0 limit:  100
2021/08/08 21:02:05 send offset:  100 limit:  100
2021/08/08 21:02:05 send total:  1

Debugger finished with the exit code 0
```