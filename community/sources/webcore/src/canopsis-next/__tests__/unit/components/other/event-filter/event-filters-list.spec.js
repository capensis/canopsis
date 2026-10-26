import { generateShallowRenderer } from '@unit/utils/vue';

import EventFiltersList from '@/components/other/event-filter/event-filters-list.vue';

const CAdvancedDataTableStub = {
  name: 'CAdvancedDataTable',
  props: ['headers'],
  template: '<div />',
};

const stubs = {
  'c-advanced-data-table': CAdvancedDataTableStub,
  'c-table-mass-actions-panel': true,
  'c-enabled': true,
  'c-circle-badge': true,
  'c-action-btn': true,
  'c-db-export-btn': true,
};

describe('event-filters-list', () => {
  const factory = generateShallowRenderer(EventFiltersList, { stubs });

  it('shows the name column after the ID column', () => {
    const wrapper = factory({
      props: {
        options: {},
      },
    });

    const { headers } = wrapper.findComponent(CAdvancedDataTableStub).props();

    expect(headers.slice(0, 3)).toEqual([
      { text: 'Id', value: '_id' },
      { text: 'Name', value: 'name' },
      { text: 'Type', value: 'type' },
    ]);
  });
});
