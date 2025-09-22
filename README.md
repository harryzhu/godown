# godown
* download files or execute shell commands concurrently
* 并行下载文件、并行执行shell命令，可以手动指定同时运行的任务数量，下载任务默认为同时下载3个，shell命令默认为同时执行CPU核数个任务。
* 适用于命令彼此独立、互不依赖场景下的并行执行，方便其他程序调用执行。
* 此工具目前正在被作者用于日常的文件批量下载、使用 `ffmpeg` 将视频转为图片。

## 并行下载文件：
1）用 `python` 生成需要下载文件的清单，格式： 三段信息，使用 | 分隔（英文竖线，前后有空格）：

```
文件下载的URL | 本地保存路径 | 文件的unix时间戳
```

2）将上述格式的下载文件信息，每行一个，保存在文件 `filelist.txt` 中，例如 1000 行 就表示有 1000 个文件需要下载。

3）运行命令：
```
./godown get
```
即可开始并行下载，
* `--workers=5` 可指定同时下载文件的任务数为`5`，默认为`3`；
* `--debug` 可以显示运行时信息；
* `--user-agent` 可以指定发出request的UA；
* `--header` 可以添加发出request的头部信息，支持多个；


## 并行执行 `shell` 命令
1）用 `python` 生成你需要执行命令的清单,保存至 `cmdlist.txt`，格式： 一行一个独立任务：
比如用 `ffprobe` 来检查视频文件是否有效（能获取到视频帧数则认为有效），并输出结果（如果任务无需收集结果则不需要导出结果）：
```
/Users/harryzhu/ffmpeg/ffprobe -v error -select_streams v:0 -show_entries stream=nb_frames -of default=nokey=1:noprint_wrappers=1  "/Users/harryzhu/1.mp4">/Users/harryzhu/1.mp4.txt
```

2）运行命令：
```
./godown shell
```
例如文件`cmdlist.txt`中有1000行这样的命令，CPU是16核心，那么 `godown` 就会默认以16个任务来并行检测这1000个视频是否都有效，并把结果写入你指定的文件 `>/Users/harryzhu/*.txt`。注意，不同命令的输出应该写入各自独立的结果输出文件，因为命令是并行的，如果所有结果写人同一个文件，会导致不可预知的错误。
* `--workers=20` 可指定同时下载文件的任务数为`20`；
* `--debug` 可以显示运行时信息；

3）等待命令运行完成，使用其他工具对每个结果进行收集、分析。

