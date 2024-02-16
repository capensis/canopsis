import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import GeomapPreview from '@/components/other/map/partials/geomap-preview.vue';
import { COLOR_INDICATOR_TYPES, ENTITIES_STATES, PBEHAVIOR_TYPE_TYPES } from '@/constants';
import { geomapPointToForm } from '@/helpers/entities/map/form';

const fitBounds = jest.fn();

const stubs = {
  geomap: {
    template: '<div v-on="$listeners" class="geomap"><slot/></div>',
    computed: {
      mapObject() {
        return {
          fitBounds,
        };
      },
    },
  },
  'geomap-control-zoom': true,
  'geomap-control-layers': true,
  'geomap-tile-layer': true,
  'geomap-control': true,
  'geomap-cluster-group': true,
  'geomap-marker': true,
  'geomap-icon': true,
  'point-popup-dialog': true,
  'point-icon': true,
  'c-help-icon': true,
  'c-zoom-overlay': true,
};

const selectPointMarkers = wrapper => wrapper.findAll('geomap-marker-stub');
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);
const selectPointPopupDialog = wrapper => wrapper.find('point-popup-dialog-stub');

const triggerPointClick = async (wrapper, index, rect = { left: 0, top: 0, width: 0 }) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  const getBoundingClientRect = jest.spyOn(pointMarker.element, 'getBoundingClientRect')
    .mockReturnValue(rect);

  await pointMarker.vm.$emit('click', {
    originalEvent: {
      target: pointMarker.element,
    },
  });

  getBoundingClientRect.mockClear();

  return rect;
};

describe('geomap-preview', () => {
  const points = [
    {
      _id: 'point-1',
      coordinates: { lat: 1, lng: 1 },
      entity: {
        state: ENTITIES_STATES.minor,
        impact_state: 1,
      },
    },
    {
      _id: 'point-2',
      coordinates: { lat: 2, lng: 2 },
      entity: {
        state: ENTITIES_STATES.critical,
        impact_state: 2,
      },
    },
    {
      _id: 'point-3',
      coordinates: { lat: 3, lng: 3 },
      map: {
        _id: 'map-1',
      },
    },
    {
      _id: 'point-4',
      coordinates: { lat: 4, lng: 4 },
      entity: {
        state: ENTITIES_STATES.ok,
        impact_state: 2,
        pbehavior_info: {
          canonical_type: PBEHAVIOR_TYPE_TYPES.active,
        },
      },
    },
  ];

  const factory = generateShallowRenderer(GeomapPreview, { stubs });
  const snapshotFactory = generateRenderer(GeomapPreview, { stubs });

  beforeEach(() => {
    fitBounds.mockClear();
  });

  test('Show alarms emitted after trigger point popup dialog', async () => {
    const point = geomapPointToForm();

    const wrapper = factory({
      propsData: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.vm.$emit('show:alarms');

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
      propsData: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await flushPromises();

    expect(fitBounds).toBeCalledWith({
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

    expect(fitBounds).toBeCalledWith({
      _northEast: secondPoint.coordinates,
      _southWest: point.coordinates,
    });
  });

  test('Show map emitted after trigger point popup dialog', async () => {
    const point = geomapPointToForm({ map: 'map' });

    const wrapper = factory({
      propsData: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    pointPopupDialog.vm.$emit('show:map');

    expect(wrapper).toEmit('show:map', point.map);
  });

  test('Point popup dialog closed after trigger close event', async () => {
    const point = geomapPointToForm({ map: 'map' });

    const wrapper = factory({
      propsData: {
        map: {
          parameters: {
            points: [point],
          },
        },
      },
    });

    await triggerPointClick(wrapper, 0);

    const pointPopupDialog = selectPointPopupDialog(wrapper);
    await pointPopupDialog.vm.$emit('close');

    expect(selectPointPopupDialog(wrapper).element).toBeFalsy();
  });

  test('Renders `geomap-preview` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        map: {
          parameters: {
            points: [],
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `geomap-preview` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        map: {
          parameters: {
            points: [geomapPointToForm()],
          },
        },
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

  test('Renders `geomap-preview` with all point types and color indicator impact state', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        map: {
          parameters: {
            points,
          },
        },
        colorIndicator: COLOR_INDICATOR_TYPES.impactState,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `geomap-preview` with all point types and color indicator state', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        map: {
          parameters: {
            points,
          },
        },
        colorIndicator: COLOR_INDICATOR_TYPES.state,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `geomap-preview` with all point types and pbehavior enabled', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        map: {
          parameters: {
            points,
          },
        },
        pbehaviorEnabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
