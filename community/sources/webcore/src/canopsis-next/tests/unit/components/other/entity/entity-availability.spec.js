import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAvailabilityModule, createMockedStoreModules } from '@unit/utils/store';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import EntityAvailability from '@/components/other/entity/entity-availability.vue';

const stubs = {
  'c-progress-overlay': true,
  'availability-bar': true,
};

const selectAvailabilityBar = wrapper => wrapper.find('availability-bar-stub');

describe('entity-availability', () => {
  jest.useFakeTimers({ now: 1386435500000 });

  const factory = generateShallowRenderer(EntityAvailability, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(EntityAvailability, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  const { availabilityModule, fetchAvailabilityWithoutStore } = createAvailabilityModule();

  const store = createMockedStoreModules([
    availabilityModule,
  ]);

  test('Availability fetched after mount', async () => {
    const entityId = Faker.datatype.string();

    factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    await flushPromises();

    expect(fetchAvailabilityWithoutStore).toHaveBeenCalledTimes(1);
    expect(fetchAvailabilityWithoutStore).toBeDispatchedWith({
      id: entityId,
      params: {
        from: 1386284400,
        to: 1386284400,
      },
    });
  });

  test('Availability fetched after trigger update interval with different time range', async () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    await flushPromises();

    fetchAvailabilityWithoutStore.mockClear();

    await selectAvailabilityBar(wrapper).triggerCustomEvent('update:interval', {
      from: QUICK_RANGES.last30Days.start,
      to: QUICK_RANGES.last30Days.stop,
    });

    expect(fetchAvailabilityWithoutStore).toBeDispatchedWith({
      id: entityId,
      params: {
        from: 1383865200,
        to: 1386370800,
      },
    });
  });

  test('Availability fetched after trigger update interval', async () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    await flushPromises();

    fetchAvailabilityWithoutStore.mockClear();

    await selectAvailabilityBar(wrapper).triggerCustomEvent('update:interval', {
      from: QUICK_RANGES.last6Months.start,
      to: QUICK_RANGES.last6Months.stop,
    });

    expect(fetchAvailabilityWithoutStore).toBeDispatchedWith({
      id: entityId,
      params: {
        from: 1370124000,
        to: 1385852400,
      },
    });
  });

  test('Fetching availability data for a different entity ID and time range', async () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    await flushPromises();

    fetchAvailabilityWithoutStore.mockClear();

    await selectAvailabilityBar(wrapper).triggerCustomEvent('update:interval', {
      from: QUICK_RANGES.last30Days.start,
      to: QUICK_RANGES.last30Days.stop,
    });

    expect(fetchAvailabilityWithoutStore).toBeDispatchedWith({
      id: entityId,
      params: {
        from: 1383865200,
        to: 1386370800,
      },
    });
  });

  test('Renders `entity-availability` with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {
          _id: 'entity-id',
        },
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `entity-availability` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {
          _id: 'entity-id',
        },
        defaultTimeRange: QUICK_RANGES.last6Months.value,
        defaultShowType: AVAILABILITY_SHOW_TYPE.duration,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
