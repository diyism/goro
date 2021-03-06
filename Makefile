#!/bin/make

all:
	goimports -w -l .
	go build -v

sapi/php-cli/php-cli:
	make -C sapi/php-cli

sapi/php-httpd/php-httpd:
	make -C sapi/php-httpd

sapi: sapi/php-cli/php-cli sapi/php-httpd/php-httpd

test:
	go test

http:
	$(RM) sapi/php-httpd/php-httpd
	make sapi/php-httpd/php-httpd
	./sapi/php-httpd/php-httpd

buildext:
	$(RM) sapi/php-cli/php-cli
	make sapi/php-cli/php-cli
	./sapi/php-cli/php-cli tools/buildext.php

deps:
	go get -v .

testdeps:
	go get -v -t .

travis:
	# get deps
	go get -v -t .
	go build -v
	make -C sapi/php-cgi
	make -C sapi/php-cli
	make -C sapi/php-fpm
	make -C sapi/php-httpd

.PHONY: test http buildext sapi testdeps deps travis
