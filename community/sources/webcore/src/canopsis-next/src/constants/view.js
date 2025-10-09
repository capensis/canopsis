export const VIEW_SCREEN_MODES = {
  default: 'default',
  fullscreen: 'fullscreen',
  kiosk: 'kiosk',
  kioskFullscreen: 'kioskFullscreen',
};

export const FULLSCREEN_MODES_TO_DEFAULT_SCREEN_MODES = {
  [VIEW_SCREEN_MODES.fullscreen]: VIEW_SCREEN_MODES.default,
  [VIEW_SCREEN_MODES.kioskFullscreen]: VIEW_SCREEN_MODES.kiosk,
};

export const KEYS_TO_VIEW_SCREEN_MODES = {
  Digit1: VIEW_SCREEN_MODES.default,
  Digit2: VIEW_SCREEN_MODES.fullscreen,
  Digit3: VIEW_SCREEN_MODES.kiosk,
  Digit4: VIEW_SCREEN_MODES.kioskFullscreen,
};
