import { ROUTES_NAMES, VIEW_SCREEN_MODES } from '@/constants';

import { shownHeaderMixin } from '@/mixins/shown-header';

const computeHeaderVisibility = ({
  screenMode = VIEW_SCREEN_MODES.default,
  route = { meta: {} },
  showHeaderOnKioskMode = false,
} = {}) => {
  const context = {
    activeViewScreenMode: screenMode,
    activeViewIsKioskScreenMode: [
      VIEW_SCREEN_MODES.kiosk,
      VIEW_SCREEN_MODES.kioskFullscreen,
    ].includes(screenMode),
    $route: route,
    showHeaderOnKioskMode,
  };

  context.isFullscreenMode = shownHeaderMixin.computed.isFullscreenMode.call(context);
  context.isKioskMode = shownHeaderMixin.computed.isKioskMode.call(context);

  return shownHeaderMixin.computed.shownHeader.call(context);
};

describe('shown-header mixin', () => {
  test.each([VIEW_SCREEN_MODES.fullscreen, VIEW_SCREEN_MODES.kioskFullscreen])(
    'hides the header in %s mode',
    (screenMode) => {
      expect(computeHeaderVisibility({ screenMode, showHeaderOnKioskMode: true })).toBe(false);
    },
  );

  test('uses the kiosk header setting in kiosk mode', () => {
    expect(
      computeHeaderVisibility({
        screenMode: VIEW_SCREEN_MODES.kiosk,
        showHeaderOnKioskMode: true,
      }),
    ).toBe(true);
    expect(
      computeHeaderVisibility({
        screenMode: VIEW_SCREEN_MODES.kiosk,
        showHeaderOnKioskMode: false,
      }),
    ).toBe(false);
  });

  test('uses the kiosk header setting on the kiosk view route', () => {
    expect(
      computeHeaderVisibility({
        route: { name: ROUTES_NAMES.viewKiosk, meta: {} },
        showHeaderOnKioskMode: true,
      }),
    ).toBe(true);
  });

  test('uses route metadata in default mode', () => {
    expect(computeHeaderVisibility()).toBe(true);
    expect(computeHeaderVisibility({ route: { meta: { hideHeader: true } } })).toBe(false);
  });
});
