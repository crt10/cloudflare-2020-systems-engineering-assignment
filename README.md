# Cloudflare Workers - 2020 Systems Engineering Assignment

CLI tool to make, time and measure GET HTTP requests, done in golang.

## Build
- Execute `make` in the root directory

## Usage
- Execute `html-request`

### Arguments
- `help` Lists possible arguments
- `--url`	string	URL to perform GET HTML request (ex: http://example.com/)
- `--profile` int Number of times to request from the URL

## Screenshots

### Cloudflare Workers
- `https://link-tree.tennyson-cheng.workers.dev/links`
![](/1.PNG)

### Apache
- `http://apache.org/`
![](/2.PNG)
![](/3.PNG)

### MIT
- `http://web.mit.edu/`
![](/4.PNG)

## Comparison

Dividing each response size with their respective mean response time:
- `http://apache.org` = 0.0025 ms/byte
- `http://web.mit.edu/` = 0.0008 ms/byte
- `https://link-tree.tennyson-cheng.workers.dev/links` = 0.2 ms/byte
Interestingly, Cloudflare Workers had a significantly longer response time (about 150x longer).