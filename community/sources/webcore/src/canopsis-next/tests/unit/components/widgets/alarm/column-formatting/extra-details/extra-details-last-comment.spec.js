import { mount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import CClickableTooltip from '@/components/common/clickable-tooltip/c-clickable-tooltip.vue';
import ExtraDetailsLastComment from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-last-comment.vue';

const localVue = createVueInstance();

const stubs = {
  'c-clickable-tooltip': CClickableTooltip,
};

const snapshotFactory = (options = {}) => mount(ExtraDetailsLastComment, {
  localVue,
  stubs,

  ...options,
});

describe('extra-details-last-comment', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const lastComment = {
    a: 'lastComment-author',
    t: prevDateTimestamp,
    m: 'lastComment-message',
  };

  it('Renders `extra-details-last-comment` with full last comment', () => {
    const wrapper = snapshotFactory({
      propsData: {
        lastComment,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-last-comment` with date in previous month', () => {
    const wrapper = snapshotFactory({
      propsData: {
        lastComment: {
          ...lastComment,
          t: prevMonthDateTimestamp,
        },
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
