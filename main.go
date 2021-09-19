package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-ping/ping"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const version string = "0.1.0"

type Flags struct {
	name  string
	value string
	usage string
}

var (
	// input args
	IP = Flags{"IP", "dns.google",
		"IPv4 address or Domain name. Example: 8.8.8.8 or dns.google"}
	IPv4 = flag.String(IP.name, IP.value, IP.usage)

	CountF    = Flags{"C", "5", "Number of echo requests to send"}
	CountV, _ = strconv.Atoi(CountF.value)
	Count     = flag.Int(CountF.name, CountV, CountF.usage)

	IntervalF    = Flags{"I", "1", "Interval is the wait time between each packet send. Default is 1s."}
	IntervalV, _ = strconv.Atoi(IntervalF.value)
	Interval     = flag.Int(IntervalF.name, IntervalV, IntervalF.usage)

	TimeoutF = Flags{"T", "-1",
		"Timeout specifies a timeout before ping exits, regardless of how many packets have been receive"}
	TimeoutV, _ = strconv.Atoi(TimeoutF.value)
	Timeout     = flag.Int(TimeoutF.name, TimeoutV, TimeoutF.usage)

	PacketSizeF    = Flags{"L", "24", "Size of packet being sent"}
	PacketSizeV, _ = strconv.Atoi(PacketSizeF.value)
	PacketSize     = flag.Int(PacketSizeF.name, PacketSizeV, PacketSizeF.usage)

	SourceF = Flags{"S", "", "Source is the source IP address"}
	Source  = flag.String(SourceF.name, SourceF.value, SourceF.usage)

	// output args
	ArgsF = Flags{"O", "", "Returns one argument from the returned result"}
	Args  = flag.String(ArgsF.name, ArgsF.value, ArgsF.usage)

	RawF    = Flags{"R", "false", "Display RAW result. Incompatible with argument 'O'"}
	RawV, _ = strconv.ParseBool(RawF.value)
	Raw     = flag.Bool(RawF.name, RawV, RawF.usage)
)

func PrintDefaults() {
	spacesl1 := strings.Repeat(" ", 2)
	spacesl2 := spacesl1 + strings.Repeat(" ", 1)
	spacesl3 := spacesl2 + strings.Repeat(" ", 6)
	fmt.Printf(spacesl1 + "Input args:\n")
	fmt.Printf(spacesl2+"-%s %T\n%s%s", IP.name, IP.name, spacesl3, IP.usage)

}

var Usage = func() {
	fmt.Printf("Usage of ping_wrapper (ver %s) [<input args>] [<output args>]\n", version)
	PrintDefaults()
}

func main() {
	flag.Usage = Usage
	// Init variables

	//var PacketsSent = flag.Int()
	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
	} else if *Args != "" && *Raw != false {
		fmt.Println("Output parameters can't have RAW and processed format at the same time")
	} else {
		pinger, err := ping.NewPinger(*IPv4)
		if err != nil {
			panic(err)
		}
		// Assign values from command line instead of defaults
		pinger.Count = *Count
		if *Timeout != -1 {
			pinger.Timeout = (time.Duration)(*Timeout)
		}
		if *Interval != 1 {
			pinger.Interval = (time.Duration)(*Interval)
		}
		if *PacketSize != 24 {
			pinger.Size = *PacketSize
		}
		// добавить валидацию на IP
		// добавить параметры выхода
		if *Source != "" {
			result, err := ipv4Validator(*Source)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if (err == nil) && (result == "") {
				os.Exit(1)
			} else {
				if result == "VALID" {
					pinger.Source = *Source
				}
			}
		}
		// change UDP proto to ICMP
		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		err = pinger.Run()
		if err != nil {
			panic(err)
		}
		if *Args == "Avg" {
			fmt.Println(pinger.Statistics().AvgRtt)
		} else if *Args == "Min" {
			fmt.Println(pinger.Statistics().MinRtt)
		} else if *Args == "Max" {
			fmt.Println(pinger.Statistics().MaxRtt)
		} else if *Args == "Std" {
			fmt.Println(pinger.Statistics().StdDevRtt)
		} else if *Args == "Loss" {
			fmt.Println(pinger.Statistics().PacketLoss)
		} else if *Args == "Sent" {
			fmt.Println(pinger.Statistics().PacketsSent)
		} else if *Args == "Recv" {
			fmt.Println(pinger.Statistics().PacketsRecv)
		} else if *Args == "Availability" {
			if pinger.Statistics().PacketLoss == 0 {
				fmt.Println("1")
			} else if pinger.Statistics().PacketLoss == 100 {
				fmt.Println("0")
			}
		} else {
			fmt.Println("Sorry, but specified parameter is unknown")
		}

	}

}
func ipv4Validator(s string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered in ipv4Validator: ", err)
		}
	}()
	tmp := strings.Split(s, ".")
	if len(tmp) != 4 {
		return "", errors.New("Incorrect IPv4 address")
	} else {
		for i := 0; i <= 3; i++ {
			octet, err := strconv.Atoi(tmp[i])
			if err != nil {
				panic(err)
			} else {
				if octet > 255 {
					return "", errors.New(fmt.Sprintf("Incorrect IPv4 address: octet (%d) can't be more than 255", i+1))
				}
			}
		}
	}
	return "VALID", nil
}
