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

package startrpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/liony823/tools/utils/datautil"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/kdisc"

	"github.com/liony823/tools/discovery"
	"github.com/liony823/tools/errs"
	"github.com/liony823/tools/log"
	"github.com/liony823/tools/mw"
	"github.com/liony823/tools/system/program"
	"github.com/liony823/tools/utils/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Start rpc server.
func Start[T any](ctx context.Context, discovery *config.Discovery, listenIP,
	registerIP string, rpcPorts []int, index int, rpcRegisterName string, share *config.Share, config T, rpcFn func(ctx context.Context,
	config T, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {

	rpcPort, err := datautil.GetElemByIndex(rpcPorts, index)
	if err != nil {
		return err
	}
	log.CInfo(ctx, "RPC server is initializing", "rpcRegisterName", rpcRegisterName, "rpcPort", rpcPort)
	rpcTcpAddr := net.JoinHostPort(network.GetListenIP(listenIP), strconv.Itoa(rpcPort))
	listener, err := net.Listen(
		"tcp",
		rpcTcpAddr,
	)
	if err != nil {
		return errs.WrapMsg(err, "listen err", "rpcTcpAddr", rpcTcpAddr)
	}

	defer listener.Close()
	client, err := kdisc.NewDiscoveryRegister(discovery)
	if err != nil {
		return err
	}
	defer client.Close()
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	registerIP, err = network.GetRpcRegisterIP(registerIP)
	if err != nil {
		return err
	}

	options = append(options, mw.GrpcServer())
	srv := grpc.NewServer(options...)
	once := sync.Once{}
	defer func() {
		once.Do(srv.GracefulStop)
	}()

	err = rpcFn(ctx, config, client, srv)
	if err != nil {
		return err
	}

	if err := client.Register(rpcRegisterName, registerIP, rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	}

	var (
		netDone = make(chan struct{}, 2)
		netErr  error
	)

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			netErr = errs.WrapMsg(err, "rpc start err: ", rpcTcpAddr)
			netDone <- struct{}{}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := gracefulStopWithCtx(ctx, srv.GracefulStop); err != nil {
			return err
		}
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return nil
	case <-netDone:
		close(netDone)
		return netErr
	}
}

func gracefulStopWithCtx(ctx context.Context, f func()) error {
	done := make(chan struct{}, 1)
	go func() {
		f()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return errs.New("timeout, ctx graceful stop")
	case <-done:
		return nil
	}
}
