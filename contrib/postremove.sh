#!/bin/sh

# References:
# https://www.debian.org/doc/debian-policy/ch-maintainerscripts.html
# https://docs.fedoraproject.org/en-US/packaging-guidelines/Scriptlets/

init_system=

if [ -e /proc/1/comm ]; then
  init_system=$(cat /proc/1/comm)
fi

case "$init_system" in
  systemd)
    systemctl daemon-reload || :
    ;;
  *)
    echo 'Unknown init system; not managing service status' >&2
    exit 0
    ;;
esac

exit 0

# vim: set sw=2 sts=2 et :
