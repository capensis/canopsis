import { generateShallowRenderer } from '@unit/utils/vue';

import ServiceEntityInfoTab from '@/components/other/service/partials/service-entity-info-tab.vue';

const stubs = {
  'service-entity-actions': true,
  'service-entity-template': true,
};

const selectServiceEntityActions = wrapper => wrapper.find('service-entity-actions-stub');

describe('service-entity-info-tab', () => {
  const factory = generateShallowRenderer(ServiceEntityInfoTab, {
    stubs,
  });

  test('Action applied after trigger entity actions', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        actions: [{}],
      },
    });

    selectServiceEntityActions(wrapper).triggerCustomEvent('apply');

    expect(wrapper).toHaveBeenEmit('apply');
  });

  test('Instruction executed after trigger entity actions', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        actions: [{}],
      },
    });

    selectServiceEntityActions(wrapper).triggerCustomEvent('execute');

    expect(wrapper).toHaveBeenEmit('execute');
  });
});
