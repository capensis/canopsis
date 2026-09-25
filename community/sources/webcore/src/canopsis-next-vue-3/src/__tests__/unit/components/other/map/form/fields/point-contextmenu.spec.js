import { generateShallowRenderer } from '@unit/utils/vue';

import PointContextmenu from '@/components/other/map/form/fields/point-contextmenu.vue';

const stubs = {
  'v-menu': {
    template: '<div><slot /></div>',
  },
  'v-list': {
    template: '<div><slot /></div>',
  },
  'v-list-item': {
    name: 'VListItem',
    emits: ['click'],
    template: '<button @click="$emit(\'click\')"><slot /></button>',
  },
  'v-list-item-title': {
    template: '<span><slot /></span>',
  },
};

const selectListItems = wrapper => wrapper.findAllComponents({ name: 'VListItem' });
const selectListItemByIndex = (wrapper, index) => selectListItems(wrapper).at(index);

describe('mermaid-contextmenu', () => {
  const factory = generateShallowRenderer(PointContextmenu, { stubs });

  test('Add point event emitted after trigger click on add point', async () => {
    const wrapper = factory({
      props: {
        positionX: 1,
        positionY: 2,
      },
    });

    await selectListItemByIndex(wrapper, 0).trigger('click');

    expect(wrapper.emitted('add:point')).toHaveLength(1);
  });

  test('Edit point event emitted after trigger click on edit point', async () => {
    const wrapper = factory({
      props: {
        positionX: 1,
        positionY: 2,
        editing: true,
      },
    });

    await selectListItemByIndex(wrapper, 0).trigger('click');

    expect(wrapper.emitted('edit:point')).toHaveLength(1);
  });

  test('Remove point event emitted after trigger click on remove point', async () => {
    const wrapper = factory({
      props: {
        positionX: 1,
        positionY: 2,
        editing: true,
      },
    });

    await selectListItemByIndex(wrapper, 1).trigger('click');

    expect(wrapper.emitted('remove:point')).toHaveLength(1);
  });
});
