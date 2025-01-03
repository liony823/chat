// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kdisc

import (
	"time"

	"github.com/liony823/tools/discovery"
	"github.com/liony823/tools/discovery/etcd"
	"github.com/liony823/tools/discovery/kubernetes"
	"github.com/liony823/tools/errs"

	"github.com/openimsdk/chat/pkg/common/config"
)

const (
	ETCDCONST       = "etcd"
	KUBERNETESCONST = "kubernetes"
	DIRECTCONST     = "direct"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(discovery *config.Discovery, runtimeEnv string) (discovery.SvcDiscoveryRegistry, error) {
	if runtimeEnv == KUBERNETESCONST {
		return kubernetes.NewKubernetesConnManager(discovery.Kubernetes.Namespace)
	}

	switch discovery.Enable {
	case ETCDCONST:
		return etcd.NewSvcDiscoveryRegistry(
			discovery.Etcd.RootDirectory,
			discovery.Etcd.Address,
			etcd.WithDialTimeout(60*time.Second),
			etcd.WithMaxCallSendMsgSize(20*1024*1024),
			etcd.WithUsernameAndPassword(discovery.Etcd.Username, discovery.Etcd.Password))
	default:
		return nil, errs.New("unsupported discovery type", "type", discovery.Enable).Wrap()
	}
}
