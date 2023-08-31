package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var stdImports = []string{
	"archive",
	"archive/tar",
	"archive/zip",
	"bufio",
	"builtin",
	"bytes",
	"compress",
	"compress/bzip2",
	"compress/flate",
	"compress/gzip",
	"compress/lzw",
	"compress/zlib",
	"container",
	"container/heap",
	"container/list",
	"container/ring",
	"context",
	"crypto",
	"crypto/aes",
	"crypto/cipher",
	"crypto/des",
	"crypto/dsa",
	"crypto/ecdsa",
	"crypto/ed25519",
	"crypto/elliptic",
	"crypto/hmac",
	"crypto/md5",
	"crypto/rc4",
	"crypto/rsa",
	"crypto/sha1",
	"crypto/sha256",
	"crypto/sha512",
	"crypto/subtle",
	"crypto/tls",
	"crypto/x509",
	"crypto/x509/pkix",
	"database",
	"database/sql",
	"database/sql/driver",
	"debug",
	"debug/dwarf",
	"debug/elf",
	"debug/gosym",
	"debug/macho",
	"debug/pe",
	"debug/plan9obj",
	"encoding",
	"encoding/ascii85",
	"encoding/asn1",
	"encoding/base32",
	"encoding/base64",
	"encoding/binary",
	"encoding/csv",
	"encoding/gob",
	"encoding/hex",
	"encoding/json",
	"encoding/pem",
	"encoding/xml",
	"errors",
	"expvar",
	"flag",
	"fmt",
	"go",
	"go/ast",
	"go/build",
	"go/constant",
	"go/doc",
	"go/format",
	"go/importer",
	"go/parser",
	"go/printer",
	"go/scanner",
	"go/token",
	"go/types",
	"hash",
	"hash/adler32",
	"hash/crc32",
	"hash/crc64",
	"hash/fnv",
	"hash/maphash",
	"html",
	"html/template",
	"image",
	"image/color",
	"image/color/palette",
	"image/draw",
	"image/gif",
	"image/jpeg",
	"image/png",
	"index",
	"index/suffixarray",
	"io",
	"io/ioutil",
	"log",
	"log/syslog",
	"math",
	"math/big",
	"math/bits",
	"math/cmplx",
	"math/rand",
	"mime",
	"mime/multipart",
	"mime/quotedprintable",
	"net",
	"net/http",
	"net/http/cgi",
	"net/http/cookiejar",
	"net/http/fcgi",
	"net/http/httptest",
	"net/http/httptrace",
	"net/http/httputil",
	"net/http/pprof",
	"net/mail",
	"net/rpc",
	"net/rpc/jsonrpc",
	"net/smtp",
	"net/textproto",
	"net/url",
	"os",
	"os/exec",
	"os/signal",
	"os/user",
	"path",
	"path/filepath",
	"plugin",
	"reflect",
	"regexp",
	"regexp/syntax",
	"runtime",
	"runtime/cgo",
	"runtime/debug",
	"runtime/msan",
	"runtime/pprof",
	"runtime/race",
	"runtime/trace",
	"sort",
	"strconv",
	"strings",
	"sync",
	"sync/atomic",
	"syscall",
	"syscall/js",
	"testing",
	"testing/iotest",
	"testing/quick",
	"text",
	"text/scanner",
	"text/tabwriter",
	"text/template",
	"text/template/parseCommand",
	"time",
	"unicode",
	"unicode/utf16",
	"unicode/utf8",
	"unsafe",
}

func addImports(src string) string {
	for _, imp := range stdImports {
		sp := strings.Split(src, "\n")
		if !strings.Contains(src, `"`+imp+`"`) && strings.Contains(src, imp) {
			src = sp[0] + "\nimport \"" + imp + "\"\n" + strings.Join(sp[1:], "\n")
		}
	}
	return src
}

func parseCodeBlock(data string) string {
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "```")
	data = strings.TrimSuffix(data, "```")
	data = strings.TrimPrefix(data, "go")
	data = strings.TrimSpace(data)
	return data
}

func ease(src string) string {
	if !strings.HasPrefix(src, "package main") {
		src = "package main\n" + src
	}
	src = addImports(src)
	return src
}

var Command_go = Command{
	Name:        "go",
	Description: "Run Go code",
	Aliases:     []string{"golang"},
	Execute: func(message *events.MessageCreate, args []string) {
		code := ease(parseCodeBlock(strings.Join(args, " ")))
		fmt.Println(code)
		ms, _ := CreateMessage(message, Message{
			Reply:   true,
			Content: "Running...",
		})
		data := struct {
			Body    string
			WithVet bool
		}{code, false}
		b, _ := json.Marshal(data)
		resp, _ := http.Post("https://play.golang.org/compile", "application/json", bytes.NewReader(b))
		d, _ := io.ReadAll(resp.Body)
		var resps map[string]interface{}
		json.Unmarshal(d, &resps)
		var stdout string
		var errors = resps["Errors"].(string)
		var stderr string
		if m, ok := resps["Events"].([]interface{}); ok {
			for _, msg := range m {
				if msg, ok := msg.(map[string]interface{}); ok {
					if msg["Kind"] == "stdout" {
						stdout += msg["Message"].(string)
					}
					if msg["Kind"] == "stderr" {
						stderr += msg["Message"].(string)
					}
				}
			}
		}
		embed := discord.NewEmbedBuilder().SetColor(color).SetTitle("Code")
		if errors != "" {
			embed.AddField("Compile Errors", "```"+errors+"```", false)
		}
		if stdout != "" {
			embed.AddField("Output", "```"+stdout+"```", false)
		}
		if stderr != "" {
			embed.AddField("Errors", "```"+stderr+"```", false)
		}
		EditMessage(message.Client(), message.ChannelID, ms.ID, Message{
			Embeds: []discord.Embed{embed.Build()},
		})
	},
}
