## REdis Serialization Protocal (RESP)

### 正常回复

* 以 "+" 开头，以 "\r\n" 结尾的字符串形式

```sh
+OK\r\n
```

### 错误回复

* 以 "-" 开头，以 "\r\n" 结尾的字符串形式

```sh
-Error message\r\n
```

### 整数

* 以 ":" 开头，以 "\r\n" 结尾的字符串形式

```sh
:123456\r\n
```

### 多行字符串sh

* 以 "$" 开头，后跟实际发送字节数，以 "\r\n" 结尾

```sh
# "redisx"
$6\r\nredisx\r\n

# ""
$0\r\n\r\n

# "hello\r\nworld"
$10\r\nhello\r\nworld\r\n
```

### 数组

* 以 "*" 开头，后跟成员个数

```sh
# SET key value
*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
```

