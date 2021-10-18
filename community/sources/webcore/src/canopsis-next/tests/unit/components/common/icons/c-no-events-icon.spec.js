import { shallowMount, createVueInstance } from '@unit/utils/vue';
import { stubDateNow } from '@unit/utils/stub-hooks';

import CNoEventsIcon from '@/components/common/icons/c-no-events-icon.vue';

const localVue = createVueInstance();

const mockData = {
  timestamp: 1386435600000, // Sun Dec 08 2013 00:00:00
  now: 1600362123456, // Fri Sep 18 2020 00:02:03
};

describe('c-no-events-icon', () => {
  stubDateNow(mockData.now);

  it('Renders `c-no-events-icon` with default props correctly', () => {
    const wrapper = shallowMount(CNoEventsIcon, { localVue });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-no-events-icon` with value correctly', () => {
    const wrapper = shallowMount(CNoEventsIcon, {
      localVue,
      propsData: {
        value: mockData.timestamp,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-no-events-icon` with custom props correctly', () => {
    const wrapper = shallowMount(CNoEventsIcon, {
      localVue,
      propsData: {
        value: mockData.timestamp,
        color: 'secondary',
        maxWidth: 1120,
        top: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
