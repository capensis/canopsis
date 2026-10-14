import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import { mermaidPointToForm } from '@/helpers/entities/map/form';

import MermaidPointsPreview from '@/components/other/map/partials/mermaid-points-preview.vue';

const stubs = {
  'mermaid-point-marker': true,
  'point-popup-dialog': true,
};

const selectPointMarkers = wrapper => wrapper.findAll('mermaid-point-marker-stub');
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);
const selectPointPopupDialog = wrapper => wrapper.find('point-popup-dialog-stub');

const triggerPointClick = async (wrapper, index, rect = { left: 0, top: 0, width: 0 }) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  const getBoundingClientRect = jest.spyOn(pointMarker.element, 'getBoundingClientRect')
    .mockReturnValue(rect);

  await pointMarker.triggerCustomEvent('click', {
    target: pointMarker.element,
  });

  getBoundingClientRect.mockClear();

  return rect;
};

describe('mermaid-points-preview', () => {
  const factory = generateShallowRenderer(MermaidPointsPreview, { stubs });
  const snapshotFactory = generateRenderer(MermaidPointsPreview, { stubs });

  test('Show alarms emitted after trigger point popup dialog', async () => {
    const point = mermaidPointToForm({ x: 1, y: 2 });

    const wrapper = factory({
      propsData: {
        points: [point],
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.triggerCustomEvent('show:alarms');

    expect(wrapper).toEmit('show:alarms', point);
  });

  test('Show map emitted after trigger point popup dialog', async () => {
    const point = mermaidPointToForm({ x: 1, y: 2, map: 'map' });

    const wrapper = factory({
      propsData: {
        points: [point],
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.triggerCustomEvent('show:map');

    expect(wrapper).toEmit('show:map', point.map);
  });

  test('Point popup dialog closed after trigger close event', async () => {
    const point = mermaidPointToForm({ x: 1, y: 2, map: 'map' });

    const wrapper = factory({
      propsData: {
        points: [point],
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    await pointPopupDialog.triggerCustomEvent('close');

    expect(selectPointPopupDialog(wrapper).element).toBeFalsy();
  });

  test('Renders `mermaid-points-preview` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        points: [mermaidPointToForm({ x: 1, y: 2 })],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mermaid-points-preview` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        points: [mermaidPointToForm({ x: 1, y: 2 })],
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
