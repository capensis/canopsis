import { mount, createVueInstance } from '@unit/utils/vue';

import CategoryLinks from '@/components/widgets/alarm/columns-formatting/category-links/category-links.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CategoryLinks, {
  localVue,

  ...options,
});

describe('category-links', () => {
  it('Renders `category-links` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `category-links` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        links: [
          { link: 'https://custom-link.link/1', label: 'Custom link' },
          { link: 'https://custom-link.link/2', label: 'Custom link 2' },
        ],
        label: 'Custom label',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
