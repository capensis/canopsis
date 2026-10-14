import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import PbehaviorEntities from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-entities.vue';

const { fetchPbehaviorEntitiesListWithoutStore } = vi.hoisted(() => ({
  fetchPbehaviorEntitiesListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { total_count: 0 },
  }),
}));

vi.mock('@/hooks/store/modules/pbehavior-entities', () => ({
  usePbehaviorEntities: () => ({ fetchPbehaviorEntitiesListWithoutStore }),
}));

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
  'pbehavior-entities-expand-item': true,
};

describe('pbehavior-entities', () => {
  const factory = generateShallowRenderer(PbehaviorEntities, { stubs });

  afterEach(() => {
    fetchPbehaviorEntitiesListWithoutStore.mockClear();
  });

  test('Pbehavior entities fetched after mount', async () => {
    const pbehaviorId = Faker.datatype.string();
    factory({
      propsData: {
        pbehavior: {
          _id: pbehaviorId,
        },
      },
    });

    await flushPromises();

    expect(fetchPbehaviorEntitiesListWithoutStore).toHaveBeenCalledWith({
      id: pbehaviorId,
      params: {
        limit: 10,
        page: 1,
      },
    });
  });
});
