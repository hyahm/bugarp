package handle

import (
	"errors"
	"itflow/cache"
	"itflow/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var NotFoundToken = errors.New("not found token")

func sortpermlist(permlist []string) []string {
	l := len(cache.CacheSidStatus)

	newlist := make([]string, 0)
	for i := 0; i < l; i++ {
		for _, v := range permlist {
			if cache.CacheSidStatus[int64(i)] == v {
				newlist = append(newlist, v)
			}
		}
	}
	return newlist
}

// []string{ "admin(admin)" }  前端传过来的数据换成数组插入到data
func formatUserlistToData(userlist []string, bid int) (string, []string, [][]interface{}) {
	ul := ""
	l := len(userlist)
	nicknamelist := make([]string, 0)
	args := make([][]interface{}, 0)
	for j, v := range userlist {

		onearg := make([]interface{}, 0)
		i := strings.Index(v, "(")

		nickname := v[:i]
		if j == l-1 {
			ul = ul + strconv.FormatInt(cache.CacheNickNameUid[v[:i]], 10)
		} else {
			ul = ul + strconv.FormatInt(cache.CacheNickNameUid[v[:i]], 10) + ","
		}
		onearg = append(onearg, bid)
		onearg = append(onearg, cache.CacheNickNameUid[v[:i]])
		args = append(args, onearg)
		nicknamelist = append(nicknamelist, nickname)
	}

	return ul, nicknamelist, args
}

// bugs表中的spuser返回昵称和真实用户名的组合（直接显示在前端）

// bugs表中的spuser返回真实用户名
func formatUserlistToRealname(userlist string) []string {
	al := make([]string, 0)
	ul := strings.Split(userlist, ",")
	for _, v := range ul {
		uid, _ := strconv.Atoi(v)
		al = append(al, cache.CacheUidRealName[int64(uid)])
	}
	return al
}

// 插入到log表中
func insertlog(classify string, content string, r *http.Request) error {
	logsql := "insert into log(exectime,classify,content,ip) values(?,?,?,?)"
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if ip != "127.0.0.1" {
		_, err := db.Mconn.Insert(logsql, time.Now().Unix(), classify, content, ip)
		if err != nil {
			return err
		}
	}

	return nil
}