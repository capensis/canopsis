import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import ServiceEntityTemplate from '@/components/other/service/partials/service-entity-template.vue';

const localVue = createVueInstance();

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'service-entity-links': true,
};

describe('service-entity-template', () => {
  const snapshotFactory = generateRenderer(ServiceEntityTemplate, {
    localVue,
    stubs,
  });

  test('Renders `service-entity-template` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-template` with custom props', async () => {
    const wrapper = snapshotFactory({
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
