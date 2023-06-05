import { mount, createVueInstance } from '@unit/utils/vue';

import MapBreadcrumbs from '@/components/widgets/map/partials/map-breadcrumbs.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(MapBreadcrumbs, {
  localVue,

  ...options,
});

const selectButtons = wrapper => wrapper.findAll('.v-btn');
const selectButtonByIndex = (wrapper, index) => selectButtons(wrapper).at(index);

describe('map-breadcrumbs', () => {
  test('Click emitted after button click was triggered', async () => {
    const previousMaps = [{}, {}, {}];
    const wrapper = snapshotFactory({
      propsData: {
        previousMaps,
      },
    });

    const index = 2;
    const button = selectButtonByIndex(wrapper, index);

    await button.vm.$emit('click');

    expect(wrapper).toEmit('click', { index: 2 });
  });

  test('Renders `map-breadcrumbs` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
