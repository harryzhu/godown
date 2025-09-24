package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	taskq     chan string
	taskTotal int
)

func init() {

}

func GetRun() {
	flist := LoadFileList(FileList)
	DebugInfo("count", len(flist))

	taskq = make(chan string, Workers)
	wg := sync.WaitGroup{}
	defer close(taskq)

	var curNum int32
	var idx int32

	timeShellStart := GetTimeNow()
	var progressShowThreshold int = 10
	if len(flist) > 999 {
		progressShowThreshold = 100
	}
	timeRoundStart := GetTimeUnix()
	for k, v := range flist {
		taskq <- v
		if IsDebug || k%progressShowThreshold == 0 || atomic.LoadInt32(&curNum) == int32(len(flist)) {
			mem := GetMemStats()
			timeRoundElapse := GetTimeUnix() - timeRoundStart
			fmt.Printf("GetRun: %v sec => %02d => %v/%v [Mem: Total=%v MB, Sys=%v MB, NumGC=%v]\n",
				timeRoundElapse,
				atomic.LoadInt32(&idx),
				atomic.LoadInt32(&curNum),
				len(flist),
				mem["TotalAlloc"], mem["Sys"], mem["NumGC"],
			)
			timeRoundStart = GetTimeUnix()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&curNum, 1)

			line := <-taskq

			if !strings.Contains(line, "|") {
				return
			}
			segments := strings.Split(line, "|")
			if len(segments) != 3 {
				return
			}

			ftarget := strings.TrimSpace(segments[0])
			furl := strings.TrimSpace(segments[1])
			ftime := strings.TrimSpace(segments[2])

			DownloadFile(furl, ftarget, ftime)

			atomic.AddInt32(&idx, -1)
			if IsDebug {
				fmt.Printf("GetRun: %v / %v: %v \n", curNum, taskTotal, ftarget)
			}
		}()

		atomic.AddInt32(&idx, 1)

		for {
			if atomic.LoadInt32(&idx) >= int32(Workers) {
				time.Sleep(time.Millisecond * 300)
			} else {
				break
			}
		}

	}

	wg.Wait()

	fmt.Printf("GetRun: %02d => %v/%v \n", atomic.LoadInt32(&idx), atomic.LoadInt32(&curNum), len(flist))

	fmt.Println("\n *** ShellRun Elapse:", time.Since(timeShellStart), "***")
	time.Sleep(2 * time.Second)
}

func GetTimeNow() time.Time {
	return time.Now()
}

func GetTimeUnix() int64 {
	return time.Now().Unix()
}

func GetMemStats() map[string]uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	var MB uint64 = 1024 * 1024
	result := make(map[string]uint64)

	result["Alloc"] = mem.Alloc / MB
	result["TotalAlloc"] = mem.TotalAlloc / MB
	result["Sys"] = mem.Sys / MB
	result["HeapAlloc"] = mem.HeapAlloc / MB
	result["HeapSys"] = mem.HeapSys / MB
	result["NumGC"] = uint64(mem.NumGC)
	return result
}

func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func Str2Int64(n string) int64 {
	s, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return -1
	}
	return s
}

func ParseHeader() (headers map[string]string) {
	headers = make(map[string]string)
	if len(Header) > 0 {
		for _, kv := range Header {
			idxColon := strings.Index(kv, ":")
			if idxColon > 0 {
				k := strings.TrimSpace(kv[:idxColon])
				v := strings.TrimSpace(kv[idxColon+1:])
				if k != "" {
					headers[k] = v
				}
			}
		}
	}
	return headers
}

func LoadFileList(fpath string) (flist []string) {
	fp, err := os.Open(fpath)
	FatalError("LoadFileList", err)
	fcontent, err := io.ReadAll(fp)
	lines := strings.Split(string(fcontent), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			flist = append(flist, line)
		}
	}
	return flist
}

func DownloadFile(furl, ftarget, mtime string) error {
	furl = strings.TrimSpace(furl)
	ftarget = strings.TrimSpace(ftarget)

	_, err := os.Stat(ftarget)
	if err == nil {
		DebugInfo("DownloadFile", "SKIP")
		return nil
	}

	downloadClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", furl, nil)
	req.Header.Set("User-Agent", UserAgent)
	if len(kvHeaders) > 0 {
		for k, v := range kvHeaders {
			req.Header.Set(k, v)
		}
	}
	req.Close = true

	resp, err := downloadClient.Do(req)
	if err != nil {
		PrintError("DownloadFile", err)
		return err
	}
	defer resp.Body.Close()

	fdir := filepath.Dir(ftarget)
	finfo, err := os.Stat(fdir)
	if err != nil || finfo.IsDir() == false {
		os.MkdirAll(fdir, os.ModePerm)
	}

	ftemp := strings.Join([]string{ftarget, "ing"}, ".")
	out, err := os.OpenFile(ftemp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	PrintError("DownloadFile:os.Create(ftemp)", err)

	_, errCopy := io.Copy(out, resp.Body)
	out.Close()

	tempInfo, err := os.Stat(ftemp)
	if err != nil {
		PrintError("DownloadFile:os.Stat(ftemp)", err)
		return err
	}

	if IsPurgeErrorFile == true && errCopy != nil {
		os.Remove(ftemp)
		return errCopy
	}

	if MinSize != 0 {
		if tempInfo.Size() <= MinSize {
			os.Remove(ftemp)
			return nil
		}
	}

	err = os.Rename(ftemp, ftarget)
	if err != nil {
		PrintError("DownloadFile:os.Rename", err)
		return err
	}

	err = os.Chtimes(ftarget, SecondToTime(Str2Int64(mtime)), SecondToTime(Str2Int64(mtime)))

	PrintError("DownloadFile:os.Chtimes", err)

	return err
}

func ShellCommand(fcmd string) error {
	cmdBash := exec.Command("bash", "-c", fcmd)
	if runtime.GOOS == "windows" {
		cmdBash = exec.Command("cmd", "/Q", "/C", fcmd)
	}
	if IsDebug {
		out, err := cmdBash.CombinedOutput()
		if err != nil {
			PrintError("ShellCommand.10", err)
		}
		DebugInfo("ShellCommand.20", string(out))
	} else {
		_, err := cmdBash.CombinedOutput()
		if err != nil {
			PrintError("ShellCommand.30", err)
		}
	}
	return nil
}

func ShellRun() {
	flist := LoadFileList(CmdList)
	DebugInfo("count", len(flist))

	taskq = make(chan string, Workers)
	wg := sync.WaitGroup{}
	defer close(taskq)

	var curNum int32
	var idx int32

	timeShellStart := GetTimeNow()
	var progressShowThreshold int = 10
	if len(flist) > 999 {
		progressShowThreshold = 100
	}
	timeRoundStart := GetTimeUnix()
	for k, v := range flist {
		taskq <- v
		if IsDebug || k%progressShowThreshold == 0 || atomic.LoadInt32(&curNum) == int32(len(flist)) {
			mem := GetMemStats()
			timeRoundElapse := GetTimeUnix() - timeRoundStart
			fmt.Printf("ShellRun: %v sec => %02d => %v/%v [Mem: Total=%v MB, Sys=%v MB, NumGC=%v]\n",
				timeRoundElapse,
				atomic.LoadInt32(&idx),
				atomic.LoadInt32(&curNum),
				len(flist),
				mem["TotalAlloc"], mem["Sys"], mem["NumGC"],
			)
			timeRoundStart = GetTimeUnix()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&curNum, 1)

			line := <-taskq
			ShellCommand(line)

			atomic.AddInt32(&idx, -1)
		}()

		atomic.AddInt32(&idx, 1)

		for {
			if atomic.LoadInt32(&idx) >= int32(Workers) {
				time.Sleep(time.Millisecond * 300)
			} else {
				break
			}
		}
	}

	wg.Wait()

	fmt.Printf("ShellRun: %02d => %v/%v \n", atomic.LoadInt32(&idx), atomic.LoadInt32(&curNum), len(flist))

	fmt.Println("\n *** ShellRun Elapse:", time.Since(timeShellStart), "***")
	time.Sleep(2 * time.Second)
}
