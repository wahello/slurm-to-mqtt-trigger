BINDIR	= $(CURDIR)/bin

depend:
	# handled by go mod

build:
	cd $(CURDIR)/src/slurm-to-mqtt-trigger && go get slurm-to-mqtt-trigger/src/slurm-to-mqtt-trigger && go build -o $(CURDIR)/bin/slurm-to-mqtt-trigger

destdirs:
	mkdir -p -m 0755 $(DESTDIR)/usr/bin

strip: build
	strip --strip-all $(BINDIR)/slurm-to-mqtt-trigger

ifneq (, $(shell which upx 2>/dev/null))
	upx -9 $(BINDIR)/slurm-to-mqtt-trigger
endif

install: strip destdirs install-bin

install-bin:
	install -m 0755 $(BINDIR)/slurm-to-mqtt-trigger $(DESTDIR)/usr/bin

clean:
	/bin/rm -f bin/slurm-to-mqtt-trigger

distclean: clean

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: depend build strip install

