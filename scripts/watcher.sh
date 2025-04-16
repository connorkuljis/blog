#!/bin/bash

find static/ templates/ cmd/ internal/ | entr make site

