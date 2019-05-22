package api


import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// logging will wrap all handlers for logging
func logging(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := &responseLogger{ResponseWriter: w}
		inner.ServeHTTP(l, r)
		log.Println(buildLogLine(l, r, start))
	})
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP status
// code and body size
type responseLogger struct {
	http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// Set status if WriteHeader has not been called
		l.status = http.StatusOK
	}

	size, err := l.ResponseWriter.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.ResponseWriter.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	if l.status == 0 {
		// This can happen if we never actually write data, but only set response headers.
		l.status = http.StatusOK
	}
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

// Common Log Format: http://en.wikipedia.org/wiki/Common_Log_Format

// buildLogLine creates a common log format
// in addition to the common fields, we also append referrer, user agent,
// request ID and response time (microseconds)
// ie, in apache mod_log_config terms:
//     %h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-agent}i\"" %L %D
func buildLogLine(l *responseLogger, r *http.Request, start time.Time) string {

	username := parseUsername(r)

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	if xff := r.Header["X-Forwarded-For"]; xff != nil {
		addrs := append(xff, host)
		host = strings.Join(addrs, ",")
	}

	uri := r.URL.RequestURI()

	referer := r.Referer()

	userAgent := r.UserAgent()

	return fmt.Sprintf(`%s - %s [%s] "%s %s %s" %s %s "%s" "%s" %s %d`,
		host,
		detect(username, "-"),
		start.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		uri,
		r.Proto,
		detect(strconv.Itoa(l.Status()), "-"),
		strconv.Itoa(l.Size()),
		detect(referer, "-"),
		detect(userAgent, "-"),
		r.Header.Get("Request-Id"),
		// response time, report in microseconds because this is consistent
		// with apache's %D parameter in mod_log_config
		int64(time.Since(start)/time.Microsecond))
}

// detect detects the first presence of a non blank string and returns it
func detect(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// parses the username either from the url or auth header
func parseUsername(r *http.Request) string {
	url := r.URL

	// get username from the url if passed there
	if url.User != nil {
		if name := url.User.Username(); name != "" {
			return name
		}
	}

	// Try to get it from the authorization header if set there
	if u, _, ok := r.BasicAuth(); ok {
		return u
	}

	// If your system uses a specific parameter, you can hard code them here
	for _, v := range []string{"u", "user"} {
		// Try to get the username from the query param 'u'
		q := url.Query()
		if u := q.Get(v); u != "" {
			return u
		}
	}
	return ""
}
