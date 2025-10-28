import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createEntityModule, createMockedStoreModules } from '@unit/utils/store';
import { fakeAlarm } from '@unit/data/alarm';

import AlarmsExpandPanelEntityEnrichments from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-entity-enrichments.vue';

const stubs = {
  'entity-infos-logs-list': true,
};

const snapshotStubs = {
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

  const {
    entityModule,
    fetchEntityInfosLogsListWithoutStore,
  } = createEntityModule();

  const store = createMockedStoreModules([entityModule]);

  const factory = generateShallowRenderer(AlarmsExpandPanelEntityEnrichments, { stubs, store });
  const snapshotFactory = generateRenderer(AlarmsExpandPanelEntityEnrichments, { stubs: snapshotStubs, store });

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
      expect.any(Object),
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

  test('Renders `alarms-expand-panel-entity-enrichments` with default props', () => {
    const alarm = fakeAlarm();

    fetchEntityInfosLogsListWithoutStore.mockResolvedValue({
      data: mockItems,
      meta: mockMeta,
    });

    const wrapper = snapshotFactory({
      propsData: {
        alarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-expand-panel-entity-enrichments` with empty data', () => {
    const alarm = fakeAlarm();

    fetchEntityInfosLogsListWithoutStore.mockResolvedValue({
      data: [],
      meta: { total_count: 0 },
    });

    const wrapper = snapshotFactory({
      propsData: {
        alarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-expand-panel-entity-enrichments` with minimal alarm prop', () => {
    const minimalAlarm = { entity: { _id: 'test-id' } };

    fetchEntityInfosLogsListWithoutStore.mockResolvedValue({
      data: [],
      meta: { total_count: 0 },
    });

    const wrapper = snapshotFactory({
      propsData: {
        alarm: minimalAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
