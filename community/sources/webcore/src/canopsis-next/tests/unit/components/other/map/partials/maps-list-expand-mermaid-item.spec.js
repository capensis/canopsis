import { mount, createVueInstance } from '@unit/utils/vue';

import MapsListExpandMermaidItem from '@/components/other/map/partials/maps-list-expand-mermaid-item.vue';

const localVue = createVueInstance();

const stubs = {
  'mermaid-preview': true,
};

const snapshotFactory = (options = {}) => mount(MapsListExpandMermaidItem, {
  localVue,
  stubs,

  ...options,
});

describe('maps-list-expand-mermaid-item', () => {
  const map = {
    name: 'Map',
    parameters: {},
  };

  test('Renders `maps-list-expand-mermaid-item` with map', async () => {
    const wrapper = snapshotFactory({
      propsData: { map },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
