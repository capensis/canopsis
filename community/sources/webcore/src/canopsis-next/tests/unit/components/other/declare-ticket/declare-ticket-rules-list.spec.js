import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import DeclareTicketRulesList from '@/components/other/declare-ticket/declare-ticket-rules-list.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search-field': true,
  'v-checkbox': true,
  'v-checkbox-functional': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-enabled': true,
  'c-table-pagination': true,
};

describe('declare-ticket-rules-list', () => {
  const totalItems = 11;

  const declareTicketRules = range(totalItems).map(index => ({
    _id: `c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf-${index}`,
    name: `name-${index}`,
    system_name: `system-name-${index}`,
    enabled: !!(index % 2),
    created: 1614861888 + index,
    updated: 1614861888 + index,
    author: {
      name: `author-name-${index}`,
    },
  }));

  const snapshotFactory = generateRenderer(DeclareTicketRulesList, { stubs });

  it('Renders `declare-ticket-rules-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pagination: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `declare-ticket-rules-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        declareTicketRules,
        pagination: {
          page: 2,
          rowsPerPage: 10,
          search: 'Rule',
          sortBy: 'created',
          descending: true,
        },
        totalItems: 50,
        pending: true,
        removable: true,
        updatable: true,
        duplicable: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
