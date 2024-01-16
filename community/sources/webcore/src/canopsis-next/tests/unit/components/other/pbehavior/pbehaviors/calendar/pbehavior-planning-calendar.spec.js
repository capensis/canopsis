import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTimespanModule, createPbehaviorTypesModule } from '@unit/utils/store';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

const localVue = createVueInstance();

const stubs = {
  'c-action-btn': true,
  'c-progress-overlay': true,
  'c-calendar': {
    template: `
      <div class="c-calendar">
        <slot name="form-event" :close="close" :event="event" />
        <slot name="menu-right" />
      </div>
    `,
    computed: {
      event() {
        return {};
      },
    },
    methods: {
      close: jest.fn(),
    },
  },
  'pbehavior-create-event': true,
  'pbehavior-planning-calendar-legend': true,
};

describe('pbehavior-planning-calendar', () => {
  jest.useFakeTimers({ now: 1386435500000 });

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

    expect(wrapper).toMatchSnapshot();
  });
});
