import { mockModals, mockRouter } from '@unit/utils/mock-hooks';
import {
  createInfoModule,
  createMockedStoreModules,
  createPbehaviorCommentModule,
  createPbehaviorModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MODALS } from '@/constants';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-action-fab-btn': true,
  'c-enabled': true,
  'c-table-pagination': true,
  'c-timezone-field': true,
  'pbehavior-actions': true,
};

const selectAddButton = wrapper => wrapper.find('c-action-fab-btn-stub[icon="add"]');
const selectCalendarButton = wrapper =>
  wrapper.find('c-action-fab-btn-stub[icon="calendar_today"]');

describe('pbehaviors-simple-list', () => {
  const $modals = mockModals();
  const $router = mockRouter();

  const { infoModule } = createInfoModule();
  const { pbehaviorModule, fetchPbehaviorsByEntityIdWithoutStore } = createPbehaviorModule();
  const { pbehaviorCommentModule } = createPbehaviorCommentModule();

  const store = createMockedStoreModules([infoModule, pbehaviorModule, pbehaviorCommentModule]);

  const factory = generateShallowRenderer(PbehaviorsSimpleList, {
    stubs,
    mocks: { $modals, $router },
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

    expect(fetchPbehaviorsByEntityIdWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: entityId,
      params: { with_flags: true },
    });
  });

  test('Pbehavior create modal opened after trigger create button', async () => {
    const entity = {
      _id: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        entity,
        addable: true,
      },
    });

    await flushPromises();
    fetchPbehaviorsByEntityIdWithoutStore.mockClear();

    selectAddButton(wrapper).triggerCustomEvent('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.pbehaviorPlanning,
      config: {
        entityPattern: createEntityIdPatternByValue(entity._id),
        entities: [entity],
        withInherited: false,
        defaultInherited: false,
        afterSubmit: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    await config.afterSubmit();
    expect(fetchPbehaviorsByEntityIdWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: entity._id,
      params: { with_flags: true },
    });
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
        defaultInherited: true,
        withInherited: true,
        addable: true,
      },
    });

    selectCalendarButton(wrapper).triggerCustomEvent('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.pbehaviorsCalendar,
      config: {
        title: `Periodic behaviors - ${entityName}`,
        entityId,
        withInherited: true,
        defaultInherited: true,
      },
    });
  });
});
