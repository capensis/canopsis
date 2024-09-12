import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import PbehaviorsList from '@/components/other/pbehavior/pbehaviors/pbehaviors-list.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-advanced-search-field': true,
  'c-expand-btn': true,
  'pbehavior-actions': true,
  'c-table-pagination': true,
  'c-ellipsis': true,
  'c-enabled': true,
  'pbehaviors-mass-actions-panel': true,
  'pbehaviors-list-expand-item': true,
};

const selectExpandButtonByRow = (wrapper, index) => wrapper
  .findAll('tbody > tr')
  .at(index)
  .find('c-expand-btn-stub');

describe('pbehaviors-list', () => {
  const totalItems = 11;

  const pbehaviorsItems = range(totalItems).map(index => ({
    _id: `id-pbehavior-${index}`,
    name: `name-${index}`,
    enabled: !!(index % 2),
    tstart: 1614861000 + index,
    tstop: 1614861200 + index,
    last_alarm_date: 1614861250 + index,
    created: 1614861888 + index,
    updated: 1614861888 + index,
    rrule: index % 2 ? 'RRULWE' : null,
    rrule_end: index % 4 ? 1614861888 + index : null,
    is_active_status: !(index % 2),
    type: {
      icon_name: `type-icon-name-${index}`,
    },
  }));

  const snapshotFactory = generateRenderer(PbehaviorsList, {
    stubs,
    parentComponent: {
      provide: {
        $system: {
          timezone: process.env.TZ,
        },
      },
    },
  });

  test('Renders `pbehaviors-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviors: [],
        pagination: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviors: pbehaviorsItems,
        pagination: {
          page: 2,
          rowsPerPage: 10,
          search: 'Filter',
          sortBy: 'created',
          descending: true,
        },
        totalItems: 50,
        pending: true,
        removable: true,
        updatable: true,
        duplicable: true,
        enablable: true,
        disablable: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list` with expanded panel', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviors: pbehaviorsItems,
        pagination: {
          page: 1,
          rowsPerPage: 10,
          search: '',
          sortBy: '',
          descending: false,
        },
        totalItems: 50,
      },
    });

    await selectExpandButtonByRow(wrapper, 0).vm.$emit('expand');

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list` with updatable and old_mongo_query', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviors: pbehaviorsItems.map(item => ({ ...item, old_mongo_query: true })),
        pagination: {
          page: 1,
          rowsPerPage: 10,
          search: '',
          sortBy: '',
          descending: false,
        },
        updatable: true,
        totalItems: 50,
      },
    });

    await selectExpandButtonByRow(wrapper, 0).vm.$emit('expand');

    expect(wrapper.element).toMatchSnapshot();
  });
});
