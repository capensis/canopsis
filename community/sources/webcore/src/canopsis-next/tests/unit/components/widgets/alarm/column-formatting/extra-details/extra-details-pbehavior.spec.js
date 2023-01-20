import { omit } from 'lodash';
import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import ExtraDetailsPbehavior from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-pbehavior.vue';

const localVue = createVueInstance();

describe('extra-details-pbehavior', () => {
  const prevDateStartTimestamp = 1386382400000;
  const prevDateStopTimestamp = 1386392400000;
  const prevMonthDateStartTimestamp = 1375884800000;
  const prevMonthDateStopTimestamp = 1375894800000;

  const pbehavior = {
    name: 'pbehavior-name',
    author: {
      _id: 'pbehavior-author',
      name: 'pbehavior-author',
    },
    tstart: prevDateStartTimestamp,
    tstop: prevDateStopTimestamp,
    rrule: 'rrule',
    type: {
      name: 'pbehavior-type-name',
      icon_name: 'pbehavior-type-icon',
    },
    reason: {
      name: 'pbehavior-reason-name',
    },
    last_comment: {
      _id: 'pbehavior-comment-1-id',
      author: {
        name: 'pbehavior-comment-1-author',
      },
      message: 'pbehavior-comment-1-message',
    },
  };
  const pbehaviorInfo = {
    type_name: 'type-name',
    icon_name: 'icon-name',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsPbehavior, {
    localVue,
    attachTo: document.body,
  });

  it('Renders `extra-details-pbehavior` with full pbehavior', async () => {
    snapshotFactory({
      propsData: {
        pbehavior,
        pbehaviorInfo,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without reason', async () => {
    snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['reason']),
        pbehaviorInfo,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without tstop', async () => {
    snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['tstop']),
        pbehaviorInfo,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without rrule', async () => {
    snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['rrule']),
        pbehaviorInfo,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without comments', async () => {
    snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['comments']),
        pbehaviorInfo,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` with pbehavior started in previous month', async () => {
    snapshotFactory({
      propsData: {
        pbehavior: {
          ...pbehavior,
          tstart: prevMonthDateStartTimestamp,
          tstop: prevMonthDateStopTimestamp,
        },
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
