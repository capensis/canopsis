import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import ExtraDetailsPbehavior from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-pbehavior.vue';

describe('extra-details-pbehavior', () => {
  const prevDateStartTimestamp = 1386382400000;
  const prevDateStopTimestamp = 1386392400000;
  const prevMonthDateStartTimestamp = 1375884800000;
  const prevMonthDateStopTimestamp = 1375894800000;

  const pbehavior = {
    name: 'pbehavior-name',
    author: {
      _id: 'pbehavior-author',
      display_name: 'pbehavior-author',
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
        display_name: 'pbehavior-comment-1-author',
      },
      message: 'pbehavior-comment-1-message',
    },
  };
  const pbehaviorInfo = {
    type_name: 'type-name',
    icon_name: 'icon-name',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsPbehavior);

  it('Renders `extra-details-pbehavior` with full pbehavior', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior,
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-pbehavior` without reason', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['reason']),
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-pbehavior` without tstop', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['tstop']),
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-pbehavior` without rrule', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['rrule']),
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-pbehavior` without comment', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['last_comment']),
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-pbehavior` without comment author', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          ...pbehavior,
          last_comment: {
            author: null,
            ...pbehavior.last_comment,
          },
        },
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-pbehavior` with pbehavior started in previous month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          ...pbehavior,
          tstart: prevMonthDateStartTimestamp,
          tstop: prevMonthDateStopTimestamp,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
