import { mount, createVueInstance } from '@unit/utils/vue';

import AlarmColumnValueLinks from '@/components/widgets/alarm/columns-formatting/alarm-column-value-links.vue';

const localVue = createVueInstance();

const stubs = {
  'category-links': true,
};

const snapshotFactory = (options = {}) => mount(AlarmColumnValueLinks, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-column-value-links', () => {
  it('Renders `alarm-column-value-links` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-links` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        links: [
          { category: 'Category', links: [] },
        ],
        label: 'Custom label',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
