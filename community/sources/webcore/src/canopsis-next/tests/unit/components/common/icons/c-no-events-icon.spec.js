import { generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import CNoEventsIcon from '@/components/common/icons/c-no-events-icon.vue';

const mockData = {
  firstTimestamp: 1386435600000, // Sun Dec 08 2013 00:00:00
  secondTimestamp: 1600362123456, // Fri Sep 18 2020 00:02:03
};

describe('c-no-events-icon', () => {
  mockDateNow(mockData.secondTimestamp);

  const factory = generateRenderer(CNoEventsIcon, {
    attachTo: document.body,
  });

  it('Renders `c-no-events-icon` with default props correctly', () => {
    const wrapper = factory();

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `c-no-events-icon` with value correctly', () => {
    const wrapper = factory({
      propsData: {
        value: mockData.firstTimestamp,
      },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `c-no-events-icon` with custom props correctly', () => {
    const wrapper = factory({
      propsData: {
        value: mockData.firstTimestamp,
        color: 'secondary',
        maxWidth: 1120,
        top: true,
      },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
