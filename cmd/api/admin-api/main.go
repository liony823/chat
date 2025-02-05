// Copyright © 2023 OpenIM open source community. All rights reserved.
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

package main

import (
	"github.com/openimsdk/chat/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

// @title          飞宏IM后台管理接口
// @version        1.0.0
// @description    飞宏IM后台管理接口

// @contact.name   飞宏IM技术支持
// @contact.url    https://github.com/liony823/chat
// @contact.email  ilovecoding@foxmail.com

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:10009
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       token
// @description               "Type 'Bearer' followed by a space and JWT token"

func main() {
	if err := cmd.NewAdminApiCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
