import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { CRUD_ACTIONS, USERS_PERMISSIONS } from '@/constants';

import KpiRatingSettings from '@/components/other/kpi/rating-settings/kpi-rating-settings.vue';

const stubs = {
  'kpi-rating-settings-list': true,
};

const defaultStore = createMockedStoreModules([{
  name: 'ratingSettings',
  getters: {
    pending: false,
    items: [],
    meta: {
      total_count: 0,
    },
  },
  actions: {
    fetchList: jest.fn(),
  },
}, {
  name: 'auth',
  getters: {
    currentUserPermissionsById: {},
  },
}]);

describe('kpi-rating-settings', () => {
  const factory = generateShallowRenderer(KpiRatingSettings, { stubs });
  const snapshotFactory = generateRenderer(KpiRatingSettings, { stubs });

  it('Rating settings fetched after mount', async () => {
    const fetchRatingSettings = jest.fn();
    factory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: fetchRatingSettings,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);
    expect(fetchRatingSettings).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: 10,
          page: 1,
        },
      },
      undefined,
    );
  });

  it('Rating settings fetched after change query', async () => {
    const fetchRatingSettings = jest.fn();
    const initialItemsPerPage = Faker.datatype.number();
    const wrapper = factory({
      data() {
        return {
          query: {
            itemsPerPage: initialItemsPerPage,
          },
        };
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: fetchRatingSettings,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
    });

    await flushPromises();

    fetchRatingSettings.mockReset();

    const kpiRatingSettingsListElement = wrapper.find('kpi-rating-settings-list-stub');

    const itemsPerPage = Faker.datatype.number({ max: initialItemsPerPage });
    const page = Faker.datatype.number();

    kpiRatingSettingsListElement.vm.$emit('update:options', {
      itemsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);
    expect(fetchRatingSettings).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: itemsPerPage,
          page,
        },
      },
      undefined,
    );
  });

  it('Renders `kpi-rating-settings` with default props', () => {
    const wrapper = snapshotFactory({
      store: defaultStore,
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `kpi-rating-settings` with full permissions', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 10,
          },
        },
        actions: {
          fetchList: jest.fn(),
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {
            [USERS_PERMISSIONS.technical.kpiRatingSettings]: {
              actions: [
                CRUD_ACTIONS.update,
                CRUD_ACTIONS.read,
              ],
            },
          },
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });
});
