# Fans `make <target>` out to every subfolder that has its own makefile.
# Adding a makefile to an example is all it takes to wire it into `make validate`.

APPS := $(patsubst %/makefile,%,$(wildcard */makefile))

.PHONY: validate all install_mm $(APPS)

validate all install_mm: $(APPS)

$(APPS):
	@echo "==> $(MAKECMDGOALS) in $@"
	@$(MAKE) -C $@ $(MAKECMDGOALS)
