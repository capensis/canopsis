import { TOP_BAR_HEIGHT } from '@/config';

import { widgetStickyAlarmMixin } from '@/mixins/widget/rows/alarm-sticky';

const computeTopBarHeight = ({
  shownHeader = true,
  isFullscreenMode = false,
  isKioskMode = false,
} = {}) =>
  widgetStickyAlarmMixin.computed.topBarHeight.call({
    shownHeader,
    isFullscreenMode,
    isKioskMode,
  });

describe('widget-sticky-alarm mixin', () => {
  test('uses the toolbar height in the default view mode', () => {
    expect(computeTopBarHeight()).toBe(TOP_BAR_HEIGHT);
  });

  test.each([
    { isFullscreenMode: true, mode: 'fullscreen' },
    { isKioskMode: true, mode: 'kiosk' },
  ])('does not add top spacing in $mode mode', (mode) => {
    expect(computeTopBarHeight(mode)).toBe(0);
  });

  test('does not add top spacing when the toolbar is hidden', () => {
    expect(computeTopBarHeight({ shownHeader: false })).toBe(0);
  });

  test('recalculates sticky positions after view screen mode changes', () => {
    const scheduleStickyScrollHandler = vi.fn();

    widgetStickyAlarmMixin.watch.activeViewScreenMode.call({
      someSticky: true,
      scheduleStickyScrollHandler,
    });

    expect(scheduleStickyScrollHandler).toHaveBeenCalledOnce();
  });

  test('does not recalculate positions without sticky elements', () => {
    const scheduleStickyScrollHandler = vi.fn();

    widgetStickyAlarmMixin.watch.activeViewScreenMode.call({
      someSticky: false,
      scheduleStickyScrollHandler,
    });

    expect(scheduleStickyScrollHandler).not.toHaveBeenCalled();
  });
});
