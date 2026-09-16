import { createMockedStoreModules, createPbehaviorEntitiesModule } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import PbehaviorEntities from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-entities.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
  'pbehavior-entities-expand-item': true,
};

describe('pbehavior-entities', () => {
  const { pbehaviorEntitiesModule, fetchPbehaviorEntitiesListWithoutStore } =
    createPbehaviorEntitiesModule();

  const store = createMockedStoreModules([pbehaviorEntitiesModule]);

  const factory = generateShallowRenderer(PbehaviorEntities, { stubs });

  test('Pbehavior entities fetched after mount', async () => {
    const pbehaviorId = Faker.datatype.string();
    factory({
      store,
      propsData: {
        pbehavior: {
          _id: pbehaviorId,
        },
      },
    });

    await flushPromises();

    expect(fetchPbehaviorEntitiesListWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: pbehaviorId,
      params: {
        limit: 10,
        page: 1,
      },
    });
  });
});
