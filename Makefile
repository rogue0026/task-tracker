.PHONY build:
build:
	go build -o ./cmd/bin/bot ./cmd/task-tracker/main.go

.PHONY run:
run:
	./cmd/bin/bot -c ./configs/test_cfg.yaml