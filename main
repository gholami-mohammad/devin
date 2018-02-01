package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func gitCommand(dir string, args ...string) []byte {
	log.Println(args)
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		log.Println("ERR ", err)
	}
	return out
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

func main() {
	app := iris.New()
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	sub := app.Party("/{username}/{repo}")
	{
		sub.Get("/info/refs", func(ctx iris.Context) {
			service := strings.TrimPrefix(ctx.Request().URL.Query().Get("service"), "git-")
			dir := "/Users/mgh/gogs-repositories/" + ctx.Params().Get("username") + "/" + ctx.Params().Get("repo")
			refs := gitCommand(dir, service, "--stateless-rpc", "--advertise-refs", dir)

			ctx.ResponseWriter().Header().Set("WWW-Authenticate", "Basic realm=\".\"")
			ctx.ResponseWriter().Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
			ctx.ResponseWriter().WriteHeader(http.StatusOK)
			ctx.ResponseWriter().Write(packetWrite("# service=git-" + service + "\n"))
			ctx.ResponseWriter().Write([]byte("0000"))
			ctx.ResponseWriter().Write(refs)
		})

		sub.Get("HEAD", func(ctx iris.Context) {

		})

	}

	// {regexp.MustCompile("(.*?)/git-upload-pack$"), "POST", serviceUploadPack},
	// {regexp.MustCompile("(.*?)/git-receive-pack$"), "POST", serviceReceivePack},
	// {regexp.MustCompile("(.*?)/info/refs$"), "GET", getInfoRefs},
	// {regexp.MustCompile("(.*?)/HEAD$"), "GET", getTextFile},
	// {regexp.MustCompile("(.*?)/objects/info/alternates$"), "GET", getTextFile},
	// {regexp.MustCompile("(.*?)/objects/info/http-alternates$"), "GET", getTextFile},
	// {regexp.MustCompile("(.*?)/objects/info/packs$"), "GET", getInfoPacks},
	// {regexp.MustCompile("(.*?)/objects/info/[^/]*$"), "GET", getTextFile},
	// {regexp.MustCompile("(.*?)/objects/[0-9a-f]{2}/[0-9a-f]{38}$"), "GET", getLooseObject},
	// {regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.pack$"), "GET", getPackFile},
	// {regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.idx$"), "GET", getIdxFile},

	app.Run(iris.Addr(":3000"), iris.WithoutServerError(iris.ErrServerClosed))
}
