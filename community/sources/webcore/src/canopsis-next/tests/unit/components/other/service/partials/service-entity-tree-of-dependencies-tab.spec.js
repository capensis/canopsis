import { generateRenderer } from '@unit/utils/vue';

import ServiceEntityTreeOfDependenciesTab from '@/components/other/service/partials/service-entity-tree-of-dependencies-tab.vue';

const stubs = {
  'service-dependencies': true,
};

describe('service-entity-tree-of-dependencies-tab', () => {
  const snapshotFactory = generateRenderer(ServiceEntityTreeOfDependenciesTab, {

    stubs,
  });

  test('Renders `service-entity-tree-of-dependencies-tab` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
