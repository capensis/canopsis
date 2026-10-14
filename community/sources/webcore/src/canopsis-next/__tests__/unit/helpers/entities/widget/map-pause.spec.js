import { mapWidgetParametersToForm } from '@/helpers/entities/widget/forms/map';

describe('map pause settings', () => {
  test('enables pause display by default', () => {
    expect(mapWidgetParametersToForm().entities_under_pause_enabled).toBe(true);
  });

  test('preserves disabled pause display', () => {
    expect(
      mapWidgetParametersToForm({
        entities_under_pause_enabled: false,
      }).entities_under_pause_enabled,
    ).toBe(false);
  });
});
