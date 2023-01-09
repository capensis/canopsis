import { range } from 'lodash';
import Faker from 'faker';
import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorModule } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import { selectRowEditButtonByIndex, selectRowRemoveButtonByIndex } from '@unit/utils/table';
import { MODALS } from '@/constants';
import { createEntityIdPatternByValue } from '@/helpers/pattern';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-simple-list.vue';

const localVue = createVueInstance();

const stubs = {
  'c-search-field': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-action-fab-btn': true,
  'c-enabled': true,
  'c-table-pagination': true,
};

const selectAddButton = wrapper => wrapper.find('c-action-fab-btn-stub[icon="add"]');
const selectCalendarButton = wrapper => wrapper.find('c-action-fab-btn-stub[icon="calendar_today"]');

describe('pbehaviors-simple-list', () => {
  const $modals = mockModals();

  const totalItems = 7;
  const pbehaviorsItems = range(totalItems).map(index => ({
    _id: `id-pbehavior-${index}`,
    name: `name-${index}`,
    author: {
      name: `author-${index}`,
    },
    enabled: !!(index % 2),
    editable: !!(index % 2),
    tstart: 1614861000 + index,
    tstop: 1614861200 + index,
    type: {
      name: `type-name-${index}`,
      icon_name: `type-icon-name-${index}`,
    },
    reason: {
      name: `reason-name-${index}`,
    },
    is_active_status: !(index % 2),
  }));

  const { pbehaviorModule, fetchPbehaviorsByEntityIdWithoutStore, removePbehavior } = createPbehaviorModule();

  const store = createMockedStoreModules([pbehaviorModule]);

  const factory = generateShallowRenderer(PbehaviorsSimpleList, {
    localVue,
    stubs,
    mocks: { $modals },
    parentComponent: {
      provide: {
        $system: {
          timezone: process.env.TZ,
        },
      },
    },
  });
  const snapshotFactory = generateRenderer(PbehaviorsSimpleList, {
    localVue,
    stubs,
    mocks: { $modals },
    parentComponent: {
      provide: {
        $system: {
          timezone: process.env.TZ,
        },
      },
    },
  });

  test('Pbehaviors fetched after mount', async () => {
    const entityId = Faker.datatype.string();
    factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
        },
      },
    });

    await flushPromises();

    expect(fetchPbehaviorsByEntityIdWithoutStore).toBeCalledWith(
      expect.any(Object),
      { id: entityId, params: { with_flags: true } },
      undefined,
    );
  });

  test('Pbehavior create modal opened after trigger create button', async () => {
    const entityId = Faker.datatype.string();
    const wrapper = factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
        },
        addable: true,
      },
    });

    await flushPromises();
    fetchPbehaviorsByEntityIdWithoutStore.mockClear();

    selectAddButton(wrapper).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(entityId),
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    await config.afterSubmit();
    expect(fetchPbehaviorsByEntityIdWithoutStore).toBeCalledWith(
      expect.any(Object),
      { id: entityId, params: { with_flags: true } },
      undefined,
    );
  });

  test('Pbehavior calendar modal opened after trigger calendar button', () => {
    const entityId = Faker.datatype.string();
    const entityName = Faker.datatype.string();
    const wrapper = factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
          name: entityName,
        },
        addable: true,
      },
    });

    selectCalendarButton(wrapper).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorsCalendar,
        config: {
          title: `Periodic behaviors - ${entityName}`,
          entityId,
        },
      },
    );
  });

  test('Confirmation remove pbehavior modal opened after trigger remove button', async () => {
    fetchPbehaviorsByEntityIdWithoutStore.mockResolvedValueOnce(pbehaviorsItems);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([pbehaviorModule]),
      propsData: {
        entity: {},
        removable: true,
      },
    });

    await flushPromises();
    fetchPbehaviorsByEntityIdWithoutStore.mockClear();

    const removingIndex = 1;
    const removingPbehavior = pbehaviorsItems[removingIndex];

    selectRowRemoveButtonByIndex(wrapper, removingIndex).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(removePbehavior).toBeCalledWith(
      expect.any(Object),
      { id: removingPbehavior._id },
      undefined,
    );
    expect(fetchPbehaviorsByEntityIdWithoutStore).toBeCalled();
  });

  test('Edit pbehavior modal opened after trigger edit button', async () => {
    fetchPbehaviorsByEntityIdWithoutStore.mockResolvedValueOnce(pbehaviorsItems);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([pbehaviorModule]),
      propsData: {
        entity: {},
        updatable: true,
      },
    });

    await flushPromises();
    fetchPbehaviorsByEntityIdWithoutStore.mockClear();

    const editingIndex = 2;
    const editingPbehavior = pbehaviorsItems[editingIndex];

    selectRowEditButtonByIndex(wrapper, editingIndex).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPlanning,
        config: {
          pbehaviors: [editingPbehavior],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    await config.afterSubmit();

    expect(fetchPbehaviorsByEntityIdWithoutStore).toBeCalled();
  });

  test('Renders `pbehaviors-simple-list` without pbehaviors', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {},
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehaviors-simple-list` with pbehaviors', async () => {
    fetchPbehaviorsByEntityIdWithoutStore.mockResolvedValueOnce(pbehaviorsItems);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([pbehaviorModule]),
      propsData: {
        entity: {},
        withActiveStatus: true,
        updatable: true,
        removable: true,
        addable: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
