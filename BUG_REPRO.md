# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	theatrecontrol/cmd/controlroom	[no test files]
?   	theatrecontrol/internal/model	[no test files]
?   	theatrecontrol/internal/store	[no test files]
ok  	theatrecontrol/internal/audit	0.011s
ok  	theatrecontrol/internal/integration	0.023s
ok  	theatrecontrol/internal/rehearsal	0.012s
ok  	theatrecontrol/internal/report	0.007s
--- FAIL: TestRoleAllowsEmptyMenuSelection (0.00s)
panic: assignment to entry in nil map [recovered]
	panic: assignment to entry in nil map

goroutine 21 [running]:
testing.tRunner.func1.2({0x169440, 0x1cc310})
	/usr/local/go/src/testing/testing.go:1631 +0x1c4
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1634 +0x33c
panic({0x169440?, 0x1cc310?})
	/usr/local/go/src/runtime/panic.go:770 +0x124
theatrecontrol/internal/roles.(*Service).SaveRole(0x4000105ed8, {{0x190e58, 0xa}, {0x18ff34, 0x6}, {0x18fd57, 0x5}, {0x333e20, 0x0, 0x0}, ...})
	/app/internal/roles/service.go:30 +0x1e4
theatrecontrol/internal/roles.TestRoleAllowsEmptyMenuSelection(0x40000d09c0)
	/app/internal/roles/roles_test.go:18 +0x184
testing.tRunner(0x40000d09c0, 0x19dba0)
	/usr/local/go/src/testing/testing.go:1689 +0xec
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:1742 +0x318
FAIL	theatrecontrol/internal/roles	0.009s
ok  	theatrecontrol/internal/show	0.009s
ok  	theatrecontrol/internal/wardrobe	0.009s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/controlroom): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/controlroom): exit `0`
