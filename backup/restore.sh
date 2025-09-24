#!/bin/bash

#pg_restore -U emre -d bitirme /backup/backup.dump
pg_restore -U postgres -d dvdrental /backup/dvdrental.tar


