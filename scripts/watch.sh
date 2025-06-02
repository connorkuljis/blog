#!/bin/bash

find assets/ templates/ cmd/ internal/ | entr make site

