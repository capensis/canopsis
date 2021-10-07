export const REQUEST_MESSAGES_TYPES = {
  join: 0,
  leave: 1,
  authenticate: 2,
  heartbeat: 3,
};

export const RESPONSE_MESSAGES_TYPES = {
  ok: 0,
  error: 1,
  close: 2,
};

export const MAX_RECONNECTS_COUNT = 10;

export const HEARTBEAT_TIMEOUT = 5000;
