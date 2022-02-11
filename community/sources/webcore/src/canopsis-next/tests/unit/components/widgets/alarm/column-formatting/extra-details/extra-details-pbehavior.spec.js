import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';

import ExtraDetailsPbehavior from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-pbehavior.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsPbehavior, {
  localVue,

  ...options,
});

describe('extra-details-pbehavior', () => {
  const prevDateStartTimestamp = 1386382400000;
  const prevDateStopTimestamp = 1386392400000;
  const prevMonthDateStartTimestamp = 1375884800000;
  const prevMonthDateStopTimestamp = 1375894800000;

  const pbehavior = {
    name: 'pbehavior-name',
    author: 'pbehavior-author',
    tstart: prevDateStartTimestamp,
    tstop: prevDateStopTimestamp,
    rrule: 'rrule',
    type_name: 'type-name',
    icon_name: 'icon-name',
    type: {
      name: 'pbehavior-type-name',
      icon_name: 'pbehavior-type-icon',
    },
    reason: 'pbehavior-reason-name',
  };
  const comments = [
    {
      _id: 'pbehavior-comment-1-id',
      author: 'pbehavior-comment-1-author',
      message: 'pbehavior-comment-1-message',
    },
    {
      _id: 'pbehavior-comment-2-id',
      author: 'pbehavior-comment-2-author',
      message: 'pbehavior-comment-2-message',
    },
  ];

  it('Renders `extra-details-pbehavior` with full pbehavior', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without reason', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['reason']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without tstop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['tstop']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without rrule', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: omit(pbehavior, ['rrule']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` with comments', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior,
        comments,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` with pbehavior started in previous month', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          ...pbehavior,
          tstart: prevMonthDateStartTimestamp,
          tstop: prevMonthDateStopTimestamp,
        },
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
