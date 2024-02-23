import { generateRenderer } from '@unit/utils/vue';

import MapBreadcrumbs from '@/components/widgets/map/partials/map-breadcrumbs.vue';

const selectButtons = wrapper => wrapper.findAll('.v-btn');
const selectButtonByIndex = (wrapper, index) => selectButtons(wrapper).at(index);

describe('map-breadcrumbs', () => {
  const snapshotFactory = generateRenderer(MapBreadcrumbs);

  test('Click emitted after button click was triggered', async () => {
    const previousMaps = [{}, {}, {}];
    const wrapper = snapshotFactory({
      propsData: {
        previousMaps,
      },
    });

    const index = 2;
    const button = selectButtonByIndex(wrapper, index);

    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('click', { index: 2 });
  });

  test('Renders `map-breadcrumbs` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `map-breadcrumbs` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pending: true,
        activeMap: {
          name: 'Active map',
        },
        previousMaps: [
          {
            name: 'First map',
          },
          {
            name: 'Second map',
          },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
