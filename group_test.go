package group_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	group "github.com/Drelf2018/gin-group"
	"github.com/gin-gonic/gin"
)

func GetPing(ctx *gin.Context) (any, error) {
	return "pong", nil
}

func GetCover(ctx *gin.Context) (any, error) {
	mid := ctx.Query("mid")

	if !strings.HasPrefix(mid, "BV") {
		return 1, fmt.Errorf("example: %v is an invalid bvid", mid)
	}

	resp, err := http.Get(fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", mid))
	if err != nil {
		return 2, err
	}
	defer resp.Body.Close()

	p, err := io.ReadAll(resp.Body)
	if err != nil {
		return 3, err
	}

	var view struct {
		Data struct {
			Pic string
		}
	}
	err = json.Unmarshal(p, &view)
	if err != nil {
		return 4, err
	}

	return view.Data.Pic, nil
}

type File struct {
	MIME string
	Data []byte
}

var cache = make(map[string]string)
var files = make(map[string]File)

var ErrNoFile = errors.New("example: the file does not exist")

func GetResourceFile(ctx *gin.Context) (any, error) {
	file, ok := files[ctx.Param("file")]
	if !ok {
		return 1, ErrNoFile
	}
	ctx.Data(http.StatusOK, file.MIME, file.Data)
	return nil, nil
}

func GetDownload(ctx *gin.Context) (any, error) {
	url := ctx.Query("url")

	if file, ok := cache[url]; ok {
		return file, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return 1, err
	}

	p, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return 2, err
	}

	s := fmt.Sprint(rand.Float64())[2:]
	cache[url] = s
	files[s] = File{
		MIME: resp.Header.Get("Content-Type"),
		Data: p,
	}
	return s, nil
}

func init() {
	api := group.Group{
		Middleware: group.CORS,
		Handlers: []group.H{
			GetPing,
			GetCover,
		},
		Groups: []group.Group{{
			Path: "admin",
			Middleware: func(ctx *gin.Context) {
				if ctx.Query("name") != "admin" {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, group.Response{
						Code:  1,
						Error: "you are not administrator!",
					})
				}
			},
			Handlers: []group.H{
				GetDownload,
				GetResourceFile,
			},
		}},
	}
	// gin.SetMode(gin.ReleaseMode)
	go api.Default().Run("localhost:8080")
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func TestPing(t *testing.T) {
	b, err := get("http://localhost:8080/ping")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestCover(t *testing.T) {
	b, err := get("http://localhost:8080/cover?mid=abc123")
	if err != nil {
		t.Fatal(err)
	}
	r, err := group.Unmarshal[string](b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r.Data)
}

func TestDownload(t *testing.T) {
	b, err := get("http://localhost:8080/cover?mid=BV1hxmwYDEJ6")
	if err != nil {
		t.Fatal(err)
	}
	r, err := group.Unmarshal[string](b)
	if err != nil {
		t.Fatal(err)
	}
	cover := r.Data
	t.Log("cover:", cover)

	b, err = get("http://localhost:8080/admin/download?name=admin&url=" + cover)
	if err != nil {
		t.Fatal(err)
	}
	r, err = group.Unmarshal[string](b)
	if err != nil {
		t.Fatal(string(b), err)
	}
	t.Log("uuid:", r.Data)

	b, err = get("http://localhost:8080/admin/resource/" + r.Data + "?name=admin")
	if err != nil {
		t.Fatal(err)
	}

	_, file := filepath.Split(cover)
	err = os.WriteFile(file, b, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
