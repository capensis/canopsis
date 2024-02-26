import { range } from 'lodash';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTimespanModule } from '@unit/utils/store';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import { createSelectInputStub } from '@unit/stubs/input';

import { MODALS, PBEHAVIOR_RRULE_PERIODS_RANGES } from '@/constants';

import PbehaviorRecurrenceRulePeriods from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-recurrence-rule-periods.vue';

const stubs = {
  'c-action-btn': true,
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'c-action-btn': true,
};

const selectTextField = wrapper => wrapper.find('select.v-select');
const selectActionButton = wrapper => wrapper.find('c-action-btn-stub');

describe('pbehavior-recurrence-rule-periods', () => {
  const nowTimestamp = 1386435600000;
  const $modals = mockModals();
  mockDateNow(nowTimestamp);

  const totalItems = 11;

  const pbehavior = {
    _id: 'pbehavior',
    type: {
      _id: 'pbehavior-type-id',
    },
  };

  const pbehaviorTimespans = range(totalItems).map(index => ({
    _id: `id-pbehavior-timespan-${index}`,
    from: 1614861788 + index,
    to: 1614862788 + index,
  }));
  const defaultTimespansData = {
    by_date: false,
    exceptions: [],
    exdates: [],
    rrule: undefined,
    start_at: undefined,
    type: pbehavior.type._id,
  };

  const { pbehaviorTimespanModule, fetchTimespansListWithoutStore } = createPbehaviorTimespanModule();

  const store = createMockedStoreModules([pbehaviorTimespanModule]);

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
  const snapshotFactory = generateRenderer(PbehaviorRecurrenceRulePeriods, {

    stubs: snapshotStubs,
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
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    expect(fetchTimespansListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          ...defaultTimespansData,
          view_from: 1385938800,
          view_to: 1386543599,
        },
      },
      undefined,
    );
  });

  test('Pbehavior timespans fetched after pbehavior updated', async () => {
    const wrapper = factory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    fetchTimespansListWithoutStore.mockClear();

    await wrapper.setProps({ pbehavior: { ...pbehavior } });

    expect(fetchTimespansListWithoutStore).toBeCalled();
  });

  test('Pbehavior timespans fetched after range changed to next week', async () => {
    const wrapper = factory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    fetchTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.triggerCustomEvent('input', PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek);
    await selectField.triggerCustomEvent('change');

    expect(fetchTimespansListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          ...defaultTimespansData,
          view_from: 1386543600,
          view_to: 1387148399,
        },
      },
      undefined,
    );
  });

  test('Pbehavior timespans fetched after range changed to next two week', async () => {
    const wrapper = factory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    fetchTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.triggerCustomEvent('input', PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks);
    await selectField.triggerCustomEvent('change');

    expect(fetchTimespansListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          ...defaultTimespansData,
          view_from: 1386543600,
          view_to: 1387753199,
        },
      },
      undefined,
    );
  });

  test('Pbehavior timespans fetched after range changed to this month', async () => {
    const wrapper = factory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    fetchTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.triggerCustomEvent('input', PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth);
    await selectField.triggerCustomEvent('change');

    expect(fetchTimespansListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          ...defaultTimespansData,
          view_from: 1385852400,
          view_to: 1388530799,
        },
      },
      undefined,
    );
  });

  test('Pbehavior timespans fetched after range changed to next month', async () => {
    const wrapper = factory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    fetchTimespansListWithoutStore.mockClear();

    const selectField = selectTextField(wrapper);

    await selectField.triggerCustomEvent('input', PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth);
    await selectField.triggerCustomEvent('change');

    expect(fetchTimespansListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          ...defaultTimespansData,
          view_from: 1388530800,
          view_to: 1391209199,
        },
      },
      undefined,
    );
  });

  test('Pbehavior recurrence rule modal opened after click action', async () => {
    const wrapper = factory({
      store,
      propsData: {
        pbehavior,
      },
    });

    const actionButton = selectActionButton(wrapper);
    actionButton.triggerCustomEvent('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorRecurrenceRule,
        config: {
          pbehavior,
        },
      },
    );
  });

  test('Renders `pbehavior-recurrence-rule-periods` without entities', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehavior-recurrence-rule-periods` with entities', async () => {
    fetchTimespansListWithoutStore.mockResolvedValueOnce(pbehaviorTimespans);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([pbehaviorTimespanModule]),
      propsData: {
        pbehavior,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
