#!/usr/bin/env bash
cd paire && \
go test -v && \
cd .. && \
go test -v *unit* && \
go test -v *integration* && \
cd paire/cmd/push && \
gox -osarch "darwin/amd64" -output="paire-push_linux_amd64"
mv paire
cd ../pull
gox -osarch "darwin/amd64" -output="paire-pull_linux_amd64"
cd ../../.. && \
go test -v *e2e*
