import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'c-enabled-field': true,
};

const selectFilterSelectField = wrapper => wrapper.find('select.v-select');

describe('filter-selector', () => {
  const filters = [
    {
      _id: '3',
      title: 'Filter 1',
      filter: { _id: 1 },
    },
    {
      _id: '4',
      title: 'Filter 2',
      filter: { _id: 2 },
    },
  ];

  const factory = generateShallowRenderer(FilterSelector, { stubs });

  it('Filter changed after trigger the select field', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: filters[0]._id,
      },
    });

    const filterSelectField = selectFilterSelectField(wrapper);

    const [, filter] = filters;

    filterSelectField.triggerCustomEvent('input', filter._id);

    expect(wrapper).toEmitInput(filter._id);
  });
});
