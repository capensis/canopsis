import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import { EVENT_ENTITY_TYPES } from '@/constants';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import DeclaredTicketsList from '@/components/other/declare-ticket/declared-tickets-list.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search-field': true,
  'v-checkbox': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-enabled': true,
  'c-table-pagination': true,
  'c-help-icon': true,
};

describe('declared-tickets-list', () => {
  const totalItems = 11;

  const tickets = range(totalItems).map(index => ({
    _id: `c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf-${index}`,
    ticket_url: index % 2 === 0 ? `https://ticket-url.com/${index}` : `ticket-url-${index}`,
    ticket: `ticket-${index}`,
    ticket_system_name: `ticket-system-name-${index}`,
    ticket_rule_name: `ticket-rule-name-${index}`,
    t: 1614861888 + index,
    _t: index % 2 === 0
      ? EVENT_ENTITY_TYPES.declareTicket
      : EVENT_ENTITY_TYPES.declareTicketFail,
    ticket_meta_alarm_id: index % 2 === 0
      ? `ticket-meta-alarm-id-${index}`
      : undefined,
    a: `author-${index}`,
    ticket_comment: `ticket-comment-${index}`,
  }));

  const snapshotFactory = generateRenderer(DeclaredTicketsList, { stubs });

  it('Renders `declared-tickets-list` with tickets', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tickets,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `declared-tickets-list` without tickets', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tickets: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `declared-tickets-list` with parant alarm id', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tickets: [],
        parentAlarmId: 'ticket-meta-alarm-id-2',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
