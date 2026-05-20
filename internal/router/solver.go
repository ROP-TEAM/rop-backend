package router

import (
	"github.com/ROP-TEAM/rop-algorithm/solver"
	grpcsolver "github.com/ROP-TEAM/rop-algorithm/solver/grpc"
	"github.com/ROP-TEAM/rop-algorithm/solver/process"
)

func buildSolver(binaryPath string) (solver.Solver, func(), error) {
	if binaryPath == "" {
		return solver.NewStub(), nil, nil
	}

	handle, err := process.Start(binaryPath)
	if err != nil {
		return nil, nil, err
	}

	return grpcsolver.New(handle.Conn), handle.Stop, nil
}
