#!/bin/bash

tcpdump -i eth0 -e -l -v > dump.log
