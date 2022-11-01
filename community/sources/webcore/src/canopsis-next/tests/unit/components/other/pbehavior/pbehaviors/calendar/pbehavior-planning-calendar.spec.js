import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTimespanModule, createPbehaviorTypesModule } from '@unit/utils/store';
import { pbehaviorToForm } from '@/helpers/forms/planning-pbehavior';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

const localVue = createVueInstance();

const snapshotStubs = {
  'c-progress-overlay': true,
  'pbehavior-create-event': true,
  'pbehavior-planning-calendar-legend': true,
};

describe('pbehavior-planning-calendar', () => {
  const { pbehaviorTimespanModule } = createPbehaviorTimespanModule();
  const { pbehaviorTypesModule } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([pbehaviorTimespanModule, pbehaviorTypesModule]);

  const pbehaviorsById = {
    pbehavior: pbehaviorToForm({
      _id: 'pbehavior',
      type: { _id: 'pbehavior-type-id' },
    }),
  };
  const addedPbehaviorsById = {
    'added-pbehavior': pbehaviorToForm({
      _id: 'added-pbehavior',
      type: { _id: 'added-pbehavior-type-id' },
    }),
  };
  const removedPbehaviorsById = {
    'removed-pbehavior': pbehaviorToForm({
      _id: 'removed-pbehavior',
      type: { _id: 'removed-pbehavior-type-id' },
    }),
  };
  const changedPbehaviorsById = {
    'changed-pbehavior': pbehaviorToForm({
      _id: 'changed-pbehavior',
      type: { _id: 'changed-pbehavior-type-id' },
    }),
  };

  const snapshotFactory = generateRenderer(PbehaviorPlanningCalendar, {
    localVue,
    stubs: snapshotStubs,
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
