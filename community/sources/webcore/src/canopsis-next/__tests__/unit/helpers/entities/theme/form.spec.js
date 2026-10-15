import { COLORS } from '@/config';
import { DEFAULT_WIDGET_BORDER_RADIUS } from '@/constants';

import { formToTheme, themeIconsColorsToForm, themeToForm } from '@/helpers/entities/theme/form';

describe('themeIconsColorsToForm', () => {
  it('uses the default extra details colors', () => {
    expect(themeIconsColorsToForm()).toEqual(COLORS.alarmExtraDetails);
  });

  it('preserves configured extra details colors', () => {
    const icons = {
      ack: '#111111',
      canceled: '#222222',
      children: '#333333',
      parent: '#444444',
      comment: '#555555',
      snooze: '#666666',
      tickets: '#777777',
      cancel_delay: '#888888',
    };

    expect(themeIconsColorsToForm(icons)).toEqual(icons);
  });

  it('adds defaults to themes created before extra details customization', () => {
    expect(themeToForm({ colors: {} }).colors.icons).toEqual(COLORS.alarmExtraDetails);
  });

  it('adds required widget colors to the API payload when they are missing', () => {
    const theme = formToTheme({ name: 'Theme', font_size: 2, colors: {} });

    expect(theme.colors.widget).toEqual({
      title_background: COLORS.widget.titleBackground,
      content_background: COLORS.widget.contentBackground,
      border: false,
      border_color: COLORS.widget.borderColor,
      border_radius: DEFAULT_WIDGET_BORDER_RADIUS,
    });
    expect(theme.colors.icons).toEqual(COLORS.alarmExtraDetails);
  });

  it('replaces blank required widget colors in the API payload', () => {
    const theme = formToTheme({
      name: 'Theme',
      font_size: 2,
      colors: {
        widget: {
          title_background: '',
          content_background: '',
        },
      },
    });

    expect(theme.colors.widget.title_background).toBe(COLORS.widget.titleBackground);
    expect(theme.colors.widget.content_background).toBe(COLORS.widget.contentBackground);
  });
});
