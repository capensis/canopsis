import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import KpiRatingSettings from '@/components/other/kpi/rating-settings/kpi-rating-settings.vue';

const stubs = {
  'kpi-rating-settings-list': true,
};

describe('kpi-rating-settings', () => {
  const factory = generateShallowRenderer(KpiRatingSettings, { stubs });

  it('Rating settings fetched after mount', async () => {
    const fetchRatingSettings = vi.fn();
    factory({
      store: createMockedStoreModules([
        {
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
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);
    expect(fetchRatingSettings).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        limit: 10,
        page: 1,
      },
    });
  });

  it('Rating settings fetched after change query', async () => {
    const fetchRatingSettings = vi.fn();
    const initialItemsPerPage = Faker.datatype.number();
    const wrapper = factory({
      data() {
        return {
          query: {
            itemsPerPage: initialItemsPerPage,
          },
        };
      },
      store: createMockedStoreModules([
        {
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
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
    });

    await flushPromises();

    fetchRatingSettings.mockReset();

    const kpiRatingSettingsListElement = wrapper.find('kpi-rating-settings-list-stub');

    const itemsPerPage = Faker.datatype.number({ max: initialItemsPerPage });
    const page = Faker.datatype.number();

    kpiRatingSettingsListElement.triggerCustomEvent('update:options', {
      itemsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);
    expect(fetchRatingSettings).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        limit: itemsPerPage,
        page,
      },
    });
  });
});
