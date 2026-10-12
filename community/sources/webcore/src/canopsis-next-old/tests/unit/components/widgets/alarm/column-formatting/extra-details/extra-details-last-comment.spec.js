import { generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsLastComment from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-last-comment.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
  'c-simple-tooltip': true,
};

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

  const snapshotFactory = generateRenderer(ExtraDetailsLastComment, {

    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-last-comment` with full last comment', () => {
    const wrapper = snapshotFactory({
      propsData: {
        lastComment,
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
