import { flushPromises, generateRenderer } from '@unit/utils/vue';

import AlarmPbehaviorCounters from '@/components/widgets/service-weather/alarm-pbehavior-counters.vue';

describe('alarm-pbehavior-counters', () => {
  const counters = [{
    count: 1,
    type: {
      _id: 'type-1-id',
      name: 'type-1-name',
      icon_name: 'assignment_returned',
    },
  }, {
    count: 3,
    type: {
      _id: 'type-2-id',
      name: 'type-2-name',
      icon_name: 'add_photo_alternate',
    },
  }, {
    count: 5,
    type: {
      _id: 'type-3-id',
      name: 'type-3-name',
      icon_name: 'assistant_direction',
    },
  }, {
    count: 7,
    type: {
      _id: 'type-4-id',
      name: 'type-4-name',
      icon_name: 'assignment',
    },
  }];
  const types = ['type-1-id', 'type-2-id'];
  const snapshotFactory = generateRenderer(AlarmPbehaviorCounters, {

    attachTo: document.body,
  });

  test('Renders `alarm-pbehavior-counters` with default props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  test('Renders `alarm-pbehavior-counters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        types,
        counters,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
