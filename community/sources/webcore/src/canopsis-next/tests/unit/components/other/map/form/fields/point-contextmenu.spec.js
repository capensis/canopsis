import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import PointContextmenu from '@/components/other/map/form/fields/point-contextmenu.vue';

const selectListTiles = wrapper => wrapper.findAll('v-list-item-stub');
const selectListTileByIndex = (wrapper, index) => selectListTiles(wrapper).at(index);

describe('mermaid-contextmenu', () => {
  const factory = generateShallowRenderer(PointContextmenu);
  const snapshotFactory = generateRenderer(PointContextmenu);

  test('Add point event emitted after trigger click on add point', () => {
    const wrapper = factory({
      propsData: {
        positionX: 1,
        positionY: 2,
      },
    });

    const addPointTile = selectListTileByIndex(wrapper, 0);

    addPointTile.triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('add:point');
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

    editPointTile.triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('edit:point');
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

    removePointTile.triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('remove:point');
  });

  test('Renders `mermaid-contextmenu` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        positionX: 1,
        positionY: 2,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mermaid-contextmenu` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        positionX: 1,
        positionY: 2,
        editing: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
