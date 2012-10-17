package main
import (
    "fmt"
    "log"
    "time"
    "github.com/cloudfoundry/gosigar"
    "github.com/jteeuwen/go-pkg-mpd"
)

var err error
var c   *mpd.Client

func resources() string {

   mem    := sigar.Mem{}
   cpu    := sigar.Cpu{}

   mem.Get()
   ramChange := float32(mem.ActualUsed) / float32(mem.ActualFree) * 100

   cpu.Get()
   cpuTotal1 := cpu.User + cpu.Nice + cpu.Sys + cpu.Idle + cpu.Wait + cpu.Irq + cpu.SoftIrq
   cpuWork1  := cpu.User + cpu.Nice + cpu.Sys

   time.Sleep(30 * time.Millisecond) // Minimum time interval for my machine

   cpu.Get()
   cpuTotal2 := cpu.User + cpu.Nice + cpu.Sys + cpu.Idle + cpu.Wait + cpu.Irq + cpu.SoftIrq
   cpuWork2  := cpu.User + cpu.Nice + cpu.Sys
   totalDifference := (cpuTotal2 - cpuTotal1)
   workDifference := (cpuWork2 - cpuWork1)

   cpuChange := float32(0)
   if (workDifference > 0 && totalDifference > 0) {
       cpuChange = (float32(workDifference) / float32(totalDifference)) * 100
   }

   hdd := sigar.FileSystemUsage{}
   hdd.Get("/")
   hddChange := uint32(hdd.UsePercent())

   dateTime := time.Now().Format("Mon, 2 Jan 3:04:05 PM")
   delimiter := "^s[right;#445544; :: ]"

   result := ""
   result = fmt.Sprintf("^s[right;#AABBAA;CPU: ] ^s[right;#E07C00;% 3d%% ]  ^g[right;80;10;%d;100;#445544;#AABBAA;ckycpu] %s ^s[right;#AABBAA;RAM:] ^s[right;#43E000; % 3d%% ] ^p[right;8;10;0;%d;100;#445544;#AABBAA;ckyhdd] %s ^s[right;#AABBAA;HDD:] ^s[right;#008BE0; % 3d%% ] ^p[right;8;10;0;%d;100;#445544;#AABBAA;ckyhdd] %s ^s[right;#AABBAA;%s] %s",
       uint32(cpuChange),
       uint32(cpuChange),
       delimiter,
       uint32(ramChange),
       uint32(ramChange),
       delimiter,
       hddChange,
       hddChange,
       delimiter,
       dateTime,
       delimiter,
   )
   return result
}

func mpdClient() string {
    if current, err := c.Current(); err == nil {
        result := fmt.Sprintf("^s[right;#AABBAA;MPD:] ^s[right;#008BE0; %s - %s - %s] ^s[right;#445544; :: ] ", current["Album"], current["Artist"], current["Title"])
        return result
    } else {
        return ""
        log.Fatal(err)
    }
    return ""
}

func main() {

    if c, err = mpd.Dial("127.0.0.1:6600", ""); err == nil {
        defer c.Close()

        /**
        * Print all the results.
        */

        for {

            mpd := mpdClient()
            res := resources()
            fmt.Printf("^s[right;#445544; :: ] %s %s\n",
                mpd,
                res,
            )
        }
    } else {
        log.Fatal(err)
    }

}
