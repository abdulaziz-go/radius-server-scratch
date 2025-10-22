#!/bin/bash
set -Eeuo pipefail

REDIS_HOST=${REDIS_HOST:-127.0.0.1}
REDIS_PORT=${REDIS_PORT:-6379}
REDIS_CONNECT_INTERVAL=1   # seconds between PINGs
INDEX_CREATE_INTERVAL=2    # seconds between FT.CREATE retries
INDEX_NAME_SUBS=index_subs
INDEX_NAME_NAS=index_nas
###############################################################################

create_subscriber_index() {
  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" FT.CREATE "$INDEX_NAME_SUBS" ON HASH PREFIX 1 "subscriber:" SCHEMA \
    subscriber_id NUMERIC SORTABLE \
    ip TEXT \
    ip_version TAG \
    session_id TEXT \
    last_updated_time NUMERIC SORTABLE
}

create_nas_index() {
  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" FT.CREATE "$INDEX_NAME_NAS" ON HASH PREFIX 1 "radius_nas:" SCHEMA \
    ip_address TAG \
}

index_exists() {
  local index_name=$1
  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" --raw FT._LIST | grep -Fxq "$index_name"
}

create_indexes_if_missing() {
  for index_name in "$INDEX_NAME_SUBS" "$INDEX_NAME_NAS"; do
    while true; do
      if index_exists "$index_name"; then
        echo "[Redis] Index '$index_name' already exists."
        break
      fi
      echo "[Redis] Attempting to create index '$index_name' ..."
      if [[ "$index_name" == "$INDEX_NAME_SUBS" ]]; then
        create_subscriber_index
      else
        create_nas_index
      fi
      echo "[Redis] Indexes created successfully."
      break
    done
  done
}

wait_for_redis_and_create_indexes() {
  echo -n "[Redis] Waiting for Redis to accept connections"
  until redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" PING &>/dev/null; do
    echo -n "."
    sleep "$REDIS_CONNECT_INTERVAL"
  done
  echo " up!"

  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" SET sanity_test ok >/dev/null

  create_indexes_if_missing
}

# Run in background
(wait_for_redis_and_create_indexes)&

echo "[Redis] Starting Redis Stack ..."
cd /data
exec /usr/local/bin/docker-entrypoint.sh /usr/local/lib/redis.conf