package docparser

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/morehao/golib/glog"
)

func stringOr(val, fallback string) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return fallback
	}
	return val
}

func parseBoolOr(val string, fallback bool) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	if val == "" {
		return fallback
	}
	return val == "true" || val == "1" || val == "yes"
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func sleepCtx(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

func logResponseStructure(label string, obj interface{}, prefix string) {
	switch v := obj.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		glog.Infof(nil, "[%s] %s = {object with %d keys: %s}", label, prefix, len(v), strings.Join(keys, ", "))
		for _, key := range keys {
			val := v[key]
			path := prefix + "." + key
			switch inner := val.(type) {
			case map[string]interface{}:
				logResponseStructure(label, inner, path)
			case []interface{}:
				glog.Infof(nil, "[%s] %s = [array with %d items]", label, path, len(inner))
				if len(inner) > 0 {
					glog.Infof(nil, "[%s] %s[0] type=%T", label, path, inner[0])
					if len(inner) <= 3 {
						for i, item := range inner {
							logResponseStructure(label, item, fmt.Sprintf("%s[%d]", path, i))
						}
					} else {
						logResponseStructure(label, inner[0], path+"[0]")
						glog.Infof(nil, "[%s] ... and %d more items in %s", label, len(inner)-1, path)
					}
				}
			case string:
				if len(inner) > 200 {
					glog.Infof(nil, "[%s] %s = string(%d chars): %.200s...", label, path, len(inner), inner)
				} else {
					glog.Infof(nil, "[%s] %s = %q", label, path, inner)
				}
			case float64:
				glog.Infof(nil, "[%s] %s = %v (number)", label, path, inner)
			case bool:
				glog.Infof(nil, "[%s] %s = %v (bool)", label, path, inner)
			case nil:
				glog.Infof(nil, "[%s] %s = null", label, path)
			default:
				glog.Infof(nil, "[%s] %s = %v (%T)", label, path, val, val)
			}
		}
	case []interface{}:
		glog.Infof(nil, "[%s] %s = [array with %d items]", label, prefix, len(v))
		if len(v) > 0 {
			if len(v) <= 3 {
				for i, item := range v {
					logResponseStructure(label, item, fmt.Sprintf("%s[%d]", prefix, i))
				}
			} else {
				logResponseStructure(label, v[0], prefix+"[0]")
				glog.Infof(nil, "[%s] ... and %d more items in %s", label, len(v)-1, prefix)
			}
		}
	default:
		glog.Infof(nil, "[%s] %s = %v (%T)", label, prefix, v, v)
	}
}
