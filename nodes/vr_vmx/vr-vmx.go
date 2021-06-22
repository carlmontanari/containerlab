// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package vr_vmx

import (
	"context"
	"fmt"

	"github.com/srl-labs/containerlab/nodes"
	"github.com/srl-labs/containerlab/runtime"
	"github.com/srl-labs/containerlab/types"
	"github.com/srl-labs/containerlab/utils"
)

const (
	nodeKind = "vr-vmx"
)

func init() {
	nodes.Register(nodeKind, func() nodes.Node {
		return new(vrVMX)
	})
}

type vrVMX struct {
	cfg  *types.NodeConfig
	mgmt *types.MgmtNet
}

func (s *vrVMX) Init(cfg *types.NodeConfig, opts ...nodes.NodeOption) error {
	s.cfg = cfg
	for _, o := range opts {
		o(s)
	}
	// env vars are used to set launch.py arguments in vrnetlab container
	defEnv := map[string]string{
		"USERNAME":           "admin",
		"PASSWORD":           "admin@123",
		"CONNECTION_MODE":    nodes.VrDefConnMode,
		"DOCKER_NET_V4_ADDR": s.mgmt.IPv4Subnet,
		"DOCKER_NET_V6_ADDR": s.mgmt.IPv6Subnet,
	}
	s.cfg.Env = utils.MergeStringMaps(defEnv, s.cfg.Env)

	if s.cfg.Env["CONNECTION_MODE"] == "macvtap" {
		// mount dev dir to enable macvtap
		s.cfg.Binds = append(s.cfg.Binds, "/dev:/dev")
	}

	s.cfg.Cmd = fmt.Sprintf("--username %s --password %s --hostname %s --connection-mode %s --trace",
		s.cfg.Env["USERNAME"], s.cfg.Env["PASSWORD"], s.cfg.ShortName, s.cfg.Env["CONNECTION_MODE"])

	return nil
}

func (s *vrVMX) Config() *types.NodeConfig { return s.cfg }

func (s *vrVMX) PreDeploy(configName, labCADir, labCARoot string) error {
	utils.CreateDirectory(s.cfg.LabDir, 0777)
	return nil
}

func (s *vrVMX) Deploy(ctx context.Context, r runtime.ContainerRuntime) error {
	return r.CreateContainer(ctx, s.cfg)
}

func (s *vrVMX) PostDeploy(ctx context.Context, r runtime.ContainerRuntime, ns map[string]nodes.Node) error {
	return nil
}

func (s *vrVMX) WithMgmtNet(mgmt *types.MgmtNet) { s.mgmt = mgmt }