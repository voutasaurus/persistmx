# persistmx

A tool for checking global DNS servers for errors looking up MX records.

# install

	go get -u https://github.com/voutasaurus/persistmx

# usage

	persistmx -d myemaildomain.com <nameservers.txt 2>output.txt

Where nameservers.txt contains a newline separated list of Domain Name System
nameservers.

# caveats

By default persistmx runs 300 workers concurrently so make sure you aren't
running into "too many files" issues.

To specify a custom number of workers, use the n flag.

	persistmx -d myemaildomain.com -n 10 <nameservers.txt 2>output.txt

persistmx uses the Go DNS resolver so on platforms that Go uses the local
resolver instead persistmx will not query different DNS servers.

# contributions

Welcome
