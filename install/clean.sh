#!/bin/bash

echo -e "input\noutput\nerror"|xargs -n 1 -I {} find /app/judge/src/workdir/users/ -name "{}*"|xargs rm -rf
