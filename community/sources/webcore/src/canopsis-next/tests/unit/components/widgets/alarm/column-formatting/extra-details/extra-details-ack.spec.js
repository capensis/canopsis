import { omit } from 'lodash';
import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsAck from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ack.vue';

describe('extra-details-ack', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const snapshotFactory = generateRenderer(ExtraDetailsAck, {

    attachTo: document.body,
  });

  const ack = {
    a: 'ack-author',
    t: prevDateTimestamp,
    initiator: 'ack-initiator',
    m: 'ack-message',
  };

  it('Renders `extra-details-ack` with full ack', async () => {
    snapshotFactory({
      propsData: {
        ack,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ack` without initiator', async () => {
    snapshotFactory({
      propsData: {
        ack: omit(ack, ['initiator']),
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ack` without message', async () => {
    snapshotFactory({
      propsData: {
        ack: omit(ack, ['m']),
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ack` with date in previous month', async () => {
    snapshotFactory({
      propsData: {
        ack: {
          ...ack,
          t: prevMonthDateTimestamp,
        },
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
