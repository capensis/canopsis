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

  test('Previews a column resize without updating reactive column widths', () => {
    const applyColumnResizePreview = vi.fn();
    const context = {
      aggregatedMovementDiff: 24,
      resizingColumnWidth: 100,
      resizingTableWidth: 300,
      headers: [{ value: 'column' }],
      minColumnWidth: 40,
      getNormalizedWidth: widgetColumnResizingAlarmMixin.methods.getNormalizedWidth,
      applyColumnResizePreview,
    };

    widgetColumnResizingAlarmMixin.methods.resizeColumnByDiff.call(context, 0);

    expect(context.resizingColumnWidth).toBe(124);
    expect(context.resizingTableWidth).toBe(324);
    expect(context.aggregatedMovementDiff).toBe(0);
    expect(context.columnsWidthByField).toBeUndefined();
    expect(applyColumnResizePreview).toHaveBeenCalledWith(0, 124, 324);
  });

  test('Commits the previewed column width when resizing finishes', () => {
    const context = {
      aggregatedMovementDiff: 0,
      resizingColumnIndex: 0,
      resizingColumnWidth: 124,
      resizingTableWidth: 324,
      columnsWidthByField: { column: 100, actions: 200 },
      headers: [{ value: 'column' }, { value: 'actions' }],
      draggingMode: false,
      flushColumnResize: vi.fn(),
      clearColumnResizePreview: vi.fn(),
      handleColumnResize: vi.fn(),
      finishColumnResize: vi.fn(),
      suppressNextClickAfterResize: vi.fn(),
    };

    widgetColumnResizingAlarmMixin.methods.finishColumnResize.call(context, {
      suppressClick: false,
    });

    expect(context.flushColumnResize).toHaveBeenCalledOnce();
    expect(context.columnsWidthByField).toEqual({ column: 124, actions: 200 });
    expect(context.clearColumnResizePreview).toHaveBeenCalledWith(0);
    expect(context.resizingColumnIndex).toBeNull();
    expect(context.resizingColumnWidth).toBeNull();
    expect(context.resizingTableWidth).toBeNull();
  });
});
