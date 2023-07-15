package proxy

import (
	"bytes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

var client = &fasthttp.Client{
	NoDefaultUserAgentHeader: true,
	DisablePathNormalizing:   true,
}

// DoTimeout performs the given request and waits for response during the given timeout duration.
// This method can be used within a fiber.Handler
func DoTimeout(c *fiber.Ctx, addr, host string, timeout time.Duration) error {
	return doAction(c, addr, host, func(cli *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		return cli.DoTimeout(req, resp, timeout)
	})
}

func doAction(
	c *fiber.Ctx,
	addr, host string,
	action func(cli *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error,
) error {
	cli := client

	req := c.Request()
	res := c.Response()
	originalURL := utils.CopyString(c.OriginalURL())
	defer req.SetRequestURI(originalURL)

	copiedURL := utils.CopyString(addr)
	req.SetRequestURI(copiedURL)
	req.Header.SetHost(host)
	req.UseHostHeader = true

	// NOTE: if req.isTLS is true, SetRequestURI keeps the scheme as https.
	// Reference: https://github.com/gofiber/fiber/issues/1762
	if scheme := getScheme(utils.UnsafeBytes(copiedURL)); len(scheme) > 0 {
		req.URI().SetSchemeBytes(scheme)
	}

	req.Header.Del(fiber.HeaderConnection)
	if err := action(cli, req, res); err != nil {
		return err
	}
	res.Header.Del(fiber.HeaderConnection)
	return nil
}

func getScheme(uri []byte) []byte {
	i := bytes.IndexByte(uri, '/')
	if i < 1 || uri[i-1] != ':' || i == len(uri)-1 || uri[i+1] != '/' {
		return nil
	}
	return uri[:i-1]
}
