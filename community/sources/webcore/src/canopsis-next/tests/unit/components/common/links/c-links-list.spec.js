import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import CLinksList from '@/components/common/links/c-links-list.vue';

const localVue = createVueInstance();

describe('c-links-list', () => {
  const links = {
    FirstCategory: [
      { label: 'FirstCategoryLinkLabel1', url: 'FirstCategoryLink1', rule_id: '1' },
      { label: 'FirstCategoryLinkLabel2', url: 'FirstCategoryLink2', rule_id: '2' },
    ],
    SecondCategory: [
      { label: 'SecondCategoryLinkLabel1', url: 'SecondCategoryLink1', rule_id: '2' },
    ],
    ThirdCategory: [
      { label: 'ThirdCategoryLinkLabel1', url: 'ThirdCategoryLink1', rule_id: '1' },
    ],
  };

  const snapshotFactory = generateRenderer(CLinksList, {
    localVue,
  });

  test('Renders `c-links-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-links-list` with links', () => {
    const wrapper = snapshotFactory({
      propsData: {
        links,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-links-list` with links and category', () => {
    const wrapper = snapshotFactory({
      propsData: {
        category: 'SecondCategory',
        links,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
