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
    prerm() {
      if systemctl --quiet is-active prometheus-lvm-exporter.service; then
        echo 'Stopping prometheus-lvm-exporter.service' >&2
        systemctl stop prometheus-lvm-exporter.service || :
      fi

      echo 'Disabling prometheus-lvm-exporter.service' >&2
      systemctl disable prometheus-lvm-exporter.service || :
    }
    ;;
  *)
    echo 'Unknown init system; not managing service status' >&2
    exit 0
    ;;
esac

if [ "$1" = remove ]; then
  prerm
elif [ -n "$1" -a -z "${1##[0-9]}" ]; then
  if [ "$1" -eq 0 ]; then
    prerm
  fi
fi

exit 0

# vim: set sw=2 sts=2 et :
