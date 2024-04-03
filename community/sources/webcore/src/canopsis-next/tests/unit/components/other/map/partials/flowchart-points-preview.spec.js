import { keyBy } from 'lodash';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES, SHAPES } from '@/constants';

import { shapeToForm } from '@/helpers/flowchart/shapes';

import FlowchartPointsPreview from '@/components/other/map/partials/flowchart-points-preview.vue';

const stubs = {
  'point-icon': true,
  'point-popup-dialog': true,
};

const selectPoints = wrapper => wrapper.findAll('.flowchart-points-preview__point');
const selectPointByIndex = (wrapper, index) => selectPoints(wrapper).at(index);
const selectPointPopupDialog = wrapper => wrapper.find('point-popup-dialog-stub');

const triggerPointClick = async (wrapper, index, rect = { left: 0, top: 0, width: 0 }) => {
  const point = selectPointByIndex(wrapper, index);

  const getBoundingClientRect = jest.spyOn(point.element, 'getBoundingClientRect')
    .mockReturnValue(rect);

  await point.trigger('click');

  getBoundingClientRect.mockClear();

  return rect;
};

describe('flowchart-points-preview', () => {
  const rectShape = shapeToForm({ type: SHAPES.rect, _id: 'rect' });
  const lineShape = shapeToForm({ type: SHAPES.rect, _id: 'line' });
  const shapes = keyBy([rectShape, lineShape], '_id');

  const firstPoint = { _id: 'point-1', x: 12, y: 32 };
  const secondPoint = {
    _id: 'point-2',
    shape: lineShape._id,
    entity: {},
  };
  const points = [firstPoint, secondPoint];

  const factory = generateShallowRenderer(FlowchartPointsPreview, { stubs });
  const snapshotFactory = generateRenderer(FlowchartPointsPreview, { stubs });

  test('Show alarms emitted after trigger point popup dialog', async () => {
    const wrapper = factory({
      propsData: {
        shapes,
        points,
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.triggerCustomEvent('show:alarms');

    expect(wrapper).toEmit('show:alarms', firstPoint);
  });

  test('Show map emitted after trigger point popup dialog', async () => {
    const wrapper = factory({
      propsData: {
        shapes,
        points,
      },
    });

    await triggerPointClick(wrapper, 1);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.triggerCustomEvent('show:map');

    expect(wrapper).toEmit('show:map', secondPoint.map);
  });

  test('Point popup dialog closed after trigger close event', async () => {
    const wrapper = factory({
      propsData: {
        shapes,
        points,
      },
    });

    await triggerPointClick(wrapper, 1);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    await pointPopupDialog.triggerCustomEvent('close');

    expect(selectPointPopupDialog(wrapper).element).toBeFalsy();
  });

  test('Renders `flowchart-points-preview` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        points,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-points-preview` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        points,
        iconSize: 12,
        popupTemplate: '{{ entity.name }}',
        popupActions: true,
        colorIndicator: COLOR_INDICATOR_TYPES.impactState,
        pbehaviorEnabled: true,
      },
    });

    await triggerPointClick(wrapper, 0, {
      top: 1,
      left: 3,
      width: 5,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
