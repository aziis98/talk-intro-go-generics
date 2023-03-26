
subprojects = $(patsubst examples/%, %, $(wildcard examples/*))

.PHONY: usage
usage:
	@printf "Found the following subprojects:\n"
	@for project in $(subprojects); do printf "  - %s\n" $$project; done
	@printf "\n"
	@printf "Available rules\n"
	@printf "  - run-<name>: Will run the program\n"
	@printf "  - compile-<name>: Will compile the program\n"
	@printf "  - compile-noinline-<name>: Will compile the program without inlining\n"
	@printf "  - decomp-<name>: Will run lensm on the project\n"
	@printf "  - decomp-noinline-<name>: Will run lensm on the project without inlining\n"
	@printf "\n"

run-%: FORCE
	go run -v ./examples/$*

compile-%: FORCE
	go build -v -o "./bin/$*" ./examples/$*

compile-noinline-%: FORCE
	go build -v -o "./bin/$*@no-inline" -gcflags=-l ./examples/$*

decomp-%: compile-% FORCE
	lensm -text-size 16 -filter main "./bin/$*"

decomp-noinline-%: compile-noinline-% FORCE
	lensm -text-size 16 -filter main "./bin/$*@no-inline"

FORCE: ;
