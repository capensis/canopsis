import { generateShallowRenderer } from '@unit/utils/vue';

import { mermaidPointToForm } from '@/helpers/entities/map/form';

import MermaidPointsPreview from '@/components/other/map/partials/mermaid-points-preview.vue';

const stubs = {
  'mermaid-point-marker': {
    name: 'MermaidPointMarker',
    emits: ['click'],
    template: '<div />',
  },
  'point-popup-dialog': {
    name: 'PointPopupDialog',
    emits: ['show:alarms', 'show:map', 'close'],
    template: '<div />',
  },
};

const selectPointMarkers = wrapper => wrapper.findAllComponents({ name: 'MermaidPointMarker' });
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);
const selectPointPopupDialog = wrapper => wrapper.findComponent({ name: 'PointPopupDialog' });

const triggerPointClick = async (wrapper, index, rect = { left: 0, top: 0, width: 0 }) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  const getBoundingClientRect = vi
    .spyOn(pointMarker.element, 'getBoundingClientRect')
    .mockReturnValue(rect);

  pointMarker.vm.$emit('click', {
    target: pointMarker.element,
  });
  await wrapper.vm.$nextTick();

  getBoundingClientRect.mockClear();

  return rect;
};

describe('mermaid-points-preview', () => {
  const factory = generateShallowRenderer(MermaidPointsPreview, { stubs });

  test('Show alarms emitted after trigger point popup dialog', async () => {
    const point = mermaidPointToForm({ x: 1, y: 2 });

    const wrapper = factory({
      props: {
        points: [point],
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.vm.$emit('show:alarms');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('show:alarms')).toEqual([[point]]);
  });

  test('Show map emitted after trigger point popup dialog', async () => {
    const point = mermaidPointToForm({ x: 1, y: 2, map: 'map' });

    const wrapper = factory({
      props: {
        points: [point],
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.vm.$emit('show:map');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('show:map')).toEqual([[point.map]]);
  });

  test('Point popup dialog closed after trigger close event', async () => {
    const point = mermaidPointToForm({ x: 1, y: 2, map: 'map' });

    const wrapper = factory({
      props: {
        points: [point],
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.vm.$emit('close');
    await wrapper.vm.$nextTick();

    expect(selectPointPopupDialog(wrapper).exists()).toBe(false);
  });
});
