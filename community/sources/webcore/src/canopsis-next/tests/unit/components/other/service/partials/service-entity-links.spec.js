import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { USERS_PERMISSIONS } from '@/constants';

import ServiceEntityLinks from '@/components/other/service/partials/service-entity-links.vue';

const localVue = createVueInstance();

describe('service-entity-links', () => {
  const links = [
    {
      cat_name: 'FirstCategory',
      links: [
        { label: 'FirstCategoryLinkLabel1', link: 'FirstCategoryLink1' },
        { label: 'FirstCategoryLinkLabel2', link: 'FirstCategoryLink2' },
      ],
    },
    {
      cat_name: 'SecondCategory',
      links: [
        { label: 'SecondCategoryLinkLabel1', link: 'SecondCategoryLink1' },
      ],
    },
    {
      cat_name: 'ThirdCategory',
      links: [
        { label: 'ThirdCategoryLinkLabel1', link: 'ThirdCategoryLink1' },
      ],
    },
  ];

  const { authModule, currentUserPermissionsById } = createAuthModule();

  const permission = USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks;

  const snapshotFactory = generateRenderer(ServiceEntityLinks, {
    localVue,
  });

  test('Renders `service-entity-links` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-links` with links', () => {
    currentUserPermissionsById.mockReturnValue({
      [`${permission}_FirstCategory`]: { actions: [] },
      [`${permission}_ThirdCategory`]: { actions: [] },
    });
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        links,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-links` with links and category', () => {
    currentUserPermissionsById.mockReturnValue({
      [permission]: { actions: [] },
    });
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        category: 'SecondCategory',
        links,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
