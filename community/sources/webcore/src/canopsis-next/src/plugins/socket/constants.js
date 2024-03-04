export const REQUEST_MESSAGES_TYPES = {
  ping: 0,
  join: 1,
  leave: 2,
  authenticate: 3,
};

export const RESPONSE_MESSAGES_TYPES = {
  pong: 0,
  ok: 1,
  error: 2,
  close: 3,
  authenticated: 4,
};

export const ICONS_RESPONSE_MESSAGES_TYPES = {
  create: 0,
  update: 1,
  delete: 2,
};

export const MAX_RECONNECTS_COUNT = 10;

export const PING_INTERVAL = 5000;

export const RECONNECT_INTERVAL = 5000;

export const EVENTS_TYPES = {
  customClose: 'custom-close',
  closeRoom: 'close-room',
  networkError: 'network-error',
};

export const ERROR_MESSAGES = {
  authenticationFailed: 'authentication failed',
};
