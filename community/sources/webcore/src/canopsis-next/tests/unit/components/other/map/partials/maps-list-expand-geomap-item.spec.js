import { mount, createVueInstance } from '@unit/utils/vue';

import MapsListExpandGeomapItem from '@/components/other/map/partials/maps-list-expand-geomap-item.vue';

const localVue = createVueInstance();

const stubs = {
  'geomap-preview': true,
};

const snapshotFactory = (options = {}) => mount(MapsListExpandGeomapItem, {
  localVue,
  stubs,

  ...options,
});

describe('maps-list-expand-geomap-item', () => {
  const map = {
    name: 'Map',
    parameters: {},
  };

  test('Renders `maps-list-expand-geomap-item` with map', async () => {
    const wrapper = snapshotFactory({
      propsData: { map },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
