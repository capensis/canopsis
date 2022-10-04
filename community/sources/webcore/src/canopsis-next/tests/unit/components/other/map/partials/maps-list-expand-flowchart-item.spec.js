import { mount, createVueInstance } from '@unit/utils/vue';

import MapsListExpandFlowchartItem from '@/components/other/map/partials/maps-list-expand-flowchart-item.vue';

const localVue = createVueInstance();

const stubs = {
  'flowchart-preview': true,
};

const snapshotFactory = (options = {}) => mount(MapsListExpandFlowchartItem, {
  localVue,
  stubs,

  ...options,
});

describe('maps-list-expand-flowchart-item', () => {
  const map = {
    name: 'FLowchart map',
    parameters: {},
  };

  test('Renders `maps-list-expand-flowchart-item` with map', async () => {
    const wrapper = snapshotFactory({
      propsData: { map },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
