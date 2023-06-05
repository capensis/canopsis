import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import ServiceEntityTreeOfDependenciesTab from '@/components/other/service/partials/service-entity-tree-of-dependencies-tab.vue';

const localVue = createVueInstance();

const stubs = {
  'service-dependencies': true,
};

describe('service-entity-tree-of-dependencies-tab', () => {
  const snapshotFactory = generateRenderer(ServiceEntityTreeOfDependenciesTab, {
    localVue,
    stubs,
  });

  test('Renders `service-entity-tree-of-dependencies-tab` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-entity-tree-of-dependencies-tab` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {
          _id: 'service-id',
        },
        columns: [{}],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
