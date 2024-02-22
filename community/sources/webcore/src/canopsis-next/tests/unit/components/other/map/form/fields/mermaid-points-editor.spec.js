import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';

import { MODALS } from '@/constants';

import { mermaidPointToForm } from '@/helpers/entities/map/form';

import MermaidPoints from '@/components/other/map/form/fields/mermaid-points-editor.vue';

const stubs = {
  'mermaid-point-marker': true,
  'point-contextmenu': true,
  'point-form-dialog-menu': true,
};

const selectRoot = wrapper => wrapper.find('div.mermaid-points');
const selectPointContextmenu = wrapper => wrapper.find('point-contextmenu-stub');
const selectPointFormDialogMenu = wrapper => wrapper.find('point-form-dialog-menu-stub');
const selectPointMarkers = wrapper => wrapper.findAll('mermaid-point-marker-stub');
const selectPointMarkerByIndex = (wrapper, index) => selectPointMarkers(wrapper).at(index);

const checkContextMenuPosition = (wrapper, { x, y }) => {
  const menu = selectPointContextmenu(wrapper);

  expect(menu.vm.value).toBe(true);
  expect(menu.vm.positionY).toBe(y);
  expect(menu.vm.positionX).toBe(x);
};

const checkDialogMenuPosition = (wrapper, { x, y }) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  expect(pointFormDialog.vm.value).toBe(true);
  expect(pointFormDialog.vm.positionY).toBe(y);
  expect(pointFormDialog.vm.positionX).toBe(x);
};

const fillPointDialog = (
  wrapper,
  point = mermaidPointToForm({
    x: Faker.datatype.number(),
    y: Faker.datatype.number(),
  }),
) => {
  const pointFormDialog = selectPointFormDialogMenu(wrapper);

  pointFormDialog.triggerCustomEvent('submit', point);

  return point;
};

const checkMenuIsClosed = async (wrapper) => {
  await flushPromises();

  expect(wrapper.vm.editingPoint).toBeFalsy();
  expect(wrapper.vm.addingPoint).toBeFalsy();
};

const triggerContextMenuEvent = (wrapper, event) => {
  const contextMenu = selectPointContextmenu(wrapper);

  return contextMenu.triggerCustomEvent(event);
};

const triggerPointEvent = (wrapper, index, event) => {
  const pointMarker = selectPointMarkerByIndex(wrapper, index);

  return pointMarker.triggerCustomEvent(event, new Event(event));
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
  const snapshotFactory = generateRenderer(MermaidPoints, {
    stubs,
    attachTo: document.body,
  });

  test('Contextmenu opened after right click', async () => {
    const event = getEvent();

    const wrapper = factory({
      propsData: {
        points: [],
      },
    });

    await selectRoot(wrapper).trigger('contextmenu', event);

    checkContextMenuPosition(wrapper, { x: event.clientX, y: event.clientY });

    const contextMenu = selectPointContextmenu(wrapper);

    expect(contextMenu.vm.editing).toBe(false);
  });

  test('Contextmenu didn\'t change when menu already opened', async () => {
    const event = getEvent();
    const menuPosition = { x: event.clientX, y: event.clientY };

    const wrapper = factory({
      propsData: {
        points: [],
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
        points: [point, secondPoint],
      },
    });

    await triggerPointEvent(wrapper, 0, 'contextmenu');
    checkContextMenuPosition(wrapper, menuPosition);

    await triggerPointEvent(wrapper, 1, 'contextmenu');
    checkContextMenuPosition(wrapper, menuPosition);
  });

  test('Point added after double click', async () => {
    jest.useFakeTimers();
    const event = getEvent();

    const wrapper = factory({
      propsData: {
        points: [],
      },
    });

    await selectRoot(wrapper).trigger('dblclick', event);

    checkDialogMenuPosition(wrapper, { x: event.clientX, y: event.clientY });

    const newPoint = fillPointDialog(wrapper);

    expect(wrapper).toEmitInput([newPoint]);

    jest.runAllTimers();

    await checkMenuIsClosed(wrapper);

    jest.useRealTimers();
  });

  test('Point dialog menu didn\'t changed after double click when dialog already opened', async () => {
    const event = getEvent();
    const secondEvent = getEvent();
    const menuPosition = { x: event.clientX, y: event.clientY };
    const wrapper = factory({
      propsData: {
        points: [],
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
        points: [],
      },
    });

    await selectRoot(wrapper).trigger('click');

    const pointFormDialog = selectPointFormDialogMenu(wrapper);

    expect(pointFormDialog.vm.value).toBeFalsy();
  });

  test('Point added after click with prop', async () => {
    jest.useFakeTimers();
    const wrapper = factory({
      propsData: {
        addOnClick: true,
        points: [],
      },
    });

    await selectRoot(wrapper).trigger('click');

    const newPoint = fillPointDialog(wrapper);

    expect(wrapper).toEmitInput([newPoint]);

    jest.runAllTimers();
    await checkMenuIsClosed(wrapper);

    jest.useRealTimers();
  });

  test('Point edited after open form dialog by contextmenu', async () => {
    jest.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        points: [point],
      },
    });

    await triggerPointEvent(wrapper, 0, 'contextmenu');
    await triggerContextMenuEvent(wrapper, 'edit:point');
    const newPoint = fillPointDialog(wrapper, { ...point, entity: 'entity' });

    expect(wrapper).toEmitInput([newPoint]);

    jest.runAllTimers();
    await checkMenuIsClosed(wrapper);

    jest.useRealTimers();
  });

  test('Point edited after open form dialog by point double click', async () => {
    jest.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const secondPoint = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        points: [point, secondPoint],
      },
    });

    await triggerPointEvent(wrapper, 1, 'dblclick');
    const newPoint = fillPointDialog(wrapper, { ...secondPoint, entity: 'entity' });

    expect(wrapper).toEmitInput([point, newPoint]);

    jest.runAllTimers();
    await checkMenuIsClosed(wrapper);

    jest.useRealTimers();
  });

  test('Point removed after trigger remove in contextmenu', async () => {
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        points: [point],
      },
      mocks: {
        $modals,
      },
    });

    await triggerPointEvent(wrapper, 0, 'contextmenu');
    await triggerContextMenuEvent(wrapper, 'remove:point');

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

    expect(wrapper).toEmitInput([]);
  });

  test('Point moved after trigger mouse events', async () => {
    jest.useFakeTimers();
    const point = mermaidPointToForm({ x: 1, y: 1 });
    const wrapper = factory({
      propsData: {
        points: [point],
      },
      mocks: {
        $modals,
      },
    });

    await triggerPointEvent(wrapper, 0, 'mousedown');

    jest.runAllTimers();

    await selectRoot(wrapper).trigger('mousemove', {
      offsetX: 100,
      offsetY: 120,
    });

    await selectRoot(wrapper).trigger('mouseup');

    expect(wrapper).toEmitInput([{ ...point, x: 100, y: 120 }]);

    jest.runAllTimers();

    expect(wrapper.vm.moving).toBeFalsy();
    expect(wrapper.vm.movingPointIndex).toBeFalsy();

    jest.useRealTimers();
  });

  test('Point didn\'t moving when dialog opened', async () => {
    const addEventListener = jest.spyOn(window, 'addEventListener');
    const point = mermaidPointToForm({ x: 1, y: 1 });

    const wrapper = factory({
      propsData: {
        points: [point],
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    expect(addEventListener).toHaveBeenNthCalledWith(1, 'resize', expect.any(Function), { passive: true });
    expect(addEventListener).toHaveBeenNthCalledWith(2, 'resize', expect.any(Function), { passive: true });
    addEventListener.mockClear();

    await triggerPointEvent(wrapper, 0, 'dblclick');
    await triggerPointEvent(wrapper, 0, 'mousedown');

    expect(addEventListener).not.toBeCalled();

    addEventListener.mockClear();
  });

  test('Renders `mermaid-points-editor` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        points: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mermaid-points-editor` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        points: [
          mermaidPointToForm({
            x: 100,
            y: 100,
          }),
          mermaidPointToForm({
            x: 150,
            y: 150,
          }),
        ],
        addOnClick: true,
        markerSize: 16,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
