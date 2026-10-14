import { generateShallowRenderer } from '@unit/utils/vue';

import Filters from '@/components/sidebars/form/fields/filters.vue';

const stubs = {
  'widget-settings-item': true,
  'filter-selector': true,
  'field-filters-list': true,
};

const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');
const selectFiltersList = wrapper => wrapper.find('field-filters-list-stub');

describe('filters', () => {
  const filters = [
    {
      _id: 'filter-id-1',
      name: 'Filter 1',
    },
    {
      _id: 'filter-id-2',
      name: 'Filter 2',
    },
  ];

  const factory = generateShallowRenderer(Filters, { stubs });

  it('Selected filters updated after trigger input on the filter selector field', () => {
    const wrapper = factory();

    selectFilterSelectorField(wrapper).triggerCustomEvent('input', filters[0]);

    expect(wrapper).toEmitInput(filters[0]);
  });

  it('Filters updated after trigger input event on filters list', () => {
    const wrapper = factory();

    const newFilters = [...filters].reverse();

    selectFiltersList(wrapper).triggerCustomEvent('update:filters', newFilters);

    expect(wrapper).toEmit('update:filters', newFilters);
  });
});
