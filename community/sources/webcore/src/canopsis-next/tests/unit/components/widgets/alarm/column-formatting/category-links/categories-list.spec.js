import { mount, createVueInstance } from '@unit/utils/vue';

import CategoriesList from '@/components/widgets/alarm/columns-formatting/category-links/categories-list.vue';

const localVue = createVueInstance();

const stubs = {
  'category-links': true,
};

const snapshotFactory = (options = {}) => mount(CategoriesList, {
  localVue,
  stubs,

  ...options,
});

const selectMenuContent = wrapper => wrapper.find('.v-menu__content');

describe('categories-list', () => {
  it('Renders `categories-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `categories-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        categories: [
          {
            category: 'Category',
            links: [{ link: 'link', label: 'Label' }],
          },
          {
            category: 'Category2',
            links: [{ link: 'link2', label: 'Label 2' }],
          },
        ],
        limit: 1,
      },
    });

    const menuContent = selectMenuContent(wrapper);

    expect(menuContent.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });
});
