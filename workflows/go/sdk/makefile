# Fans `make <target>` out to every subfolder that has its own makefile.
# Auto-discovers examples via wildcard, so dropping a new sibling project in
# (e.g. patient-intake/ once it's runnable) requires no edits here.

APPS := $(patsubst %/makefile,%,$(wildcard */makefile))

.PHONY: validate all install_mm $(APPS)

validate all install_mm: $(APPS)

$(APPS):
	@echo "==> $(MAKECMDGOALS) in $@"
	@$(MAKE) -C $@ $(MAKECMDGOALS)
