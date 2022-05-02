## How it works

Can use any of the arguments below, but "**-IP**" is a required argument.

**Example #1**: `ping_wrapper-win.exe -IP 1.1.1.1 -L 100 -I 2 -C 10 -O Availability` will return `1` - available or `0' - not available 

**Example #2**: `ping_wrapper-win.exe -IP 1.1.1.1 -L 100 -I 2 -C 10 -R true` will return `RAW output -> {10 10 0 0 1.1.1.1 1.1.1.1 [31.5037ms 32.0743ms 31.6363ms 31.6239ms 31.6239ms 30.6653ms 31.356ms 31.4617ms 31.351ms 31.2691ms] 30.6653ms 32.0743ms 31.456522ms 339.802µs}`

### List of possible arguments
```shell
Usage of ping_wrapper (Version: 0.3.0, build info: go1.16.6 [2022-05-02 11:02:28AM UTC]), [<input args>] [<output args>]
  Input args:
   -IP string
         IPv4 address or Domain name. Example: 8.8.8.8 or dns.google (default "dns.google")
         Number of echo requests to send (default 5)
   -L string
         Size of packet being sent (default 24)
   -I string
         Interval is the wait time between each packet send. Default is 1s. (default 1)
   -T string
         Timeout specifies a timeout before ping exits, regardless of how many packets have been receive (default -1)
   -S string
         Source is the source IP address (default "")
   -L string
         Size of packet being sent (default 24)
   -I string
         Interval is the wait time between each packet send. Default is 1s. (default 1)
   -T string
         Timeout specifies a timeout before ping exits, regardless of how many packets have been receive (default -1)
   -S string
         Source is the source IP address (default "")
  Output args:
   -O string
         Returns one argument from the returned result (default "")
         Can use one of the keys: Availability, Min, Max, Std, Loss, Sent, Recv
   -R string
         Display RAW result. Incompatible with argument 'O' (default false)

```

## Attention 

> **Output args 'O' and 'R' don't work together. Only 'O' or 'R'**

## License

See the [LICENSE](LICENSE) file for license rights and limitations (MIT).

