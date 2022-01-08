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

const version string = "0.2.0"

type Flags struct {
	name  string
	value string
	usage string
}

var (
	myFlag = flag.NewFlagSet("myFlagSet", 1)
	// input args
	IP = Flags{"IP", "dns.google",
		"IPv4 address or Domain name. Example: 8.8.8.8 or dns.google"}
	IPv4 = myFlag.String(IP.name, IP.value, IP.usage)

	CountF    = Flags{"C", "5", "Number of echo requests to send"}
	CountV, _ = strconv.Atoi(CountF.value)
	Count     = myFlag.Int(CountF.name, CountV, CountF.usage)

	IntervalF    = Flags{"I", "1", "Interval is the wait time between each packet send. Default is 1s."}
	IntervalV, _ = strconv.Atoi(IntervalF.value)
	Interval     = myFlag.Int(IntervalF.name, IntervalV, IntervalF.usage)

	TimeoutF = Flags{"T", "-1",
		"Timeout specifies a timeout before ping exits, regardless of how many packets have been receive"}
	TimeoutV, _ = strconv.Atoi(TimeoutF.value)
	Timeout     = myFlag.Int(TimeoutF.name, TimeoutV, TimeoutF.usage)

	PacketSizeF    = Flags{"L", "24", "Size of packet being sent"}
	PacketSizeV, _ = strconv.Atoi(PacketSizeF.value)
	PacketSize     = myFlag.Int(PacketSizeF.name, PacketSizeV, PacketSizeF.usage)

	SourceF = Flags{"S", "", "Source is the source IP address"}
	Source  = myFlag.String(SourceF.name, SourceF.value, SourceF.usage)
	// output args
	OutputF = Flags{"O", "", "Returns one argument from the returned result"}
	Output  = myFlag.String(OutputF.name, OutputF.value, OutputF.usage)

	RawF = Flags{"R", "false", "Display RAW result. Incompatible with argument 'O'"}
	Raw  = myFlag.String(RawF.name, RawF.value, RawF.usage)

	OutArgs = []string{"Availability", "Min", "Max", "Std", "Loss", "Sent", "Recv"}
)

// Overriding default help
type myFlagSet struct {
	*flag.FlagSet
}

func (f *myFlagSet) PrintDefaults() {
	PrintDefaults()
}

func PrintDefaults() {
	spacesl1 := strings.Repeat(" ", 2)
	spacesl2 := spacesl1 + strings.Repeat(" ", 1)
	spacesl3 := spacesl2 + strings.Repeat(" ", 6)
	fmt.Printf(spacesl1 + "Input args:\n")
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default \"%s\")\n", IP.name, IP.name, spacesl3, IP.usage, IP.value)
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default %s)\n", CountF.name, CountF.name, spacesl3, CountF.usage, CountF.value)
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default %s)\n", PacketSizeF.name, PacketSizeF.name, spacesl3, PacketSizeF.usage, PacketSizeF.value)
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default %s)\n", IntervalF.name, IntervalF.name, spacesl3, IntervalF.usage, IntervalF.value)
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default %s)\n", TimeoutF.name, TimeoutF.name, spacesl3, TimeoutF.usage, TimeoutF.value)
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default \"%s\")\n", SourceF.name, SourceF.name, spacesl3, SourceF.usage, SourceF.value)
	fmt.Printf(spacesl1 + "Output args:\n")
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default \"%s\")\n", OutputF.name, OutputF.name, spacesl3, OutputF.usage, OutputF.value)
	fmt.Printf(spacesl3+"Can use one of the keys: %v\n", strings.Join(OutArgs, ", "))
	fmt.Printf(spacesl2+"-%s %T\n%s%s (default %s)\n", RawF.name, RawF.name, spacesl3, RawF.usage, RawF.value)
}

var Usage = func() {
	fmt.Printf("Usage of ping_wrapper (ver %s) [<input args>] [<output args>]\n", version)
	PrintDefaults()
}

func main() {

	var _flag = myFlagSet{myFlag}
	_flag.Usage = Usage
	err := _flag.Parse(os.Args[1:])

	if err != nil {
		//panic(err)
	}

	if len(os.Args) <= 1 {
		fmt.Printf("Usage of ping_wrapper (ver %s) [<input args>] [<output args>]\n", version)
		PrintDefaults()
	} else if *Output != "" && *Raw != "false" {
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

		if len(*Output) >= 1 {
			var correctRegistry = strings.ToLower(*Output)
			correctRegistry = strings.ToUpper(string(correctRegistry[0])) + correctRegistry[1:]

			ok := false
			for _, v := range OutArgs {
				if correctRegistry == v {
					ok = true
				}
			}
			if ok != true {
				errors.New(fmt.Sprintf("Incorrect argunemt"))
				os.Exit(1)
			}

			if correctRegistry == "Avg" {
				fmt.Println(pinger.Statistics().AvgRtt)
			} else if correctRegistry == "Min" {
				fmt.Println(pinger.Statistics().MinRtt)
			} else if correctRegistry == "Max" {
				fmt.Println(pinger.Statistics().MaxRtt)
			} else if correctRegistry == "Std" {
				fmt.Println(pinger.Statistics().StdDevRtt)
			} else if correctRegistry == "Loss" {
				fmt.Println(pinger.Statistics().PacketLoss)
			} else if correctRegistry == "Sent" {
				fmt.Println(pinger.Statistics().PacketsSent)
			} else if correctRegistry == "Recv" {
				fmt.Println(pinger.Statistics().PacketsRecv)
			} else if correctRegistry == "Availability" {
				if pinger.Statistics().PacketLoss == 0 {
					fmt.Println("1")
				} else if pinger.Statistics().PacketLoss == 100 {
					fmt.Println("0")
				}
			} else {
				fmt.Println("Sorry, but specified parameter is unknown")
			}
		} else if len(*Output) == 0 && *Raw == "true" {
			fmt.Printf("RAW output -> %v", *pinger.Statistics())
		} else if len(*Output) == 0 && *Raw != "true" {
			fmt.Println("Specify an output parameter '-O' or '-R'!")
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
