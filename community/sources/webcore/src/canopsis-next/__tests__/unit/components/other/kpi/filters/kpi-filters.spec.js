import { generateShallowRenderer } from '@unit/utils/vue';

import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';

const mocks = vi.hoisted(() => ({
  fetchFiltersList: vi.fn(),
  createFilter: vi.fn(),
  updateFilter: vi.fn(),
  removeFilter: vi.fn(),
}));

vi.mock('@/hooks/store/modules/filter', () => ({
  useFilter: () => ({
    filters: [],
    filtersMeta: { total_count: 0 },
    filtersPending: false,
    ...mocks,
  }),
}));

const stubs = { 'kpi-filters-list': true };

describe('kpi-filters', () => {
  const factory = generateShallowRenderer(KpiFilters, { stubs });

  it('fetches filters with the initial pagination options', async () => {
    const wrapper = factory();
    await wrapper.vm.$nextTick();

    expect(mocks.fetchFiltersList).toHaveBeenCalled();
    expect(wrapper.vm.options).toMatchObject({ page: 1, itemsPerPage: 10 });
  });

  it('exposes the current filter actions and query updater', () => {
    const wrapper = factory();

    expect(wrapper.vm.fetchList).toEqual(expect.any(Function));
    expect(wrapper.vm.updateOptions).toEqual(expect.any(Function));
    expect(wrapper.vm.showEditFilterModal).toEqual(expect.any(Function));
    expect(wrapper.vm.showDuplicateFilterModal).toEqual(expect.any(Function));
    expect(wrapper.vm.showDeleteFilterModal).toEqual(expect.any(Function));
  });
});
