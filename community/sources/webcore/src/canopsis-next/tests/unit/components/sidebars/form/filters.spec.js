import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';

import Filters from '@/components/sidebars/form/fields/filters.vue';

const stubs = {
  'widget-settings-item': true,
  'filter-selector': true,
  'field-filters-list': true,
};

const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');
const selectFiltersList = wrapper => wrapper.find('field-filters-list-stub');

describe('filters', () => {
  const $modals = mockModals();
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

  const factory = generateShallowRenderer(Filters, {

    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(Filters, { stubs });

  it('Selected filters updated after trigger input on the filter selector field', () => {
    const wrapper = factory();

    selectFilterSelectorField(wrapper).triggerCustomEvent('input', filters[0]);

    expect(wrapper).toEmit('input', filters[0]);
  });

  it('Filters updated after trigger input event on filters list', () => {
    const wrapper = factory();

    const newFilters = [...filters].reverse();

    selectFiltersList(wrapper).triggerCustomEvent('input', newFilters);

    expect(wrapper).toEmit('update:filters', newFilters);
  });

  it('Renders `filters` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `filters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        widgetId: Faker.datatype.string(),
        value: filters[0]._id,
        addable: false,
        editable: false,
        withAlarm: false,
        withEntity: false,
        withPbehavior: false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
