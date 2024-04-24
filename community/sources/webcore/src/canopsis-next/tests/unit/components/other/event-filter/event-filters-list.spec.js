import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import { EVENT_FILTER_TYPES } from '@/constants';

import EventFiltersList from '@/components/other/event-filter/event-filters-list.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search': true,
  'v-checkbox': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-enabled': true,
  'c-table-pagination': true,
  'pbehaviors-create-action-btn': true,
  'pbehaviors-list-action-btn': true,
  'event-filters-list-expand-panel': true,
};

describe('event-filters-list', () => {
  const totalItems = 11;
  const eventFilterTypes = Object.values(EVENT_FILTER_TYPES);

  const eventFilters = range(totalItems).map(index => ({
    _id: `c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf-${index}`,
    type: eventFilterTypes[index % eventFilterTypes.length],
    priority: index,
    enabled: !!(index % 2),
    author: {
      display_name: `author-${index}`,
    },
    created: 1614861888 + index,
    updated: 1614861888 + index,
  }));

  const snapshotFactory = generateRenderer(EventFiltersList, { stubs });

  it('Renders `event-filters-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        options: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `event-filters-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        eventFilters,
        options: {
          page: 2,
          updatedPerPage: 10,
          search: 'Filter',
          sortBy: ['created'],
          sortDesc: [true],
        },
        totalItems: 50,
        pending: true,
        removable: true,
        updatable: true,
        duplicable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `event-filters-list` with expanded panel', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        eventFilters,
        options: {
          page: 1,
          itemsPerPage: 10,
          search: '',
          sortBy: [],
          sortDesc: [],
        },
        totalItems: 50,
      },
    });

    const expandButton = wrapper
      .findAll('tr > td')
      .at(0)
      .find('.v-data-table__expand-icon');

    await expandButton.triggerCustomEvent('expand');

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `event-filters-list` with updatable and old_patterns', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        eventFilters: eventFilters.map(eventFilter => ({ ...eventFilter, old_patterns: true })),
        options: {
          page: 1,
          itemsPerPage: 10,
          search: [],
          sortBy: [],
          descending: false,
        },
        updatable: true,
        totalItems: 50,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
