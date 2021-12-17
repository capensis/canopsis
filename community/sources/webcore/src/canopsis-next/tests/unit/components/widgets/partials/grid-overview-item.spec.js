import { mount, createVueInstance } from '@unit/utils/vue';

import GridOverviewItem from '@/components/widgets/partials/grid-overview-item.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(GridOverviewItem, {
  localVue,

  ...options,
});

describe('grid-overview-item', () => {
  const widget = {
    grid_parameters: {
      mobile: {
        x: 1,
        y: 1,
        h: 2,
        w: 12,
        autoHeight: false,
      },
      tablet: {
        x: 2,
        y: 3,
        h: 3,
        w: 10,
        autoHeight: true,
      },
      desktop: {
        x: 4,
        y: 22,
        h: 12,
        w: 8,
        autoHeight: false,
      },
    },
  };

  it.each(['m', 't', 'l', 'xl'])('Renders `grid-overview-item` with %s desktop size', (size) => {
    const wrapper = snapshotFactory({
      propsData: {
        widget,
      },
      mocks: {
        $mq: size,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
