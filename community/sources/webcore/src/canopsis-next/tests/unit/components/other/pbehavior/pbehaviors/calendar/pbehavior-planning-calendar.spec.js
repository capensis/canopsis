import flushPromises from 'flush-promises';
import { Day } from 'dayspan';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTimespanModule, createPbehaviorTypesModule } from '@unit/utils/store';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { pbehaviorToForm } from '@/helpers/forms/planning-pbehavior';
import DaySpanVuetifyPlugin from '@/plugins/dayspan-vuetify';
import { convertDateToMoment } from '@/helpers/date/date';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

const localVue = createVueInstance();

localVue.use(DaySpanVuetifyPlugin);

const stubs = {
  'c-action-btn': true,
  'c-progress-overlay': true,
  'pbehavior-create-event': true,
  'pbehavior-planning-calendar-legend': true,
};

describe('pbehavior-planning-calendar', () => {
  mockDateNow(1386435500000);

  const now = new Day(convertDateToMoment(1386435500000));
  const today = now.start();
  const tomorrow = now.next();

  localVue.$dayspan.now = now;
  localVue.$dayspan.today = today;
  localVue.$dayspan.tomorrow = tomorrow;

  const { pbehaviorTimespanModule, fetchTimespansListWithoutStore } = createPbehaviorTimespanModule();
  const { pbehaviorTypesModule } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([pbehaviorTimespanModule, pbehaviorTypesModule]);

  const pbehaviorsById = {
    pbehavior: pbehaviorToForm({
      _id: 'pbehavior',
      type: { _id: 'pbehavior-type-id' },
      tstart: 1612825000,
    }),
  };
  const addedPbehaviorsById = {
    'added-pbehavior': pbehaviorToForm({
      _id: 'added-pbehavior',
      type: { _id: 'added-pbehavior-type-id' },
      tstart: 1612825200,
    }),
  };
  const removedPbehaviorsById = {
    'removed-pbehavior': pbehaviorToForm({
      _id: 'removed-pbehavior',
      type: { _id: 'removed-pbehavior-type-id' },
      tstart: 1612825200,
    }),
  };
  const changedPbehaviorsById = {
    'changed-pbehavior': pbehaviorToForm({
      _id: 'changed-pbehavior',
      type: { _id: 'changed-pbehavior-type-id' },
      tstart: 1612825100,
    }),
  };
  const timespans = [
    { from: 1612825000, to: 1612911400, type: { _id: 'type-id-1' } },
    { from: 1612825100, to: 1612911500, type: { _id: 'type-id-2' } },
    { from: 1612825200, to: 1612911600, type: { _id: 'type-id-3' } },
  ];

  const snapshotFactory = generateRenderer(PbehaviorPlanningCalendar, {
    localVue,
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `pbehavior-planning-calendar` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        pbehaviorsById: {},
        addedPbehaviorsById: {},
        removedPbehaviorsById: {},
        changedPbehaviorsById: {},
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-planning-calendar` with custom props', async () => {
    fetchTimespansListWithoutStore.mockResolvedValue(timespans);
    const wrapper = snapshotFactory({
      store,
      propsData: {
        pbehaviorsById,
        addedPbehaviorsById,
        removedPbehaviorsById,
        changedPbehaviorsById,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
