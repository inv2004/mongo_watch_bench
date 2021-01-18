# bench

Small tool to bench mongo ChangeStream latency

### Help
```
Usage of C:\Users\unknown\AppData\Local\Temp\go-build767003074\b001\exe\bench.exe:
  -r int
        amount of readers processes (read all data without filtering) (default 1)
  -rows int
        number of rows to insert per sender (default 1)
  -s int
        amount of sender processes (default 1)
```

### Run
```bash
./bench.exe -rows 1000 -r 1 -s 1
```

### Output
```
rows = 1000, senders = 1, readers = 1
[sender-0] 1000 in 96.9463ms   /// time to insert 1000 rows, insert runs on 1000ms intervals, so the number should be less
[reader-0] 13.855ms            /// maximum latency per reader
[sender-0] 1000 in 87.9497ms
[reader-0] 13.3393ms
[sender-0] 1000 in 92.0267ms
[reader-0] 13.4515ms
```

### Some Tests
#### DigitalOcean
Intel(R) Xeon(R) Gold 6140 CPU @ 2.30GHz / 64Gb

```bash
./bench -rows 3500 -s 10 -r 1
```
avg latency 238ms

```bash
./bench -rows 3500 -s 10 -r 4
```
avg latency 295ms


### ScaleWay
AMD EPYC 7281 16-Core Processor / 96Gb

```bash
./bench --rows 1350 -s 10 -r 1
```
avg latency 157ms

```bash
./bench --rows 1350 -s 10 -r 10
```
avg latency 159ms

