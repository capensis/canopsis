import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';
import { selectRowExpandButtonByIndex } from '@unit/utils/table';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityList from '@/components/other/availability/availability-list.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search-field': true,
  'v-checkbox': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
  'availability-list-column-value': true,
  'entity-column-cell': true,
  'availability-list-expand-panel': true,
};

describe('availability-list', () => {
  const totalItems = 11;

  const availabilities = range(totalItems).map(index => ({
    _id: `availability-${index}`,
    uptime_share: 10,
    uptime_duration: 1000,
    downtime_share: 90,
    downtime_duration: 9000,
    entity: {
      _id: `entity-id-${index}`,
      column: `column-value-${index}`,
    },
  }));
  const columns = [{ value: 'column', text: 'Custom column' }];

  const snapshotFactory = generateRenderer(AvailabilityList, { stubs });

  test('Renders `availability-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availabilities,
        options: {},
        interval: {},
        columns,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        availabilities,
        options: {
          page: 2,
          itemsPerPage: 10,
          search: 'Filter',
          sortBy: ['column'],
          sortDesc: [true],
        },
        totalItems,
        pending: true,
        showTrend: true,
        columns,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        showType: AVAILABILITY_SHOW_TYPE.duration,
        interval: {},
      },
    });

    await selectRowExpandButtonByIndex(wrapper, 1).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });
});
