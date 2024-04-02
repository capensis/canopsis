import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { fakeAlarmDetails } from '@unit/data/alarm';

import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimeline from '@/components/widgets/alarm/timeline/alarm-timeline.vue';

const stubs = {
  'c-enabled-field': true,
  'alarm-timeline-days': true,
  'c-pagination': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectPagination = wrapper => wrapper.find('c-pagination-stub');

describe('alarm-timeline', () => {
  const { steps: { data: steps, meta } } = fakeAlarmDetails();
  const factory = generateShallowRenderer(AlarmTimeline, { stubs });
  const snapshotFactory = generateRenderer(AlarmTimeline, { stubs });

  test('User can update group value and emit update event', () => {
    const newGroup = true;
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query: { group: false },
      },
    });

    const enabledField = selectEnabledField(wrapper);

    enabledField.triggerCustomEvent('input', newGroup);

    expect(wrapper).toEmit('update:query', { group: newGroup, page: 1 });
  });

  test('User can update page value and emit update event', () => {
    const query = { group: true };
    const newPage = 2;
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query,
      },
    });

    const pagination = selectPagination(wrapper);

    pagination.triggerCustomEvent('input', newPage);

    expect(wrapper).toEmit('update:query', { ...query, page: newPage });
  });

  test('Days computed property returns expected value', () => {
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query: {},
      },
    });

    const daysToCompare = Object.entries(groupAlarmSteps(steps)).reduce((acc, [key, items]) => {
      acc[key] = items.map(item => ({ ...item, key: expect.any(String) }));

      return acc;
    }, {});

    expect(wrapper.vm.days).toEqual(daysToCompare);
  });

  test('Renders `alarm-timeline` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        steps,
        meta,
        query: { group: true, page: 1 },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
