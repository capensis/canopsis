import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import ServiceEntityInfoTab from '@/components/other/service/partials/service-entity-info-tab.vue';

const stubs = {
  'service-entity-actions': true,
  'service-entity-template': true,
};

const selectServiceEntityActions = wrapper => wrapper.find('service-entity-actions-stub');

describe('service-entity-info-tab', () => {
  const applyAction = jest.fn();
  const executeInstruction = jest.fn();

  const snapshotFactory = generateRenderer(ServiceEntityInfoTab, {

    stubs,
    listeners: {
      apply: applyAction,
      execute: executeInstruction,
    },
  });
  const factory = generateShallowRenderer(ServiceEntityInfoTab, {

    stubs,
    listeners: {
      apply: applyAction,
      execute: executeInstruction,
    },
  });

  test('Action applied after trigger entity actions', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        actions: [{}],
      },
    });

    selectServiceEntityActions(wrapper).triggerCustomEvent('apply');

    expect(applyAction).toHaveBeenCalled();
  });

  test('Instruction executed after trigger entity actions', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        actions: [{}],
      },
    });

    selectServiceEntityActions(wrapper).triggerCustomEvent('execute');

    expect(executeInstruction).toHaveBeenCalled();
  });

  test('Renders `service-entity-info-tab` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-entity-info-tab` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {
          _id: 'service-id',
          assigned_instructions: [{}],
        },
        template: '<div>Template</div>',
        actions: [{}],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
