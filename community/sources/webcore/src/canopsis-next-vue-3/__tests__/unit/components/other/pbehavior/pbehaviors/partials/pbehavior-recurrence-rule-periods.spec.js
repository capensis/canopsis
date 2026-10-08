import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { MODALS, PBEHAVIOR_RRULE_PERIODS_RANGES } from '@/constants';

import PbehaviorRecurrenceRulePeriods from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-recurrence-rule-periods.vue';

const pbehaviorTimespan = vi.hoisted(() => ({
  fetchPbehaviorTimespansListWithoutStore: vi.fn().mockResolvedValue([]),
}));

vi.mock('@/hooks/store/modules/pbehavior-timespan', () => ({
  usePbehaviorTimespan: () => pbehaviorTimespan,
}));

vi.mock('@/hooks/store/modules/info', () => ({
  useInfo: () => ({ timezone: 'Asia/Almaty' }),
}));

const stubs = {
  'c-action-btn': true,
  'c-select': {
    props: ['modelValue', 'items'],
    emits: ['update:modelValue'],
    template: `
      <select
        class="c-select"
        :value="modelValue"
        @change="$emit('update:modelValue', $event.target.value)"
      >
        <option v-for="item in items" :key="item.value" :value="item.value">
          {{ item.text }}
        </option>
      </select>
    `,
  },
};

const selectTextField = wrapper => wrapper.find('select.c-select');
const selectActionButton = wrapper => wrapper.find('c-action-btn-stub');

describe('pbehavior-recurrence-rule-periods', () => {
  const nowTimestamp = 1386435600000;
  const $modals = mockModals();
  mockDateNow(nowTimestamp);

  const pbehavior = {
    _id: 'pbehavior',
    type: {
      _id: 'pbehavior-type-id',
    },
  };

  const defaultTimespansData = {
    by_date: false,
    exceptions: [],
    exdates: [],
    rrule: undefined,
    start_at: undefined,
    type: pbehavior.type._id,
  };

  const factory = generateShallowRenderer(PbehaviorRecurrenceRulePeriods, {
    stubs,
    mocks: { $modals },
    parentComponent: {
      provide: {
        $system: {
          timezone: process.env.TZ,
        },
      },
    },
  });

  test('Pbehavior timespans fetched after mount', async () => {
    factory({
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    expect(pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore).toHaveBeenCalledWith({
      data: {
        ...defaultTimespansData,
        view_from: 1385920800,
        view_to: 1386525599,
      },
    });
  });

  test('Pbehavior timespans fetched after pbehavior updated', async () => {
    const wrapper = factory({
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore.mockClear();

    await wrapper.setProps({ pbehavior: { ...pbehavior } });

    expect(pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore).toBeCalled();
  });

  test('Pbehavior timespans fetched after range changed to next week', async () => {
    const wrapper = factory({
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.setValue(PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek);

    expect(pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore).toHaveBeenCalledWith({
      data: {
        ...defaultTimespansData,
        view_from: 1386525600,
        view_to: 1387130399,
      },
    });
  });

  test('Pbehavior timespans fetched after range changed to next two week', async () => {
    const wrapper = factory({
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.setValue(PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks);

    expect(pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore).toHaveBeenCalledWith({
      data: {
        ...defaultTimespansData,
        view_from: 1386525600,
        view_to: 1387735199,
      },
    });
  });

  test('Pbehavior timespans fetched after range changed to this month', async () => {
    const wrapper = factory({
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.setValue(PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth);

    expect(pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore).toHaveBeenCalledWith({
      data: {
        ...defaultTimespansData,
        view_from: 1385834400,
        view_to: 1388512799,
      },
    });
  });

  test('Pbehavior timespans fetched after range changed to next month', async () => {
    const wrapper = factory({
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.setValue(PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth);

    expect(pbehaviorTimespan.fetchPbehaviorTimespansListWithoutStore).toHaveBeenCalledWith({
      data: {
        ...defaultTimespansData,
        view_from: 1388512800,
        view_to: 1391191199,
      },
    });
  });

  test('Pbehavior recurrence rule modal opened after click action', async () => {
    const wrapper = factory({
      propsData: {
        pbehavior,
      },
    });

    const actionButton = selectActionButton(wrapper);
    actionButton.triggerCustomEvent('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.pbehaviorRecurrenceRule,
      config: {
        pbehavior,
      },
    });
  });
});
