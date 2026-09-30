import { widgetColumnResizingAlarmMixin } from '@/mixins/widget/columns/alarm-resizing';

describe('widget-column-resizing-alarm', () => {
  test('Calculates a column width without the element transform', () => {
    const element = document.createElement('th');

    Object.defineProperty(element, 'offsetWidth', { value: 240 });
    vi.spyOn(element, 'getBoundingClientRect').mockReturnValue({ width: 96 });

    const context = {
      minColumnWidth: 40,
      getNormalizedWidth: widgetColumnResizingAlarmMixin.methods.getNormalizedWidth,
    };

    expect(
      widgetColumnResizingAlarmMixin.methods.calculateElementNormalizedWidth.call(
        context,
        element,
        'column',
      ),
    ).toBe(240);
  });
});
