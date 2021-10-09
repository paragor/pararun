package hacks_centos7

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func EnableUserNamespaces() error {
	data, err := ioutil.ReadFile("/proc/sys/user/max_user_namespaces")
	if err != nil {
		return err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return fmt.Errorf("cant atoi content of /proc/sys/user/max_user_namespaces (%s): %w", string(data), err)
	}
	if count == 0 {
		err = ioutil.WriteFile("/proc/sys/user/max_user_namespaces", []byte("10000"), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
