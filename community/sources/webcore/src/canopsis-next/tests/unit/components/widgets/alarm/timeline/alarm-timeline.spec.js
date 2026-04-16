import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { fakeAlarmDetails } from '@unit/data/alarm';

import { ALARM_STEPS_TYPES } from '@/constants';

import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimeline from '@/components/widgets/alarm/timeline/alarm-timeline.vue';

const stubs = {
  'c-enabled-field': true,
  'alarm-timeline-days': true,
  'c-pagination': true,
};

const selectGroupField = wrapper => wrapper.findAll('c-enabled-field-stub').at(0);
const selectOnlyCommentsField = wrapper => wrapper.findAll('c-enabled-field-stub').at(1);
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

    const groupField = selectGroupField(wrapper);

    groupField.triggerCustomEvent('input', newGroup);

    expect(wrapper).toEmit('update:query', { group: newGroup, page: 1 });
  });

  test('User can enable only comments filter and emit update event', () => {
    const query = { group: true, page: 2 };
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query,
      },
    });

    const onlyCommentsField = selectOnlyCommentsField(wrapper);

    onlyCommentsField.triggerCustomEvent('input', true);

    expect(wrapper).toEmit('update:query', {
      group: true,
      page: 1,
      type: ALARM_STEPS_TYPES.comment,
    });
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

    expect(wrapper.vm.days).toEqual(groupAlarmSteps(steps));
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
