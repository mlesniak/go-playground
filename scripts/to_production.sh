#!/bin/bash

git checkout production
git rebase master
git push
git checkout master
