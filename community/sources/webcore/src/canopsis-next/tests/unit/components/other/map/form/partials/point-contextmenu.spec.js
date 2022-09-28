import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import PointContextmenu from '@/components/other/map/form/partials/point-contextmenu.vue';

const localVue = createVueInstance();

const factory = (options = {}) => shallowMount(PointContextmenu, {
  localVue,

  ...options,
});

const snapshotFactory = (options = {}) => mount(PointContextmenu, {
  localVue,

  ...options,
});

const selectListTiles = wrapper => wrapper.findAll('v-list-tile-stub');
const selectListTileByIndex = (wrapper, index) => selectListTiles(wrapper).at(index);

describe('mermaid-contextmenu', () => {
  test('Add point event emitted after trigger click on add point', () => {
    const wrapper = factory({
      propsData: {
        positionX: 1,
        positionY: 2,
      },
    });

    const addPointTile = selectListTileByIndex(wrapper, 0);

    addPointTile.vm.$emit('click');

    expect(wrapper).toEmit('add:point');
  });

  test('Edit point event emitted after trigger click on edit point', () => {
    const wrapper = factory({
      propsData: {
        positionX: 1,
        positionY: 2,
        editing: true,
      },
    });

    const editPointTile = selectListTileByIndex(wrapper, 0);

    editPointTile.vm.$emit('click');

    expect(wrapper).toEmit('edit:point');
  });

  test('Remove point event emitted after trigger click on remove point', () => {
    const wrapper = factory({
      propsData: {
        positionX: 1,
        positionY: 2,
        editing: true,
      },
    });

    const removePointTile = selectListTileByIndex(wrapper, 1);

    removePointTile.vm.$emit('click');

    expect(wrapper).toEmit('remove:point');
  });

  test('Renders `mermaid-contextmenu` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        positionX: 1,
        positionY: 2,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `mermaid-contextmenu` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        positionX: 1,
        positionY: 2,
        editing: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
