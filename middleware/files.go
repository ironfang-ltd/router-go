package middleware

import (
	"net/http"

	"github.com/ironfang-ltd/router-go"
)

type FilesOption func(*FilesConfig)

type FilesConfig struct {
	PublicDir string
}

func WithPublicDir(publicDir string) FilesOption {
	return func(c *FilesConfig) {
		c.PublicDir = publicDir
	}
}

func Files(opts ...FilesOption) router.Middleware {

	config := &FilesConfig{
		PublicDir: "./web/static",
	}

	for _, opt := range opts {
		opt(config)
	}

	fs := http.Dir(config.PublicDir)

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		if r.Method != http.MethodGet {
			next(w, r)
			return
		}

		f, err := fs.Open(r.URL.Path)
		if err != nil {
			next(w, r)
			return
		}

		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next(w, r)
			return
		}

		if fi.IsDir() {
			next(w, r)
			return
		}

		http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	}
}
