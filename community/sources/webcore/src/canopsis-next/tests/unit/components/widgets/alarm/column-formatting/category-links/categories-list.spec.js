import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import CategoriesList from '@/components/widgets/alarm/columns-formatting/category-links/categories-list.vue';

const localVue = createVueInstance();

const stubs = {
  'category-links': true,
};

describe('categories-list', () => {
  const snapshotFactory = generateRenderer(CategoriesList, {
    localVue,
    stubs,
    attachTo: document.body,
  });

  it('Renders `categories-list` with default props', async () => {
    snapshotFactory();

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `categories-list` with custom props', async () => {
    snapshotFactory({
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

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
