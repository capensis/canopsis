import { omit } from 'lodash';
import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsCanceled from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-canceled.vue';

describe('extra-details-canceled', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const canceled = {
    a: 'cancelled-author',
    t: prevDateTimestamp,
    m: 'cancelled-message',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsCanceled, {

    attachTo: document.body,
  });

  it('Renders `extra-details-canceled` with full canceled', async () => {
    snapshotFactory({
      propsData: {
        canceled,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-canceled` without message', async () => {
    snapshotFactory({
      propsData: {
        canceled: omit(canceled, ['m']),
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-canceled` with date in previous month', async () => {
    snapshotFactory({
      propsData: {
        canceled: {
          ...canceled,
          t: prevMonthDateTimestamp,
        },
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
