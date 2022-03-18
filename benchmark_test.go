package gzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	// thanks to https://loremipsum.io/generator/?n=25&t=w
	loremIpsum25Words  = `nulla facilisi morbi tempus iaculis urna id volutpat lacus laoreet non curabitur gravida arcu ac tortor dignissim convallis aenean et tortor at risus viverra adipiscing `
	loremIpsum100Words = loremIpsum25Words + " " + loremIpsum25Words + " " + loremIpsum25Words + " " + loremIpsum25Words
)

func loremIpsumManyWords(n00 int) string {
	buf := strings.Builder{}
	for i := 0; i < n00; i++ {
		buf.WriteString(loremIpsum100Words)
	}
	return buf.String()
}

func loremIpsum1kWords() string {
	return loremIpsumManyWords(10)
}

func loremIpsum10kWords() string {
	return loremIpsumManyWords(100)
}

func loremIpsum155kWords() string {
	return loremIpsumManyWords(1550) // approx 1MiB
}

func TestBenchmarkData(t *testing.T) {
	t.Logf("Sizes\n25 words:    %d bytes\n100 words:   %d bytes\n1k words:    %d bytes\n10k words:   %d bytes\n150k words:  %d bytes\n",
		len(loremIpsum25Words), len(loremIpsum100Words), len(loremIpsum1kWords()), len(loremIpsum10kWords()), len(loremIpsum155kWords()))
}

//-------------------------------------------------------------------------------------------------

func BenchmarkGzipDefCompBlank(b *testing.B) {
	gzipBothBenchmark(b, gzip.DefaultCompression, "")
}

func BenchmarkGzipDefCompWords25(b *testing.B) {
	gzipBothBenchmark(b, gzip.DefaultCompression, loremIpsum25Words)
}

func BenchmarkGzipDefCompWords1k(b *testing.B) {
	gzipBothBenchmark(b, gzip.DefaultCompression, loremIpsum1kWords())
}

func BenchmarkGzipDefCompWords10k(b *testing.B) {
	gzipBothBenchmark(b, gzip.DefaultCompression, loremIpsum10kWords())
}

func BenchmarkGzipDefCompWords155k(b *testing.B) {
	gzipBothBenchmark(b, gzip.DefaultCompression, loremIpsum155kWords())
}

//-------------------------------------------------------------------------------------------------

func BenchmarkBothGzipC9Blank(b *testing.B) {
	gzipBothBenchmark(b, 9, "")
}

func BenchmarkBothGzipC9Words25(b *testing.B) {
	gzipBothBenchmark(b, 9, loremIpsum25Words)
}

func BenchmarkBothGzipC9Words1k(b *testing.B) {
	gzipBothBenchmark(b, 9, loremIpsum1kWords())
}

func BenchmarkBothGzipC9Words10k(b *testing.B) {
	gzipBothBenchmark(b, 9, loremIpsum10kWords())
}

func BenchmarkBothGzipC9Words155k(b *testing.B) {
	gzipBothBenchmark(b, 9, loremIpsum155kWords())
}

//-------------------------------------------------------------------------------------------------

func BenchmarkReqGzipC9Blank(b *testing.B) {
	gzipRequestOnlyBenchmark(b, 9, "")
}

func BenchmarkReqGzipC9Words25(b *testing.B) {
	gzipRequestOnlyBenchmark(b, 9, loremIpsum25Words)
}

func BenchmarkReqGzipC9Words1k(b *testing.B) {
	gzipRequestOnlyBenchmark(b, 9, loremIpsum1kWords())
}

func BenchmarkReqGzipC9Words10k(b *testing.B) {
	gzipRequestOnlyBenchmark(b, 9, loremIpsum10kWords())
}

func BenchmarkReqGzipC9Words155k(b *testing.B) {
	gzipRequestOnlyBenchmark(b, 9, loremIpsum155kWords())
}

//-------------------------------------------------------------------------------------------------

func BenchmarkResGzipC9Blank(b *testing.B) {
	compression := []int{gzip.DefaultCompression, 1, 3, 6, 9}
	for _, comp := range compression {
		name := fmt.Sprintf("C%d", comp)
		b.Run(name+"_blank", func(b *testing.B) {
			gzipResponseOnlyBenchmark(b, 9, "")
		})
		b.Run(name+"_Words25", func(b *testing.B) {
			gzipResponseOnlyBenchmark(b, 9, loremIpsum25Words)
		})
		b.Run(name+"_Words1k", func(b *testing.B) {
			gzipResponseOnlyBenchmark(b, 9, loremIpsum1kWords())
		})
		b.Run(name+"_Words10k", func(b *testing.B) {
			gzipResponseOnlyBenchmark(b, 9, loremIpsum10kWords())
		})
		b.Run(name+"_Words155k", func(b *testing.B) {
			gzipResponseOnlyBenchmark(b, 9, loremIpsum155kWords())
		})
	}
}

//-------------------------------------------------------------------------------------------------

//func BenchmarkGzipC4Blank(b *testing.B) {
//	gzipBothBenchmark(b, 4, "")
//}
//
//func BenchmarkGzipC4Words25(b *testing.B) {
//	gzipBothBenchmark(b, 4, loremIpsum25Words)
//}
//
//func BenchmarkGzipC4Words1k(b *testing.B) {
//	gzipBothBenchmark(b, 4, loremIpsum1kWords())
//}
//
//func BenchmarkGzipC4Words10k(b *testing.B) {
//	gzipBothBenchmark(b, 4, loremIpsum10kWords())
//}
//
//func BenchmarkGzipC4Words155k(b *testing.B) {
//	gzipBothBenchmark(b, 4, loremIpsum155kWords())
//}

//-------------------------------------------------------------------------------------------------

//func BenchmarkReqGzipC4Blank(b *testing.B) {
//	gzipRequestOnlyBenchmark(b, 4, "")
//}
//
//func BenchmarkReqGzipC4Words25(b *testing.B) {
//	gzipRequestOnlyBenchmark(b, 4, loremIpsum25Words)
//}
//
//func BenchmarkReqGzipC4Words1k(b *testing.B) {
//	gzipRequestOnlyBenchmark(b, 4, loremIpsum1kWords())
//}
//
//func BenchmarkReqGzipC4Words10k(b *testing.B) {
//	gzipRequestOnlyBenchmark(b, 4, loremIpsum10kWords())
//}
//
//func BenchmarkReqGzipC4Words155k(b *testing.B) {
//	gzipRequestOnlyBenchmark(b, 4, loremIpsum155kWords())
//}

//-------------------------------------------------------------------------------------------------

//func BenchmarkResGzipC4Blank(b *testing.B) {
//	gzipResponseOnlyBenchmark(b, 4, "")
//}
//
//func BenchmarkResGzipC4Words25(b *testing.B) {
//	gzipResponseOnlyBenchmark(b, 4, loremIpsum25Words)
//}
//
//func BenchmarkResGzipC4Words1k(b *testing.B) {
//	gzipResponseOnlyBenchmark(b, 4, loremIpsum1kWords())
//}
//
//func BenchmarkResGzipC4Words10k(b *testing.B) {
//	gzipResponseOnlyBenchmark(b, 4, loremIpsum10kWords())
//}
//
//func BenchmarkResGzipC4Words155k(b *testing.B) {
//	gzipResponseOnlyBenchmark(b, 4, loremIpsum155kWords())
//}

//-------------------------------------------------------------------------------------------------

//func BenchmarkGzipC1Blank(b *testing.B) {
//	gzipBothBenchmark(b, 1, "")
//}
//
//func BenchmarkGzipC1Words25(b *testing.B) {
//	gzipBothBenchmark(b, 1, loremIpsum25Words)
//}
//
//func BenchmarkGzipC1Words1k(b *testing.B) {
//	gzipBothBenchmark(b, 1, loremIpsum1kWords())
//}
//
//func BenchmarkGzipC1Words10k(b *testing.B) {
//	gzipBothBenchmark(b, 1, loremIpsum10kWords())
//}
//
//func BenchmarkGzipC1Words155k(b *testing.B) {
//	gzipBothBenchmark(b, 1, loremIpsum155kWords())
//}

//-------------------------------------------------------------------------------------------------

func BenchmarkPlainBlank(b *testing.B) {
	noGzipBenchmark(b, "")
}

func BenchmarkPlainWords25(b *testing.B) {
	noGzipBenchmark(b, loremIpsum25Words)
}

func BenchmarkPlainWords1k(b *testing.B) {
	noGzipBenchmark(b, loremIpsum1kWords())
}

func BenchmarkPlainWords10k(b *testing.B) {
	noGzipBenchmark(b, loremIpsum10kWords())
}

//-------------------------------------------------------------------------------------------------

func compressedBuffer(b *testing.B, compression int, text string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	gz, _ := gzip.NewWriterLevel(buf, compression)
	if _, err := gz.Write([]byte(text)); err != nil {
		gz.Close()
		b.Fatal(err)
	}
	gz.Close()
	return buf
}

func ginWithRequestDecompressor(compression int) *gin.Engine {
	router := gin.New()
	router.Use(Gzip(compression, WithDecompressFn(DefaultDecompressHandle)))
	return router
}

func gzipBothBenchmark(b *testing.B, compression int, text string) {
	buf := compressedBuffer(b, compression, text)
	router := ginWithRequestDecompressor(compression)
	router.POST("/", func(c *gin.Context) {
		data, _ := c.GetRawData()
		c.Data(200, "text/plain", data)
	})

	req := newRequestWithGzipResponse(buf)
	runBenchmark(b, router, req, len(text))
}

func gzipRequestOnlyBenchmark(b *testing.B, compression int, text string) {
	buf := compressedBuffer(b, compression, text)
	router := ginWithRequestDecompressor(compression)
	router.POST("/", func(c *gin.Context) {
		data, _ := c.GetRawData()
		c.Data(200, "text/plain", data)
	})

	req := newPlainRequest(buf)
	runBenchmark(b, router, req, len(text))
}

func gzipResponseOnlyBenchmark(b *testing.B, compression int, text string) {
	buf := strings.NewReader(text)
	router := gin.New()
	router.Use(Gzip(compression))
	router.POST("/", func(c *gin.Context) {
		data, _ := c.GetRawData()
		c.Data(200, "text/plain", data)
	})

	req := newRequestWithGzipResponse(buf)
	runBenchmark(b, router, req, len(text))
}

func noGzipBenchmark(b *testing.B, text string) {
	buf := strings.NewReader(text)
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		data, _ := c.GetRawData()
		c.Data(200, "text/plain", data[:0])
	})

	req := newPlainRequest(buf)
	runBenchmark(b, router, req, len(text))
}

func newPlainRequest(body io.Reader) *http.Request {
	req, _ := http.NewRequest("POST", "/", body)
	return req
}

func newRequestWithGzipResponse(body io.Reader) *http.Request {
	req, _ := http.NewRequest("POST", "/", body)
	req.Header.Add("Content-Encoding", "gzip")
	return req
}

func runBenchmark(b *testing.B, router *gin.Engine, req *http.Request, size int) {
	w := httptest.NewRecorder()
	w.Body.Grow(size)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		router.ServeHTTP(w, req)

		w.HeaderMap = make(http.Header)
		w.Body.Reset()
	}
}
