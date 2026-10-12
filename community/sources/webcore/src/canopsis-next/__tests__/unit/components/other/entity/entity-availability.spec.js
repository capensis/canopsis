import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { QUICK_RANGES } from '@/constants';

import EntityAvailability from '@/components/other/entity/entity-availability.vue';

const stubs = { 'c-progress-overlay': true, 'availability-bar': true };

describe('entity-availability', () => {
  vi.useFakeTimers({ now: 1386435500000 });

  const factory = generateShallowRenderer(EntityAvailability, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Builds availability query for the default interval', () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    expect(wrapper.vm.getQuery()).toMatchObject({ _id: entityId });
    expect(wrapper.vm.getQuery().from).toEqual(expect.any(Number));
    expect(wrapper.vm.getQuery().to).toEqual(expect.any(Number));
  });

  test('Updates query after interval change', () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    const interval = {
      from: QUICK_RANGES.last30Days.start,
      to: QUICK_RANGES.last30Days.stop,
    };
    wrapper.vm.updateInterval(interval);

    expect(wrapper.vm.query.interval).toEqual(interval);
  });

  test('Builds query for a six-month interval', () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    wrapper.vm.updateInterval({
      from: QUICK_RANGES.last6Months.start,
      to: QUICK_RANGES.last6Months.stop,
    });

    expect(wrapper.vm.getQuery()).toMatchObject({ _id: entityId });
    expect(wrapper.vm.getQuery().from).toEqual(expect.any(Number));
    expect(wrapper.vm.getQuery().to).toEqual(expect.any(Number));
  });

  test('Uses the current entity ID in the query', () => {
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        entity: {
          _id: entityId,
        },
        defaultTimeRange: QUICK_RANGES.yesterday.value,
      },
    });

    expect(wrapper.vm.getQuery()._id).toBe(entityId);
  });
});
