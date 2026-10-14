export const CLOCK_DISPLAY_TYPES = {
  digital: 'digital',
  sevenSegments: 'seven_segments',
  analog: 'analog',
};

export const CLOCK_SIZES = {
  small: 's',
  medium: 'm',
  large: 'l',
};

export const CLOCK_SPACING_TYPES = {
  automatic: 'automatic',
  fixed: 'fixed',
};

export const DEFAULT_CLOCK_WIDGET_PARAMETERS = {
  display: CLOCK_DISPLAY_TYPES.digital,
  size: CLOCK_SIZES.medium,
  spacing: CLOCK_SPACING_TYPES.automatic,
  is_12_hour_format: false,
  show_seconds: true,
  show_date: true,
  show_utc_offset: true,
  show_timezone: true,
};
