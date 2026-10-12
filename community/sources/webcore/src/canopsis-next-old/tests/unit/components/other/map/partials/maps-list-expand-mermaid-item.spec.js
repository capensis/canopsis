import { generateRenderer } from '@unit/utils/vue';

import MapsListExpandMermaidItem from '@/components/other/map/partials/maps-list-expand-mermaid-item.vue';

const stubs = {
  'mermaid-preview': true,
};

describe('maps-list-expand-mermaid-item', () => {
  const map = {
    name: 'Map',
    parameters: {},
  };

  const snapshotFactory = generateRenderer(MapsListExpandMermaidItem, { stubs });

  test('Renders `maps-list-expand-mermaid-item` with map', async () => {
    const wrapper = snapshotFactory({
      propsData: { map },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
