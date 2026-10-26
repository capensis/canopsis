import { generateShallowRenderer } from '@unit/utils/vue';

import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

import BroadcastMessagesList from '@/components/other/broadcast-message/broadcast-messages-list.vue';

const CAdvancedDataTableStub = {
  name: 'CAdvancedDataTable',
  props: ['headers', 'items'],
  template: '<div />',
};

const stubs = {
  'c-advanced-data-table': CAdvancedDataTableStub,
  'broadcast-message': true,
  'c-layout': true,
  'c-action-btn': true,
};

describe('broadcast-messages-list', () => {
  const factory = generateShallowRenderer(BroadcastMessagesList, { stubs });

  test('uses statuses returned by the backend', () => {
    const broadcastMessages = [
      {
        _id: 'message-id',
        message: 'Message',
        status: BROADCAST_MESSAGES_STATUSES.active,
      },
    ];
    const wrapper = factory({
      props: {
        broadcastMessages,
        options: {},
      },
    });

    expect(BROADCAST_MESSAGES_STATUSES).toEqual({
      active: 1,
      pending: 2,
      expired: 3,
    });
    expect(wrapper.findComponent(CAdvancedDataTableStub).props('items')).toEqual(broadcastMessages);
    expect(wrapper.vm.headers).toEqual(
      expect.arrayContaining([expect.objectContaining({ value: 'status' })]),
    );
  });
});
