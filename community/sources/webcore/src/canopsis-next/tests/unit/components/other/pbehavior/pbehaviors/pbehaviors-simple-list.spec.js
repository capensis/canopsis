import { range } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorModule } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';

import { MODALS } from '@/constants';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-action-fab-btn': true,
  'c-enabled': true,
  'c-table-pagination': true,
  'pbehavior-actions': true,
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
      display_name: `author-${index}`,
    },
    enabled: !!(index % 2),
    editable: !!(index % 2),
    tstart: 1614861000 + index,
    tstop: 1614861200 + index,
    rrule: index % 2 ? 'RRULE:' : null,
    rrule_end: index % 4 ? 1614861888 + index : null,
    type: {
      name: `type-name-${index}`,
      icon_name: `type-icon-name-${index}`,
    },
    reason: {
      name: `reason-name-${index}`,
    },
    is_active_status: !(index % 2),
  }));

  const { pbehaviorModule, fetchPbehaviorsByEntityIdWithoutStore } = createPbehaviorModule();

  const store = createMockedStoreModules([pbehaviorModule]);

  const factory = generateShallowRenderer(PbehaviorsSimpleList, {
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

    selectAddButton(wrapper).triggerCustomEvent('click');

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

    selectCalendarButton(wrapper).triggerCustomEvent('click');

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

  test('Renders `pbehaviors-simple-list` without pbehaviors', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {},
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
