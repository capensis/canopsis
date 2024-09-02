import { flushPromises, createVueInstance, generateRenderer } from '@unit/utils/vue';
import { Day } from 'dayspan';

import { createMockedStoreModules, createPbehaviorTimespanModule, createPbehaviorTypesModule } from '@unit/utils/store';
import { mockDateNow } from '@unit/utils/mock-hooks';

import { convertDateToMoment } from '@/helpers/date/date';

import DaySpanVuetifyPlugin from '@/plugins/dayspan-vuetify';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

const localVue = createVueInstance();

localVue.use(DaySpanVuetifyPlugin);

const stubs = {
  'c-action-btn': true,
  'c-progress-overlay': true,
  'pbehavior-create-event': true,
  'pbehavior-planning-calendar-legend': true,
  'calendar-app-period-picker': true,
};

describe('pbehavior-planning-calendar', () => {
  mockDateNow(1386435500000);

  const now = new Day(convertDateToMoment(1386435500000));
  const today = now.start();
  const tomorrow = now.next();

  localVue.$dayspan.now = now;
  localVue.$dayspan.today = today;
  localVue.$dayspan.tomorrow = tomorrow;

  const { pbehaviorTimespanModule } = createPbehaviorTimespanModule();
  const { pbehaviorTypesModule } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([pbehaviorTimespanModule, pbehaviorTypesModule]);

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
});
