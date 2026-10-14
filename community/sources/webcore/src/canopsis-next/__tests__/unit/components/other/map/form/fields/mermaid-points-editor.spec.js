import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MODALS } from '@/constants';

import { mermaidPointToForm } from '@/helpers/entities/map/form';

import MermaidPoints from '@/components/other/map/form/fields/mermaid-points-editor.vue';

const stubs = {
  'mermaid-point-marker': {
    name: 'MermaidPointMarker',
    props: ['x', 'y', 'entity', 'size'],
    emits: ['click', 'contextmenu', 'dblclick', 'mousedown'],
    template: '<mermaid-point-marker-stub />',
  },
  'point-contextmenu': {
    name: 'PointContextmenu',
    props: ['value', 'editing', 'positionX', 'positionY'],
    emits: ['add:point', 'edit:point', 'remove:point', 'close'],
    template: '<point-contextmenu-stub />',
  },
  'point-form-dialog-menu': {
    name: 'PointFormDialogMenu',
    props: ['value', 'point', 'positionX', 'positionY', 'editing', 'existsEntities'],
    emits: ['cancel', 'submit', 'remove'],
    template: '<point-form-dialog-menu-stub />',
  },
};

const selectRoot = wrapper => wrapper.find('div.mermaid-points');
const selectPointContextmenu = wrapper => wrapper.findComponent({ name: 'PointContextmenu' });
const selectPointFormDialogMenu = wrapper => wrapper.findComponent({ name: 'PointFormDialogMenu' });
const selectPointMarkers = wrapper => wrapper.findAllComponents({ name: 'MermaidPointMarker' });
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);

const checkContextMenuPosition = (wrapper, { x, y }) => {
  const menu = selectPointContextmenu(wrapper);

  expect(menu.props('value')).toBe(true);
  expect(menu.props('positionY')).toBe(y);
  expect(menu.props('positionX')).toBe(x);
};

const checkDialogMenuPosition = (wrapper, { x, y }) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  expect(pointFormDialog.props('value')).toBe(true);
  expect(pointFormDialog.props('positionY')).toBe(y);
  expect(pointFormDialog.props('positionX')).toBe(x);
};

const fillPointDialog = (
  wrapper,
  point = mermaidPointToForm({
    x: Faker.datatype.number(),
    y: Faker.datatype.number(),
  }),
) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  pointFormDialog.vm.$emit('submit', point);

  return point;
};

const checkMenuIsClosed = async (wrapper) => {
  await flushPromises();

  expect(wrapper.vm.editingPoint).toBeFalsy();
  expect(wrapper.vm.addingPoint).toBeFalsy();
};

const triggerContextMenuEvent = (wrapper, event) => {
  const contextMenu = selectPointContextmenu(wrapper);

  contextMenu.vm.$emit(event);

  return flushPromises();
};

const triggerPointEvent = (wrapper, index, event) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  pointMarker.vm.$emit(event, new Event(event));

  return flushPromises();
};

const getEvent = () => ({
  clientX: Faker.datatype.number(),
  clientY: Faker.datatype.number(),
  offsetX: Faker.datatype.number(),
  offsetY: Faker.datatype.number(),
});

describe('mermaid-points-editor', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(MermaidPoints, {
    stubs,
    attachTo: document.body,
  });

  test('Contextmenu opened after right click', async () => {
    const event = getEvent();

    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    await selectRoot(wrapper).trigger('contextmenu', event);

    checkContextMenuPosition(wrapper, { x: event.clientX, y: event.clientY });

    const contextMenu = selectPointContextmenu(wrapper);

    expect(contextMenu.props('editing')).toBe(false);
  });

  test('Contextmenu didn\'t change when menu already opened', async () => {
    const event = getEvent();
    const menuPosition = { x: event.clientX, y: event.clientY };

    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    await selectRoot(wrapper).trigger('contextmenu', event);
    checkContextMenuPosition(wrapper, menuPosition);

    await selectRoot(wrapper).trigger('contextmenu', {
      clientX: event.clientX + 1,
      clientY: event.clientY + 1,
    });
    checkContextMenuPosition(wrapper, menuPosition);
  });

  test('Point contextmenu didn\'t change when menu already opened', async () => {
    const point = mermaidPointToForm({
      x: Faker.datatype.number(),
      y: Faker.datatype.number(),
    });
    const secondPoint = mermaidPointToForm({
      x: Faker.datatype.number(),
      y: Faker.datatype.number(),
    });
    const menuPosition = { x: point.x, y: point.y };

    const wrapper = factory({
      propsData: {
        value: [point, secondPoint],
      },
    });

    await triggerPointEvent(wrapper, 0, 'contextmenu');
    checkContextMenuPosition(wrapper, menuPosition);

    await triggerPointEvent(wrapper, 1, 'contextmenu');
    checkContextMenuPosition(wrapper, menuPosition);
  });

  test('Point added after double click', async () => {
    vi.useFakeTimers();
    const event = getEvent();

    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    await selectRoot(wrapper).trigger('dblclick', event);

    checkDialogMenuPosition(wrapper, { x: event.clientX, y: event.clientY });

    const newPoint = fillPointDialog(wrapper);

    expect(wrapper).toEmitInput([newPoint]);

    vi.runAllTimers();

    await checkMenuIsClosed(wrapper);

    vi.useRealTimers();
  });

  test('Point dialog menu didn\'t changed after double click when dialog already opened', async () => {
    const event = getEvent();
    const secondEvent = getEvent();
    const menuPosition = { x: event.clientX, y: event.clientY };
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    await selectRoot(wrapper).trigger('dblclick', event);
    checkDialogMenuPosition(wrapper, menuPosition);

    await selectRoot(wrapper).trigger('dblclick', secondEvent);
    checkDialogMenuPosition(wrapper, menuPosition);
  });

  test('Point didn\'t add after click without prop', async () => {
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    await selectRoot(wrapper).trigger('click');

    const pointFormDialog = selectPointFormDialogMenu(wrapper);

    expect(pointFormDialog.props('value')).toBeFalsy();
  });

  test('Point added after click with prop', async () => {
    vi.useFakeTimers();
    const wrapper = factory({
      propsData: {
        addOnClick: true,
        value: [],
      },
    });

    await selectRoot(wrapper).trigger('click');

    const newPoint = fillPointDialog(wrapper);

    expect(wrapper).toEmitInput([newPoint]);

    vi.runAllTimers();
    await checkMenuIsClosed(wrapper);

    vi.useRealTimers();
  });

  test('Point edited after open form dialog by contextmenu', async () => {
    vi.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        value: [point],
      },
    });

    await triggerPointEvent(wrapper, 0, 'contextmenu');
    await triggerContextMenuEvent(wrapper, 'edit:point');
    const newPoint = fillPointDialog(wrapper, { ...point, entity: 'entity' });

    expect(wrapper).toEmitInput([newPoint]);

    vi.runAllTimers();
    await checkMenuIsClosed(wrapper);

    vi.useRealTimers();
  });

  test('Point edited after open form dialog by point double click', async () => {
    vi.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const secondPoint = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        value: [point, secondPoint],
      },
    });

    await triggerPointEvent(wrapper, 1, 'dblclick');
    const newPoint = fillPointDialog(wrapper, { ...secondPoint, entity: 'entity' });

    expect(wrapper).toEmitInput([point, newPoint]);

    vi.runAllTimers();
    await checkMenuIsClosed(wrapper);

    vi.useRealTimers();
  });

  test('Point removed after trigger remove in contextmenu', async () => {
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        value: [point],
      },
      mocks: {
        $modals,
      },
    });

    await triggerPointEvent(wrapper, 0, 'contextmenu');
    await triggerContextMenuEvent(wrapper, 'remove:point');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(wrapper).toEmitInput([]);
  });

  test('Point moved after trigger mouse events', async () => {
    vi.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        value: [point],
      },
      mocks: {
        $modals,
      },
    });

    await triggerPointEvent(wrapper, 0, 'mousedown');

    vi.runAllTimers();

    wrapper.vm.movePointByEvent({
      offsetX: 100,
      offsetY: 120,
    });

    wrapper.vm.finishMoving();
    await flushPromises();

    expect(wrapper).toEmitInput([{ ...point, x: 100, y: 120 }]);

    vi.runAllTimers();

    expect(wrapper.vm.moving).toBeFalsy();
    expect(wrapper.vm.movingPointIndex).toBeFalsy();

    vi.useRealTimers();
  });

  test('Point didn\'t moving when dialog opened', async () => {
    vi.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });

    const wrapper = factory({
      propsData: {
        value: [point],
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    await triggerPointEvent(wrapper, 0, 'dblclick');
    await triggerPointEvent(wrapper, 0, 'mousedown');

    vi.runAllTimers();

    expect(wrapper.vm.moving).toBe(false);
    expect(wrapper.vm.movingPointIndex).toBeUndefined();

    vi.useRealTimers();
  });
});
