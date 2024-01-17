import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';

import { USERS_PERMISSIONS } from '@/constants';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import ServiceEntityTemplate from '@/components/other/service/partials/service-entity-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
  'c-links-list': true,
};

describe('service-entity-template', () => {
  const { authModule, currentUserPermissionsById } = createAuthModule();
  const store = createMockedStoreModules([authModule]);

  const snapshotFactory = generateRenderer(ServiceEntityTemplate, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `service-entity-template` with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {},
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-template` with custom props', async () => {
    currentUserPermissionsById.mockReturnValue({
      [USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks]: { actions: [] },
    });

    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {
          _id: 'service-id',
          links: { test: [{ rule_id: 'id', url: 'url', label: 'label' }] },
        },
        template: '{{entity._id}}{{links category="test"}}',
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-template` with custom props without right', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: {
          _id: 'service-id',
          links: { test: [{ rule_id: 'id', url: 'url', label: 'label' }] },
        },
        template: '{{entity._id}}{{links category="test"}}',
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
