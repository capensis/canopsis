import { generateRenderer } from '@unit/utils/vue';

import { LINK_RULE_ACTIONS } from '@/constants';

import CLinksList from '@/components/common/links/c-links-list.vue';

const snapshotStubs = {
  'c-copy-wrapper': true,
};

describe('c-links-list', () => {
  const links = {
    FirstCategory: [
      { label: 'FirstCategoryLinkLabel1', url: 'FirstCategoryLink1', rule_id: '1' },
      { label: 'FirstCategoryLinkLabel2', url: 'FirstCategoryLink2', rule_id: '2' },
      { label: 'FirstCategoryLinkLabel3', url: 'FirstCategoryLink3' },
      { label: 'FirstCategoryLinkLabel4', url: 'FirstCategoryLink4' },
    ],
    SecondCategory: [
      { label: 'SecondCategoryLinkLabel1', url: 'SecondCategoryLink1', rule_id: '2' },
    ],
    ThirdCategory: [
      { label: 'ThirdCategoryLinkLabel1', url: 'ThirdCategoryLink1', rule_id: '1' },
    ],
  };

  const snapshotFactory = generateRenderer(CLinksList, {

    stubs: snapshotStubs,
  });

  test('Renders `c-links-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-links-list` with links', () => {
    const wrapper = snapshotFactory({
      propsData: {
        links,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-links-list` with links and category', () => {
    const wrapper = snapshotFactory({
      propsData: {
        category: 'SecondCategory',
        links,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-links-list` with links with copy action', () => {
    const wrapper = snapshotFactory({
      propsData: {
        links: {
          ...links,

          FirstCategory: links.FirstCategory.map(link => ({
            ...link,
            action: LINK_RULE_ACTIONS.copy,
          })),
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
