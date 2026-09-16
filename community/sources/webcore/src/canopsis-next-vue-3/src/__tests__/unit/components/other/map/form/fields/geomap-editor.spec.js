import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { LatLngBounds } from 'leaflet';

import { MODALS } from '@/constants';

import { geomapPointToForm } from '@/helpers/entities/map/form';

import GeomapEditor from '@/components/other/map/form/fields/geomap-editor.vue';

const fitBounds = vi.fn();
const fitWorld = vi.fn();
const panTo = vi.fn();
const closeContextmenu = vi.fn();
const latLngToContainerPoint = vi.fn().mockReturnValue({ x: 0, y: 0 });
const getBounds = vi.fn().mockReturnValue(new LatLngBounds([0, 0], [1, 1]));

const stubs = {
  'v-tooltip': createActivatorElementStub('v-tooltip'),
  geomap: {
    name: 'Geomap',
    props: ['style', 'disabled', 'minZoom', 'options', 'center'],
    emits: ['click', 'dblclick', 'update:center'],
    setup(props, { expose }) {
      expose({
        mapObject: {
          fitBounds,
          fitWorld,
          panTo,
          latLngToContainerPoint,
        },
      });

      return {};
    },
    template: '<div class="geomap"><slot /></div>',
  },
  'geomap-control-zoom': true,
  'geomap-control-layers': true,
  'geomap-contextmenu': {
    props: ['items', 'markerItems'],
    setup(props, { expose }) {
      expose({ close: closeContextmenu });

      return {};
    },
    template: `<div
      class="geomap-contextmenu"
      :items="items"
      :marker-items="markerItems"
    />`,
  },
  'geomap-control': true,
  'geomap-tile-layer': true,
  'geomap-cluster-group': {
    setup(props, { expose }) {
      expose({
        mapObject: {
          getBounds,
        },
      });

      return {};
    },
    template: '<div class="geomap-cluster-group"><slot/></div>',
  },
  'geomap-marker': {
    name: 'GeomapMarker',
    props: ['options', 'latLng', 'draggable'],
    emits: ['dragend', 'dblclick'],
    template: '<geomap-marker-stub><slot /></geomap-marker-stub>',
  },
  'c-zoom-overlay': true,
  'c-layout': {
    template: '<div><slot /></div>',
  },
  'geomap-icon': true,
  'point-icon': true,
  'map-editor-required-field': {
    name: 'MapEditorRequiredField',
    props: ['valid', 'message', 'name'],
    template: '<div><slot :has-error="!valid" :messages="[message]" /></div>',
  },
  'point-form-dialog-menu': {
    name: 'PointFormDialogMenu',
    props: ['value', 'point', 'positionX', 'positionY', 'editing', 'existsEntities', 'coordinates'],
    emits: ['cancel', 'submit', 'remove', 'fly:coordinates'],
    template: '<point-form-dialog-menu-stub />',
  },
  'c-help-icon': true,
  'v-btn': {
    name: 'VBtn',
    emits: ['click'],
    template: '<v-btn-stub @click="$emit(\'click\', $event)"><slot /></v-btn-stub>',
  },
  'v-messages': true,
};

const selectGeomap = wrapper => wrapper.find('div.geomap');
const selectGeomapComponent = wrapper => wrapper.findComponent({ name: 'Geomap' });
const selectAddLocationBtn = wrapper => wrapper.find('v-btn-stub');
const selectMapEditorRequiredField = wrapper =>
  wrapper.findComponent({ name: 'MapEditorRequiredField' });
const selectPointFormDialogMenu = wrapper => wrapper.findComponent({ name: 'PointFormDialogMenu' });
const selectPointMarkers = wrapper => wrapper.findAllComponents({ name: 'GeomapMarker' });
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);

const getEvent = event => ({
  latlng: {
    lat: Faker.datatype.number(),
    lng: Faker.datatype.number(),
  },
  ...event,
});

const triggerPointEvent = async (wrapper, index, event) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  pointMarker.vm.$emit(event, new Event(event));

  await flushPromises();
};

const checkDialogMenuPosition = (wrapper, { x, y }) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  expect(pointFormDialog.props('value')).toBe(true);
  expect(pointFormDialog.props('positionY')).toBe(y);
  expect(pointFormDialog.props('positionX')).toBe(x);
};

const fillPointDialog = (wrapper, point = geomapPointToForm()) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  pointFormDialog.vm.$emit('submit', point);

  return point;
};

const checkMenuIsClosed = (wrapper) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  expect(pointFormDialog.props('value')).toBeFalsy();
};

const triggerContextMenuEvent = (wrapper, event, data) => {
  const item = {
    add: wrapper.vm.mapItems[0],
    edit: wrapper.vm.markerItems[0],
    remove: wrapper.vm.markerItems[1],
  }[event];

  item.action(data);
};

const triggerItemContextMenuEvent = (wrapper, index, event) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  const marker = {
    options: pointMarker.props('options'),
  };

  triggerContextMenuEvent(wrapper, event, { marker });
};

describe('geomap-editor', () => {
  const $modals = mockModals();
  const initialForm = {
    points: [],
  };

  const factory = generateShallowRenderer(GeomapEditor, { stubs });

  test('Add on click mode enabled after trigger add location btn', async () => {
    const event = getEvent();
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const geomap = selectGeomap(wrapper);
    const getBoundingClientRect = vi
      .spyOn(geomap.element, 'getBoundingClientRect')
      .mockReturnValue({ x: 12, y: 22 });
    latLngToContainerPoint.mockReturnValueOnce({
      x: 88,
      y: 28,
    });

    const addLocationBtn = selectAddLocationBtn(wrapper);

    await addLocationBtn.triggerCustomEvent('click');

    selectGeomapComponent(wrapper).vm.$emit('click', event);
    await flushPromises();

    checkDialogMenuPosition(wrapper, { x: 100, y: 50 });

    getBoundingClientRect.mockClear();
  });

  test('Point didn\'t add after click with add on click disabled', async () => {
    const wrapper = factory({
      propsData: {
        form: {
          points: [],
        },
      },
    });

    selectGeomapComponent(wrapper).vm.$emit('click');
    await flushPromises();

    const pointFormDialog = selectPointFormDialogMenu(wrapper);

    expect(pointFormDialog.props('value')).toBeFalsy();
  });

  test('Form re-validated after change form with error', async () => {
    const wrapper = factory({
      propsData: {
        name: 're-validate-name',
        form: {
          points: [],
        },
      },
    });

    expect(selectMapEditorRequiredField(wrapper).props('valid')).toBe(false);

    wrapper.setProps({
      form: {
        points: [geomapPointToForm()],
      },
    });

    await flushPromises();

    expect(selectMapEditorRequiredField(wrapper).props('valid')).toBe(true);
  });

  test('Point added after double click', async () => {
    const event = getEvent();
    const menuPosition = {
      x: Faker.datatype.number(),
      y: Faker.datatype.number(),
    };

    const wrapper = factory({
      propsData: {
        form: {
          points: [],
        },
      },
    });

    latLngToContainerPoint.mockReturnValueOnce(menuPosition);

    selectGeomapComponent(wrapper).vm.$emit('dblclick', event);
    await flushPromises();

    checkDialogMenuPosition(wrapper, menuPosition);

    const newPoint = fillPointDialog(wrapper);

    expect(wrapper).toEmitInput({
      points: [newPoint],
    });

    await flushPromises();
    await checkMenuIsClosed(wrapper);
  });

  test('Point dialog menu didn\'t changed after double click when dialog already opened', async () => {
    const event = getEvent();
    const secondEvent = getEvent();
    const firstMenuPosition = {
      x: Faker.datatype.number(),
      y: Faker.datatype.number(),
    };
    const secondMenuPosition = {
      x: Faker.datatype.number(),
      y: Faker.datatype.number(),
    };
    latLngToContainerPoint
      .mockReturnValueOnce(firstMenuPosition)
      .mockReturnValueOnce(secondMenuPosition);

    const wrapper = factory({
      propsData: {
        form: {
          points: [],
        },
      },
    });

    selectGeomapComponent(wrapper).vm.$emit('dblclick', event);
    await flushPromises();
    checkDialogMenuPosition(wrapper, firstMenuPosition);

    selectGeomapComponent(wrapper).vm.$emit('dblclick', secondEvent);
    await flushPromises();
    checkDialogMenuPosition(wrapper, firstMenuPosition);
  });

  test('Point edited after open form dialog by contextmenu', async () => {
    const point = geomapPointToForm();
    const wrapper = factory({
      propsData: {
        form: {
          points: [point],
        },
      },
    });

    await triggerItemContextMenuEvent(wrapper, 0, 'edit');
    const newPoint = fillPointDialog(wrapper, { ...point, entity: 'entity' });

    expect(wrapper).toEmitInput({ points: [newPoint] });

    await flushPromises();
    await checkMenuIsClosed(wrapper);
  });

  test('Point added after open form dialog by contextmenu', async () => {
    const point = geomapPointToForm();
    const wrapper = factory({
      propsData: {
        form: {
          points: [point],
        },
      },
    });
    const latlng = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };

    await triggerContextMenuEvent(wrapper, 'add', { latlng });
    const newPoint = fillPointDialog(wrapper, { ...point, entity: 'entity' });

    expect(wrapper).toEmitInput({ points: [point, newPoint] });

    await flushPromises();
    await checkMenuIsClosed(wrapper);
  });

  test('Point removed form point dialog', async () => {
    const point = geomapPointToForm();
    const wrapper = factory({
      propsData: {
        form: {
          points: [point],
        },
      },
      mocks: {
        $modals,
      },
    });

    await triggerItemContextMenuEvent(wrapper, 0, 'edit');

    const pointFormDialog = selectPointFormDialogMenu(wrapper);

    pointFormDialog.vm.$emit('remove');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(wrapper).toEmitInput({ points: [] });

    await flushPromises();
    await checkMenuIsClosed(wrapper);
  });

  test('Point edited after open form dialog by point double click', async () => {
    const point = geomapPointToForm();
    const secondPoint = geomapPointToForm();
    const wrapper = factory({
      propsData: {
        form: {
          points: [point, secondPoint],
        },
      },
    });

    await triggerPointEvent(wrapper, 1, 'dblclick');
    const newPoint = fillPointDialog(wrapper, { ...secondPoint, entity: 'entity' });

    expect(wrapper).toEmitInput({
      points: [point, newPoint],
    });

    await flushPromises();
    await checkMenuIsClosed(wrapper);
  });

  test('Point removed after trigger remove in contextmenu', async () => {
    const point = geomapPointToForm();
    const secondPoint = geomapPointToForm();
    const wrapper = factory({
      propsData: {
        form: {
          points: [point, secondPoint],
        },
      },
      mocks: {
        $modals,
      },
    });

    await triggerItemContextMenuEvent(wrapper, 1, 'remove');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(wrapper).toEmitInput({ points: [point] });
  });

  test('Coordinates changed after change coordinates by point dialog', async () => {
    const coordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };
    const point = geomapPointToForm();
    const wrapper = factory({
      propsData: {
        form: {
          points: [point],
        },
      },
    });

    await triggerItemContextMenuEvent(wrapper, 0, 'edit');

    const pointFormDialog = selectPointFormDialogMenu(wrapper);

    pointFormDialog.vm.$emit('fly:coordinates', coordinates);
    await flushPromises();

    expect(pointFormDialog.props('point').coordinates).toEqual(coordinates);
  });

  test('Marker moved after dragend', async () => {
    const coordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };
    const point = geomapPointToForm();
    const secondPoint = geomapPointToForm();
    const getLatLng = vi.fn().mockReturnValue(coordinates);
    const wrapper = factory({
      propsData: {
        form: {
          points: [point, secondPoint],
        },
      },
    });

    const pointMarker = selectPointMarkerByIndex(wrapper, 0);

    pointMarker.vm.$emit('dragend', {
      target: {
        getLatLng,
        options: {
          data: point,
        },
      },
    });
    await flushPromises();

    expect(wrapper).toEmitInput({
      points: [{ ...point, coordinates }, secondPoint],
    });
  });
});
