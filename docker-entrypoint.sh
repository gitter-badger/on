#!/bin/sh
set -e

# You can set CONTINUUL_BIND_INTERFACE to the name of the interface you'd like to
# bind to and this will look up the IP and pass the proper -bind= option along
# to Continuul.
CONTINUUL_BIND=
if [ -n "$CONTINUUL_BIND_INTERFACE" ]; then
  CONTINUUL_BIND_ADDRESS=$(ip -o -4 addr list $CONTINUUL_BIND_INTERFACE | head -n1 | awk '{print $4}' | cut -d/ -f1)
  if [ -z "$CONTINUUL_BIND_ADDRESS" ]; then
    echo "Could not find IP for interface '$CONTINUUL_BIND_INTERFACE', exiting"
    exit 1
  fi

  CONTINUUL_BIND="--bind=$CONTINUUL_BIND_ADDRESS"
  echo "==> Found address '$CONTINUUL_BIND_ADDRESS' for interface '$CONTINUUL_BIND_INTERFACE', setting bind option..."
fi
if [ -n "$CONTINUUL_BIND_ADDRESS" ]; then
  CONTINUUL_BIND="--bind=$CONTINUUL_BIND_ADDRESS"
fi

CONTINUUL_DATA_DIR=/var/lib/continuul
CONTINUUL_CONFIG_DIR=/etc/continuul

if [ -n "$CONTINUUL_LOCAL_CONFIG" ]; then
        echo "$CONTINUUL_LOCAL_CONFIG" > "$CONTINUUL_CONFIG_DIR/local.json"
fi

# If the user is trying to run Continuul directly with some arguments, then
# pass them to `on`.
if [ "${1:0:1}" = '-' ]; then
    set -- on "$@"
fi

if [ "$1" = 'agent' ]; then
    shift
    set -- on agent \
        $CONTINUUL_BIND \
        "$@"
elif [ "$1" = 'version' ]; then
    # This needs a special case because there's no help output.
    set -- on "$@"
elif on --help "$1" 2>&1 | grep -q "on $1"; then
    # We can't use the return code to check for the existence of a subcommand, so
    # we have to use grep to look for a pattern in the help output.
    set -- on "$@"
fi

# If we are running Continuul, make sure it executes as the proper user.
if [ "$1" = 'on' ]; then
    # If the data or config dirs are bind mounted then chown them.
    # Note: This checks for root ownership as that's the most common case.
    if [ "$(stat -c %u /var/lib/continuul)" != "$(id -u continuul)" ]; then
        chown continuul:continuul /var/lib/continuul
    fi
    if [ "$(stat -c %u /etc/continuul)" != "$(id -u continuul)" ]; then
        chown continuul:continuul /etc/continuul
    fi

    # If requested, set the capability to bind to privileged ports before
    # we drop to the non-root user. Note that this doesn't work with all
    # storage drivers (it won't work with AUFS).
    if [ ! -z ${CONTINUUL_ALLOW_PRIVILEGED_PORTS+x} ]; then
        setcap "cap_net_bind_service=+ep" /usr/local/bin/on
    fi

    set -- gosu continuul "$@"
fi

exec "$@"
