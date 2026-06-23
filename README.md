# Douyin Downloader (douyin-go)

A high-performance, concurrent command-line utility written in Go to download Douyin videos, cover images, audio tracks (MP3), and photo slideshow albums.

---

## Features

- **Concurrent Downloads**: 
  - **Post-level Concurrency**: Downloads multiple posts in parallel when crawling user feeds or liked lists.
  - **Asset-level Concurrency**: Downloads cover images, audio tracks, and videos/slideshows concurrently for a single post using goroutines.
  - **Slideshow Concurrency Limits**: Limits concurrent image downloads (default max 5) to prevent CDNs from rate-limiting or dropping connections.
- **Cobra CLI Powered**: Structured subcommand interface (`video`, `user`, `likes`) with comprehensive help page documentation.
- **Automatic URL Extraction**: Parses sharing texts and automatically resolves redirection for `v.douyin.com` short links to extract correct IDs.
- **Config Management**: Automatically generates a `config.yaml` configuration file for persistent cookies, user-agents, output directories, and proxy setups.
- **Proxies Support**: Native support for HTTP, HTTPS, and SOCKS5 proxies.
- **Customizable Output Directory**: Organizes downloaded files cleanly with path templates like `{nickname}/{publish_date}_{title}/`.
- **Metadata JSON Saving**: Automatically saves the full raw video/post metadata as a formatted `data.json` file inside the download directory of each item.

---

## Installation

### Prerequisites
- [Go](https://go.dev/doc/install) 1.21 or later.

### Build from Source
Clone the repository and compile the binary:
```bash
git clone https://github.com/thaunq95/douyin-go.git
cd douyin-go
go build -o douyin-go ./cmd/douyin-go/main.go
```

---

## Usage

You can run the compiled binary `douyin-go` with the following commands. Use the `--help` flag on any command to view the options.

### 1. Help Information
```bash
./douyin-go --help
./douyin-go video --help
```

### 2. Download a Single Video or Photo Album
```bash
./douyin-go video "https://v.douyin.com/L4FJNR3/"
```
*Supports raw links, short redirection links, and sharing texts copied from the Douyin mobile app.*

### 3. Download All Published Posts from a User
```bash
./douyin-go user "https://www.douyin.com/user/MS4wLjABAAAA..." --count 50
```
- Use the `--count` (or `-n`) flag to define the retrieval page size.
- The crawler automatically handles pagination recursively until all posts are fetched.

### 4. Download Liked Posts from a User
```bash
./douyin-go likes "https://www.douyin.com/user/MS4wLjABAAAA..." --cookie "YOUR_COOKIE_HERE"
```
*Note: Downloading liked posts or accessing full profile items typically requires a valid logged-in session cookie to bypass Douyin restrictions.*

### 5. Using a Proxy
```bash
./douyin-go video "https://v.douyin.com/L4FJNR3/" --proxy "socks5://127.0.0.1:1080"
```

---

## Configuration (`config.yaml`)

On the first run, the tool automatically generates a `config.yaml` in the root folder. You can configure options here so you don't need to specify them as CLI flags every time:

```yaml
# Paste your Douyin session cookie (headers cookie value) here
cookie: ""
# Default User-Agent matching cryptographic signature requirements
user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome..."
# Default folder path for saved files
output_dir: "./downloads"
# Path structure template
naming_template: "{nickname}/{publish_date}_{title}"
# Fetch limit page size
count: 50
# Global proxy URL no auth (e.g. http://127.0.0.1:8080 or socks5://127.0.0.1:1080)
# Global proxy URL with auth (e.g. http://user:pass@127.0.0.1:8080 or socks5://user:pass@127.0.0.1:1080) 
proxy: ""
# Number of concurrent post downloads allowed
concurrency: 5
```

### Path Template Placeholders
You can customize the `naming_template` in `config.yaml` or via the `--template` flag using these placeholders:
- `{nickname}`: The user's nickname.
- `{publish_date}`: The release date of the post (`YYYY-MM-DD`).
- `{title}`: The post description (truncated to 50 characters).
- `{aweme_id}`: The unique identifier of the post.

---

### Inspired
- [https://github.com/jiji262/douyin-downloader](https://github.com/jiji262/douyin-downloader)
- [https://github.com/Evil0ctal/Douyin_TikTok_Download_API](https://github.com/Evil0ctal/Douyin_TikTok_Download_API)
---

## License

This project is licensed under the MIT License. Feel free to use and contribute!
