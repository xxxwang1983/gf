// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// GrpcServerConfig is the configuration for server.
type GrpcServerConfig struct {
	Address          string              // (optional) Address for server listening.
	Name             string              // (optional) Name for current service.
	Logger           *glog.Logger        // (optional) Logger for server.
	LogPath          string              // (optional) LogPath specifies the directory for storing logging files.
	LogStdout        bool                // (optional) LogStdout specifies whether printing logging content to stdout.
	ErrorStack       bool                // (optional) ErrorStack specifies whether logging stack information when error.
	ErrorLogEnabled  bool                // (optional) ErrorLogEnabled enables error logging content to files.
	ErrorLogPattern  string              // (optional) ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	AccessLogEnabled bool                // (optional) AccessLogEnabled enables access logging content to file.
	AccessLogPattern string              // (optional) AccessLogPattern specifies the error log file pattern like: access-{Ymd}.log
	Endpoints        []string            // (optional) Address for server register if null use Address value.
	Options          []grpc.ServerOption // (optional) GRPC Server options.
}

// NewConfig creates and returns a ServerConfig object with default configurations.
// Note that, do not define this default configuration to local package variable, as there are
// some pointer attributes that may be shared in different servers.
func (s modServer) NewConfig() *GrpcServerConfig {
	var (
		err    error
		ctx    = context.TODO()
		config = &GrpcServerConfig{
			Name:             defaultServerName,
			Logger:           glog.New(),
			LogStdout:        true,
			ErrorLogEnabled:  true,
			ErrorLogPattern:  "error-{Ymd}.log",
			AccessLogEnabled: false,
			AccessLogPattern: "access-{Ymd}.log",
		}
	)
	// Reading configuration file and updating the configured keys.
	if g.Cfg().Available(ctx) {
		// Server attributes configuration.
		serverConfigMap := g.Cfg().MustGet(ctx, configNodeNameGrpcServer).Map()
		if len(serverConfigMap) == 0 {
			return config
		}
		if err = gconv.Struct(serverConfigMap, &config); err != nil {
			g.Log().Error(ctx, err)
			return config
		}
		// Server logger configuration checks.
		serverLoggerConfigMap := g.Cfg().MustGet(
			ctx,
			fmt.Sprintf(`%s.logger`, configNodeNameGrpcServer),
		).Map()
		if len(serverLoggerConfigMap) == 0 && len(serverConfigMap) > 0 {
			serverLoggerConfigMap = gconv.Map(serverConfigMap["logger"])
		}
		if len(serverLoggerConfigMap) > 0 {
			if err = config.Logger.SetConfigWithMap(serverLoggerConfigMap); err != nil {
				panic(err)
			}
		}
	}
	return config
}

// SetWithMap changes current configuration with map.
// This is commonly used for changing several configurations of current object.
func (c *GrpcServerConfig) SetWithMap(m g.Map) error {
	return gconv.Struct(m, c)
}

// MustSetWithMap acts as SetWithMap but panics if error occurs.
func (c *GrpcServerConfig) MustSetWithMap(m g.Map) {
	err := c.SetWithMap(m)
	if err != nil {
		panic(err)
	}
}
