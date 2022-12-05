import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreModules } from '@unit/utils/store';
import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
} from '@/constants';
import AlarmColumnValueCategories from '@/components/widgets/alarm/columns-formatting/alarm-column-value-categories.vue';

const localVue = createVueInstance();

const stubs = {
  'categories-list': true,
};

const selectMenuContent = wrapper => wrapper.find('.v-menu__content');

describe('alarm-column-value-categories', () => {
  const links = {
    Category: [
      { link: 'Category link', label: 'Category link' },
      { link: 'Category link 2', label: 'Category link 2' },
      { link: 'Category link 3', label: 'Category link 3' },
    ],
    Category2: [
      { link: 'Category 2 link', label: 'Category 2 link' },
    ],
    Category3: [
      'Category 3 link',
    ],
  };
  const prefix = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links];
  const authModule = {
    name: 'auth',
    getters: {
      currentUserPermissionsById: () => Object
        .keys(links)
        .reduce((acc, category) => {
          acc[`${prefix}_2${category}`] = {
            actions: [],
          };

          return acc;
        }, {
          [BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links]]: {
            actions: [],
          },
        }),
    },
  };

  const snapshotFactory = generateRenderer(AlarmColumnValueCategories, {
    localVue,
    stubs,
  });

  it('Renders `alarm-column-value-categories` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = selectMenuContent(wrapper);

    expect(menuContent.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-categories` with custom props', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        links,
        asList: true,
        limit: 2,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-categories` with custom props without access to categories', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'auth',
        getters: { currentUserPermissionsById: {} },
      }]),
      propsData: {
        links,
        limit: 2,
      },
    });

    const menuContent = selectMenuContent(wrapper);

    expect(menuContent.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });
});
