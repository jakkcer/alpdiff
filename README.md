# alpdiff
show differences between two alp results

## Installation
### Prerequesties
(alp)[https://github.com/tkuchiki/alp] needed for alpdiff

### Binary distribution
You can pick your download (here)[https://github.com/jakkcer/alpdiff/releases], and install it as follows:
```sh
sudo install <download file> /usr/local/bin/alpdiff
```

## Usage
```sh
$ alpdiff --help
Usage: alpdiff <old log file> <new log file> [<flags>]

Flags:
  -m string
      same as alp -m option
```

### Example
```sh
$ ./alpdiff example_log/ltsv_access.log example_log/new_ltsv_access.log -m "/diary/entry/\d+"
+-------+-----+-----+-----+-----+-----+--------+------------------+---------+---------+---------+---------+---------+---------+---------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |       URI        |   MIN   |   MAX   |   SUM   |   AVG   |   P90   |   P95   |   P99   |
+-------+-----+-----+-----+-----+-----+--------+------------------+---------+---------+---------+---------+---------+---------+---------+
|     0 |   0 |   0 |   0 |   0 |   0 | POST   | /hoge/piyo       |   0.700 |   0.700 |   0.700 |   0.700 |   0.700 |   0.700 |   0.700 |
|     0 |   0 |   0 |   0 |   0 |   0 | GET    | /foo/bar/5xx     | -30.000 | -30.000 | -30.000 | -30.000 | -30.000 | -30.000 | -30.000 |
|     0 |   0 |   0 |   0 |   0 |   0 | GET    | /req             |  -0.200 |  -0.200 |  -0.200 |  -0.200 |  -0.200 |  -0.200 |  -0.200 |
|     0 |   0 |   0 |   0 |   0 |   0 | GET    | /foo/bar         |   0.100 |   1.000 |   1.100 |   0.550 |   1.000 |   1.000 |   1.000 |
|     0 |   0 |   0 |   0 |   0 |   0 | GET    | /diary/entry/\d+ |   0.397 |   0.103 |   0.500 |   0.251 |   0.103 |   0.103 |   0.103 |
|     0 |   0 |   0 |   0 |   0 |   0 | POST   | /foo/bar         |  -0.047 |   0.300 |   0.309 |   0.061 |   0.300 |   0.300 |   0.300 |
+-------+-----+-----+-----+-----+-----+--------+------------------+---------+---------+---------+---------+---------+---------+---------+
```
