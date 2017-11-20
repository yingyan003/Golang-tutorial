fileDemo:文件与目录的相关操作
==========================

## 包：os

### 目录操作
------------------

#### 获取系统用于保存临时文件的默认目录

```
func TempDir() string

```

如在我的mac中是：/var/folders/v_/66t4cmys0sg90gdkll8wbvpm0000gn/T/


#### 获取当前工作目录的根路径

```
func Getwd() (dir string, err error)

```

1. 当用相对路径表示时，（当前工作路径）"."表示的与函数返回的路径一致。

2. 如果用开发工具（如Goland-EAP）直接运行时，返回的是项目的根目录。<br>
eg. 项目名为web，全路径为/usr/web,返回的则是/usr/web

3. 如果用（shell）命令行方式运行，返回的是当前（shell）的工作目录，即linux下pwd命令返回的路径。<br>
eg. host:fileDemo host$ 命令。该例子模拟shell的命令行，host是你的主机，fileDemo是你当前的工作目录，
pwd返回的则是该目录的全路径。


#### 判断文件或目录是否存在

```
func Stat(name string) (fi FileInfo, err error)   or <br>
func Open(name string) (file *File, err error)
<br> + <br>
func IsNotExist(err error) bool  or  <br>
func IsExist(err error) bool

```

Stat：返回一个描述name（可以是文件或目录的路径）指定文件对象的FileInfo。如果指定的文件对象是一个
符号链接（有效的路径），返回的FileInfo描述该符号链接的文件信息。如果出错，返回错误值为*PathError类型<br><br>

Open：根据name指定的路径（可以是文件或目录）打开一个文件用于读取。如果操作成功，返回文件对象的方法可用于读取数据；<br>
对应的文件描述符具有O_RDONOLY(只读)模式。如果出错，错误底层类型是*PathError.<br><br>

IsNotExist/IsExist：返回一个bool值说明该错误是否表示一个文件或目录不存在/已经存在，参数可以是上面Stat/Open返回的err。


#### 创建目录

```
func Mkdir(name string, perm FileMode) error    or  <br><br>

func MkdirAll(path string,perm FileMode) error

```

Mkdir：创建单级目录，使用指定的名称和权限。若果出错，会返回*PathError类型错误。名称可以是目录的绝对或相对路径，但只允许
创建单级。<br>
eg. 如/usr/a/b，要创建的目录是b，则a必须存在，也就是说/usr/a是可寻址的有效路径。<br><br>

Mkdir：创建多级目录，使用指定的名称和权限。目录（包括任何上级目录）不存在则创建，存在则忽略，并返回`nil`。否则返回错误。
权限位perm会作用在每一个被创建的目录上


#### 删除文件或目录

```
func Remove(name string) error  or  <br><br>

func RemoveAll(path string) error

```

Remove：删除name指定的单个文件或`空`目录。如果出错（如指定路径不存在，目录不为空），返回错误。<br><br>

RemoveAll：删除path指定的文件或目录（包括目录的任何下级对象）。指定的路径不存在时，返回`nil`。它会尝试删除所有东西，
除非遇到错误并返回。


### 文件操作
------------------


#### 判断文件是否存在

参考上文：判断文件或目录是否存在

#### 创建文件

```
func Create(name string) (file *File, err error) <br><br>

func OpenFile(name string, flag int, perm FileMode) (file *File, err error)

```

Create：采用模式0666(任何人可读写，不可执行)创建一个名为name的文件，如果name指定路径的文件存在，
file返回已存在文件的指针，err返回`nil`。成功则返回可用于I/O的文件对象指针。对应的文件描述符具有O_RDWR模式。否则返回错误。<br><br>

OpenFile：更一般性的文件打开函数，常用Create或Open代替本函数。
使用指定的选项（如O_RDONLU等），指定的模式（如0666等）打开指定名称的文件。如果文件不存在，可通过指定flag=os.O_CREATE创建，
如果文件存在则以指定flag（os.O_RDONLY只读，os.O_RDWR可读写）的形式打开。


#### 打开文件

```
func Open(name string) (file *File, err error)   <br><br>

func OpenFile(name string, flag int, perm FileMode) (file *File, err error)

```

Open：根据name指定的路径（可以是文件或目录）打开一个文件用于读取。如果操作成功，返回文件对象的方法可用于读取数据；<br>
对应的文件描述符具有O_RDONOLY(只读)模式。如果出错（如指定文件路径不存在），错误底层类型是*PathError.<br><br>

OpenFile：参考上文


#### 删除文件

参考上文：删除文件或目录




## 包：ioutil

### 目录操作
------------------

#### 用指定前缀建立临时文件夹

```
func TempDir(dir, prefix string) (name string, err error)

```

在dir目录里创建一个新的，用prefix为前缀的临时文件夹，并返回文件夹路径（文件夹名字是指定前缀+系统生成的一串数字）。
如果dir是空字符串，TempDir使用默认用于临时文件的目录（os.TempDir返回的目录）。不同程序同时调用该函数会创建不同的临时目录，
调用本函数的程序有责任在不需要临时文件夹时摧毁它。


#### 用指定前缀建立临时文件

```
func TempFile(dir, prefix string) (name string, err error)

```

在dir目录里创建一个新的，用prefix为前缀的临时文件，以读写模式打开该文件并返回os.File指针。如果dir是空字符串，
TempDir使用默认用于临时文件的目录（os.TempDir返回的目录）。不同程序同时调用该函数会创建不同的临时文件，
调用本函数的程序有责任在不需要临时文件时摧毁它。
