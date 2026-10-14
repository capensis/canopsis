import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import ServiceEntityHeader from '@/components/other/service/partials/service-entity-header.vue';

const stubs = {
  'c-no-events-icon': true,
  'c-tooltip': true,
};

const selectAlert = wrapper => wrapper.find('v-alert-stub');

describe('service-entity-header', () => {
  const { authModule } = createAuthModule();
  const store = createMockedStoreModules([authModule]);

  const factory = generateShallowRenderer(ServiceEntityHeader, { stubs, store });

  test('Alert removed after trigger alert', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        lastActionUnavailable: true,
      },
    });

    const alert = selectAlert(wrapper);

    alert.triggerCustomEvent('input');

    expect(wrapper).toHaveBeenEmit('remove:unavailable');
  });
});
