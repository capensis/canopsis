import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';

import { DEFAULT_SERVICE_WEATHER_MODAL_TEMPLATE } from '@/constants';

import ServiceTemplate from '@/components/other/service/partials/service-template.vue';

const localVue = createVueInstance();

const stubs = {
  'c-progress-overlay': true,
  'service-entities-list': true,
  'c-table-pagination': true,
};

const snapshotFactory = (options = {}) => mount(ServiceTemplate, {
  localVue,
  stubs,

  ...options,
});

describe('service-template', () => {
  const service = {};

  const { authModule } = createAuthModule();

  const store = createMockedStoreModules([
    authModule,
  ]);

  test('Renders `service-template` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        service,
        pagination: {},
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-template` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        service,
        pagination: {},
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate: DEFAULT_SERVICE_WEATHER_MODAL_TEMPLATE,
        },
        itemsPerPage: 10,
        totalItems: 20,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
