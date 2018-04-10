default:
	@echo "No default target."

# For each test file, the number should be two
check_corresponding_output_directory:
	@find . -name '*.yaml' | \
	sed 'sW\.\(.*\)\.yamlWecho -n ".\1.yaml " \&\& cat .\1.yaml | grep "/tmp/robin\1" | wc -lW' | bash

clean:
	rm -f *.covout *-cover.bats ./robin ./robin.cover

cover: clean
	./cover.bash

test: clean
	bats *.bats
