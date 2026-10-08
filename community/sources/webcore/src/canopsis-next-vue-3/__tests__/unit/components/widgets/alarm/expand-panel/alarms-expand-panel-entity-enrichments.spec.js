import { fakeAlarm } from '@unit/data/alarm';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import AlarmsExpandPanelEntityEnrichments from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-entity-enrichments.vue';

const { fetchEntityInfosLogsListWithoutStore } = vi.hoisted(() => ({
  fetchEntityInfosLogsListWithoutStore: vi.fn(),
}));

vi.mock('@/hooks/store/modules/entity', () => ({
  useEntity: () => ({ fetchEntityInfosLogsListWithoutStore }),
}));

const stubs = {
  'entity-infos-logs-list': true,
};

describe('alarms-expand-panel-entity-enrichments', () => {
  const mockItems = [
    {
      _id: Faker.datatype.string(),
      time: Faker.date.recent(),
      rule: {
        _id: Faker.datatype.string(),
        description: Faker.lorem.sentence(),
        action_type: Faker.random.word(),
        action_description: Faker.lorem.sentence(),
      },
      name: Faker.random.word(),
      prev_value: Faker.random.word(),
      new_value: Faker.random.word(),
    },
    {
      _id: Faker.datatype.string(),
      time: Faker.date.recent(),
      rule: {
        _id: Faker.datatype.string(),
        description: Faker.lorem.sentence(),
        action_type: Faker.random.word(),
        action_description: Faker.lorem.sentence(),
      },
      name: Faker.random.word(),
      prev_value: Faker.random.word(),
      new_value: Faker.random.word(),
    },
  ];

  const mockMeta = {
    total_count: 10,
    page: 1,
    page_count: 1,
  };

  const factory = generateShallowRenderer(AlarmsExpandPanelEntityEnrichments, { stubs });

  afterEach(() => {
    fetchEntityInfosLogsListWithoutStore.mockClear();
  });

  test('Should call fetchEntityInfosLogsListWithoutStore on mount', async () => {
    const alarm = fakeAlarm();

    fetchEntityInfosLogsListWithoutStore.mockResolvedValue({
      data: mockItems,
      meta: mockMeta,
    });

    factory({
      propsData: {
        alarm,
      },
    });

    await flushPromises();

    expect(fetchEntityInfosLogsListWithoutStore).toHaveBeenCalledWith(
      expect.objectContaining({
        params: expect.objectContaining({
          _id: alarm.entity._id,
          page: 1,
          limit: 10,
        }),
      }),
    );
  });

  test('Should render entity-infos-logs-list component', () => {
    const alarm = fakeAlarm();

    fetchEntityInfosLogsListWithoutStore.mockResolvedValue({
      data: mockItems,
      meta: mockMeta,
    });

    const wrapper = factory({
      propsData: {
        alarm,
      },
    });

    const entityInfosLogsList = wrapper.find('entity-infos-logs-list-stub');

    expect(entityInfosLogsList.exists()).toBe(true);
  });
});
