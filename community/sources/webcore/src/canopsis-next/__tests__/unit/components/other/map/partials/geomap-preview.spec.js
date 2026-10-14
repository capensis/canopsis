import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { geomapPointToForm } from '@/helpers/entities/map/form';

import GeomapPreview from '@/components/other/map/partials/geomap-preview.vue';

const fitBounds = vi.fn();

const stubs = {
  geomap: {
    name: 'Geomap',
    inheritAttrs: false,
    props: ['minZoom', 'disabled'],
    setup(props, { expose }) {
      expose({
        mapObject: {
          fitBounds,
        },
      });

      return {};
    },
    template: '<div class="geomap"><slot /></div>',
  },
  'geomap-control-zoom': true,
  'geomap-control-layers': true,
  'geomap-tile-layer': true,
  'geomap-control': true,
  'geomap-cluster-group': true,
  'geomap-marker': {
    name: 'GeomapMarker',
    emits: ['click'],
    template: '<div><slot /></div>',
  },
  'geomap-icon': true,
  'point-popup-dialog': true,
  'point-icon': true,
  'c-help-icon': true,
  'c-zoom-overlay': true,
};

const selectPointMarkers = wrapper => wrapper.findAllComponents({ name: 'GeomapMarker' });
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);
const selectPointPopupDialog = wrapper => wrapper.find('point-popup-dialog-stub');

const triggerPointClick = async (wrapper, index, rect = { left: 0, top: 0, width: 0 }) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  const getBoundingClientRect = vi
    .spyOn(pointMarker.element, 'getBoundingClientRect')
    .mockReturnValue(rect);

  pointMarker.vm.$emit('click', {
    originalEvent: {
      target: pointMarker.element,
    },
  });
  await wrapper.vm.$nextTick();

  getBoundingClientRect.mockClear();

  return rect;
};

describe('geomap-preview', () => {
  const factory = generateShallowRenderer(GeomapPreview, { stubs });

  beforeEach(() => {
    fitBounds.mockClear();
  });

  test('Show alarms emitted after trigger point popup dialog', async () => {
    const point = geomapPointToForm();

    const wrapper = factory({
      props: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.triggerCustomEvent('show:alarms');

    expect(wrapper).toEmit('show:alarms', point);
  });

  test('Map fitted when points updated', async () => {
    const point = geomapPointToForm({
      coordinates: {
        lat: 1,
        lng: 2,
      },
    });

    const wrapper = factory({
      props: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await flushPromises();

    expect(fitBounds).toHaveBeenCalledWith({
      _northEast: point.coordinates,
      _southWest: point.coordinates,
    });

    const secondPoint = geomapPointToForm({
      coordinates: {
        lat: 3,
        lng: 4,
      },
    });

    const newMap = {
      parameters: {
        points: [point, secondPoint],
      },
    };

    fitBounds.mockClear();

    await wrapper.setProps({
      map: newMap,
    });

    await flushPromises();

    expect(fitBounds).toHaveBeenCalledWith({
      _northEast: secondPoint.coordinates,
      _southWest: point.coordinates,
    });
  });

  test('Show map emitted after trigger point popup dialog', async () => {
    const point = geomapPointToForm({ map: 'map' });

    const wrapper = factory({
      props: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.triggerCustomEvent('show:map');

    expect(wrapper).toEmit('show:map', point.map);
  });

  test('Point popup dialog closed after trigger close event', async () => {
    const point = geomapPointToForm({ map: 'map' });

    const wrapper = factory({
      props: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    await pointPopupDialog.triggerCustomEvent('close');

    expect(selectPointPopupDialog(wrapper).exists()).toBe(false);
  });
});
