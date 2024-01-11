import { generateRenderer } from '@unit/utils/vue';

import MapsListExpandGeomapItem from '@/components/other/map/partials/maps-list-expand-geomap-item.vue';

const stubs = {
  'geomap-preview': true,
};

describe('maps-list-expand-geomap-item', () => {
  const map = {
    name: 'Map',
    parameters: {},
  };

  const snapshotFactory = generateRenderer(MapsListExpandGeomapItem, { stubs });

  test('Renders `maps-list-expand-geomap-item` with map', async () => {
    const wrapper = snapshotFactory({
      propsData: { map },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
