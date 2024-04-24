import { range } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorEntitiesModule } from '@unit/utils/store';

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
  const totalItems = 11;

  const pbehavior = {
    _id: 'pbehavior',
  };

  const pbehaviorEntities = range(totalItems).map(index => ({
    _id: `id-pbehavior-entity-${index}`,
    name: `name-pbehavior-entity-${index}`,
    type: `type-pbehavior-entity-${index}`,
  }));

  const { pbehaviorEntitiesModule, fetchPbehaviorEntitiesListWithoutStore } = createPbehaviorEntitiesModule();

  const store = createMockedStoreModules([pbehaviorEntitiesModule]);

  const factory = generateShallowRenderer(PbehaviorEntities, { stubs });
  const snapshotFactory = generateRenderer(PbehaviorEntities, { stubs });

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

    expect(fetchPbehaviorEntitiesListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: pbehaviorId,
        params: {
          limit: 10,
          page: 1,
        },
      },
      undefined,
    );
  });

  test('Renders `pbehavior-entities` without entities', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehavior-entities` with entities', async () => {
    fetchPbehaviorEntitiesListWithoutStore.mockResolvedValueOnce({
      data: pbehaviorEntities,
      meta: { total_count: totalItems },
    });
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([pbehaviorEntitiesModule]),
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
