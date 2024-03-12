import { generateRenderer } from '@unit/utils/vue';

import KpiFiltersExpandItem from '@/components/other/kpi/filters/partials/kpi-filters-expand-item.vue';

const stubs = {
  'c-patterns-field': true,
};

describe('kpi-filters-expand-item', () => {
  const filter = {
    _id: 'id',
    entity_patterns: [{
      connector: 'connector',
    }],
  };

  const snapshotFactory = generateRenderer(KpiFiltersExpandItem, { stubs });

  it('Renders `kpi-filters-expand-item` with filter', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filter,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
