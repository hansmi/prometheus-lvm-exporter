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

    postinstall() {
      echo 'Enabling prometheus-lvm-exporter.service' >&2
      systemctl enable prometheus-lvm-exporter.service || :
      echo 'Starting prometheus-lvm-exporter.service' >&2
      systemctl start prometheus-lvm-exporter.service || :
    }

    postupgrade() {
      echo 'Restarting prometheus-lvm-exporter.service' >&2
      systemctl try-restart prometheus-lvm-exporter.service || :
    }
  ;;
  *)
    echo 'Unknown init system; not managing service status' >&2
    exit 0
    ;;
esac

if [ "$1" = configure ]; then
  if [ -z "$2" ]; then
    postinstall
  else
    postupgrade
  fi
elif [ -n "$1" -a -z "${1##[0-9]}" ]; then
  if [ "$1" -eq 1 ]; then
    postinstall
  elif [ "$1" -ge 1 ]; then
    postupgrade
  fi
fi

exit 0

# vim: set sw=2 sts=2 et :
