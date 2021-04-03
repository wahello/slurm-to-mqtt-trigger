GOPATH	= $(CURDIR)
BINDIR	= $(CURDIR)/bin

depend:
	env GOPATH=$(GOPATH) go get -u github.com/eclipse/paho.mqtt.golang
	env GOPATH=$(GOPATH) go get -u github.com/sirupsen/logrus
	env GOPATH=$(GOPATH) go get -u gopkg.in/ini.v1

build:
	env GOPATH=$(GOPATH) go install slurm-to-mqtt-trigger

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
	rm -rf src/github.com/
	rm -rf src/golang.org/
	rm -rf src/gopkg.in/
	rm -rf pkg/

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: depend build strip install

