import { generateRenderer } from '@unit/utils/vue';

import MapsListExpandFlowchartItem from '@/components/other/map/partials/maps-list-expand-flowchart-item.vue';

const stubs = {
  'flowchart-preview': true,
};

describe('maps-list-expand-flowchart-item', () => {
  const map = {
    name: 'FLowchart map',
    parameters: {},
  };

  const snapshotFactory = generateRenderer(MapsListExpandFlowchartItem, { stubs });

  test('Renders `maps-list-expand-flowchart-item` with map', async () => {
    const wrapper = snapshotFactory({
      propsData: { map },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
