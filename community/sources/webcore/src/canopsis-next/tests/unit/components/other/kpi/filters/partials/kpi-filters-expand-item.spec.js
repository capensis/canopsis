import { mount, createVueInstance } from '@unit/utils/vue';

import KpiFiltersExpandItem from '@/components/other/kpi/filters/partials/kpi-filters-expand-item.vue';

const localVue = createVueInstance();

const stubs = {
  'c-patterns-field': true,
};

const snapshotFactory = (options = {}) => mount(KpiFiltersExpandItem, {
  localVue,
  stubs,

  ...options,
});

describe('kpi-filters-expand-item', () => {
  const filter = {
    _id: 'id',
    entity_patterns: [{
      connector: 'connector',
    }],
  };

  it('Renders `kpi-filters-expand-item` with filter', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filter,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
