import { flushPromises, generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import CClickableTooltip from '@/components/common/clickable-tooltip/c-clickable-tooltip.vue';
import ExtraDetailsLastComment from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-last-comment.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';

const stubs = {
  'c-clickable-tooltip': CClickableTooltip,
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
  'c-alarm-extra-details-chip': true,
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

  it('Renders `extra-details-last-comment` with full last comment', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        lastComment,
      },
    });

    await flushPromises();

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-last-comment` with date in previous month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        lastComment: {
          ...lastComment,
          t: prevMonthDateTimestamp,
        },
      },
    });

    await flushPromises();

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
