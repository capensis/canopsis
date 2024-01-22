import Faker from 'faker';
import { LatLngBounds } from 'leaflet';
import { omit } from 'lodash';
import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { mockModals } from '@unit/utils/mock-hooks';
import { geomapPointToForm } from '@/helpers/entities/map/form';
import { MODALS } from '@/constants';

import GeomapEditor from '@/components/other/map/form/fields/geomap-editor.vue';

const fitBounds = jest.fn();
const fitWorld = jest.fn();
const panTo = jest.fn();
const closeContextmenu = jest.fn();
const latLngToContainerPoint = jest.fn().mockReturnValue({ x: 0, y: 0 });
const getBounds = jest.fn().mockReturnValue(new LatLngBounds([0, 0], [1, 1]));

const stubs = {
  'v-tooltip': createActivatorElementStub('v-tooltip'),
  geomap: {
    template: '<div v-on="$listeners" class="geomap"><slot/></div>',
    computed: {
      mapObject() {
        return {
          fitBounds,
          fitWorld,
          panTo,
          latLngToContainerPoint,
        };
      },
    },
  },
  'geomap-control-zoom': true,
  'geomap-control-layers': true,
  'geomap-contextmenu': {
    props: ['items', 'markerItems'],
    template: `<div
      class="geomap-contextmenu"
      :items="items"
      :marker-items="markerItems"
    />`,
    methods: {
      close: closeContextmenu,
    },
  },
  'geomap-control': true,
  'geomap-tile-layer': true,
  'geomap-cluster-group': {
    template: '<div class="geomap-cluster-group"><slot/></div>',
    computed: {
      mapObject() {
        return {
          getBounds,
        };
      },
    },
  },
  'geomap-marker': true,
  'c-zoom-overlay': true,
  'geomap-icon': true,
  'point-icon': true,
  'point-form-dialog': true,
  'c-help-icon': true,
};
const snapshotStubs = omit(stubs, ['v-tooltip']);

const selectGeomap = wrapper => wrapper.find('div.geomap');
const selectAddLocationBtn = wrapper => wrapper.find('v-btn-stub');
const selectGeomapContextmenu = wrapper => wrapper.find('.geomap-contextmenu');
const selectPointFormDialogMenu = wrapper => wrapper.find('point-form-dialog-menu-stub');
const selectPointMarkers = wrapper => wrapper.findAll('geomap-marker-stub');
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);

const getEvent = event => ({
  latlng: [Faker.datatype.number(), Faker.datatype.number()],
  ...event,
});

const triggerPointEvent = (wrapper, index, event) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  return pointMarker.triggerCustomEvent(event, new Event(event));
};

const checkDialogMenuPosition = (wrapper, { x, y }) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  expect(pointFormDialog.vm.value).toBe(true);
  expect(pointFormDialog.vm.positionY).toBe(y);
  expect(pointFormDialog.vm.positionX).toBe(x);
};

const fillPointDialog = (
  wrapper,
  point = geomapPointToForm(),
) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  pointFormDialog.triggerCustomEvent('submit', point);

  return point;
};

const checkMenuIsClosed = (wrapper) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  expect(pointFormDialog.vm.value).toBeFalsy();
};

const triggerContextMenuEvent = (wrapper, event, data) => {
  const geomapContextmenu = selectGeomapContextmenu(wrapper);

  const { items, markerItems } = geomapContextmenu.vm;

  const item = {
    add: items[0],
    edit: markerItems[0],
    remove: markerItems[1],
  }[event];

  item.action(data);
};

const triggerItemContextMenuEvent = (wrapper, index, event) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  const marker = {
    options: pointMarker.vm.options,
  };

  triggerContextMenuEvent(wrapper, event, { marker });
};

describe('geomap-editor', () => {
  const $modals = mockModals();
  const initialForm = {
    points: [],
  };

  const factory = generateShallowRenderer(GeomapEditor, { stubs });
  const snapshotFactory = generateRenderer(GeomapEditor, {
    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Add on click mode enabled after trigger add location btn', async () => {
    const event = getEvent();
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const geomap = selectGeomap(wrapper);
    const getBoundingClientRect = jest.spyOn(geomap.element, 'getBoundingClientRect')
      .mockReturnValue({ x: 12, y: 22 });
    latLngToContainerPoint.mockReturnValueOnce({
      x: 88,
      y: 28,
    });

    const addLocationBtn = selectAddLocationBtn(wrapper);

    await addLocationBtn.triggerCustomEvent('click');

    await geomap.trigger('click', event);

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

    await selectGeomap(wrapper).trigger('click');

    const pointFormDialog = selectPointFormDialogMenu(wrapper);

    expect(pointFormDialog.vm.value).toBeFalsy();
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

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(validator.errors.items).toHaveLength(1);

    wrapper.setProps({
      form: {
        points: [geomapPointToForm()],
      },
    });

    await flushPromises();

    expect(validator.errors.items).toHaveLength(0);
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

    await selectGeomap(wrapper).trigger('dblclick', event);

    checkDialogMenuPosition(wrapper, menuPosition);

    const newPoint = fillPointDialog(wrapper);

    expect(wrapper).toEmit('input', {
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

    await selectGeomap(wrapper).trigger('dblclick', event);
    checkDialogMenuPosition(wrapper, firstMenuPosition);

    await selectGeomap(wrapper).trigger('dblclick', secondEvent);
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

    expect(wrapper).toEmit('input', { points: [newPoint] });

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
    const latlng = [Faker.datatype.number(), Faker.datatype.number()];

    await triggerContextMenuEvent(wrapper, 'add', { latlng });
    const newPoint = fillPointDialog(wrapper, { ...point, entity: 'entity' });

    expect(wrapper).toEmit('input', { points: [point, newPoint] });

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

    pointFormDialog.triggerCustomEvent('remove');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(wrapper).toEmit('input', { points: [] });

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

    expect(wrapper).toEmit('input', {
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

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(wrapper).toEmit('input', { points: [point] });
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

    await pointFormDialog.triggerCustomEvent('fly:coordinates', coordinates);

    expect(pointFormDialog.vm.point.coordinates).toBe(coordinates);
  });

  test('Marker moved after dragend', async () => {
    const coordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };
    const point = geomapPointToForm();
    const secondPoint = geomapPointToForm();
    const getLatLng = jest.fn().mockReturnValue(coordinates);
    const wrapper = factory({
      propsData: {
        form: {
          points: [point, secondPoint],
        },
      },
    });

    const pointMarker = selectPointMarkerByIndex(wrapper, 0);

    pointMarker.triggerCustomEvent('dragend', {
      target: {
        getLatLng,
        options: {
          data: point,
        },
      },
    });

    expect(wrapper).toEmit('input', {
      points: [{ ...point, coordinates }, secondPoint],
    });
  });

  test('Renders `geomap-editor` with form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `geomap-editor` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          points: [
            geomapPointToForm({ coordinates: { lat: 1, lng: 2 } }),
          ],
        },
        minZoom: 3,
        iconSize: 12,
        name: 'custom_name',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `geomap-editor` with validation errors ', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
        name: 'validation_name',
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});
