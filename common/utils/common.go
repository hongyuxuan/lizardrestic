package utils

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hongyuxuan/lizardrestic/common/errorx"
)

func GetLizardAgentKey(key []byte) string {
	arr := strings.Split(string(key), "/")
	uid := arr[len(arr)-1]
	return strings.TrimSuffix(string(key), "/"+uid)
}

func GetServiceMata(prefix, key string) (map[string]string, error) {
	re, _ := regexp.Compile(prefix + "lizardrestic-agent\\.(.+?)\\.(.+)")
	res := re.FindStringSubmatch(key)
	if res == nil {
		return nil, errorx.NewDefaultError("No match of \"%s\" to <ServicePrefix>.lizardrestic-agent.<system>.<ip>", key)
	}
	return map[string]string{
		"Protocol": "grpc",
		"Service":  prefix + "lizardrestic-agent",
		"System":   res[1],
		"IP":       res[2],
	}, nil
}

func GetTarget(prefix, key string) (system, ip string, err error) {
	re, _ := regexp.Compile(prefix + "lizardrestic-agent\\.(.+?)\\.(.+)")
	res := re.FindStringSubmatch(key)
	if res == nil {
		return "", "", errorx.NewDefaultError("No match of \"%s\" to <ServicePrefix>.lizardrestic-agent.<system>.<ip>", key)
	}
	return res[1], res[2], nil
}

func GetPayload(ctx context.Context) (username, role string, tenant, namespaces []string) {
	payload := ctx.Value("payloads").(map[string]interface{})
	username = payload["username"].(string)
	role = payload["role"].(string)
	tenantStr := payload["tenant"].(string)
	tenant = strings.Split(tenantStr, ",")
	namespaceStr := payload["namespace"].(string)
	namespaces = strings.Split(namespaceStr, ",")
	return
}

func AnyToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case bool:
		return strconv.FormatBool(v.(bool))
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	default:
		return fmt.Sprintf("%+v", v)
	}
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
