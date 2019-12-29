@echo off
clear
SETLOCAL
go fmt
@REM set GOCACHE=off
golint ./... > .sonarqube\golint-report.out 2>&1
go vet ./... > .sonarqube\govet-report.out 2>&1
golangci-lint run --out-format checkstyle ./... > .sonarqube/linter-report.xml
go test -json -coverprofile=.sonarqube\cover.out ./... > .sonarqube\test-report.json
IF %ERRORLEVEL% NEQ 0 (
  cat  .sonarqube\test-report.json
) ELSE (
  sonar
  pushd .sonarqube
  move /Y test-report.json test-report-success.json
  del golint-report.out cover.out govet-report.out linter-report.xml
  popd
)
endlocal
