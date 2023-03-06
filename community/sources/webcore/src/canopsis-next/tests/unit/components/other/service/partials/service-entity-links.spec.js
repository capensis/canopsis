import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { USERS_PERMISSIONS } from '@/constants';

import ServiceEntityLinks from '@/components/other/service/partials/service-entity-links.vue';

const localVue = createVueInstance();

describe('service-entity-links', () => {
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

  const { authModule, currentUserPermissionsById } = createAuthModule();

  const permission = USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks;

  const snapshotFactory = generateRenderer(ServiceEntityLinks, {
    localVue,
  });

  test('Renders `service-entity-links` with default props', () => {
    currentUserPermissionsById.mockReturnValue({});

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-links` with links', () => {
    currentUserPermissionsById.mockReturnValue({ [permission]: { actions: [] } });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        links,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-links` with links and category', () => {
    currentUserPermissionsById.mockReturnValue({ [permission]: { actions: [] } });

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
