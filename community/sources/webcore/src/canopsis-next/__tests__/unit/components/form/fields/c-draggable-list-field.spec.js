import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CDraggableListField from '@/components/forms/fields/list/c-draggable-list-field.vue';

const stubs = {
  draggable: {
    template: `
      <div class="draggable" @change="$emit('change', $event)" />
    `,
  },
};

const selectDraggable = wrapper => wrapper.find('.draggable');

describe('c-draggable-list-field', () => {
  const items = [{ title: 'Item 1' }, { title: 'Item 2' }, { title: 'Item 3' }];

  const factory = generateShallowRenderer(CDraggableListField, { stubs });

  it('Item position changed after trigger draggable with moved event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const newOrder = [items[1], items[0], items[2]];

    selectDraggable(wrapper).triggerCustomEvent('input', newOrder);

    expect(wrapper).toEmitInput(newOrder);
  });

  it('Filter position changed after trigger draggable with added event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const newItem = {
      title: Faker.datatype.string(),
    };

    const newItems = [items[0], newItem, items[1], items[2]];

    selectDraggable(wrapper).triggerCustomEvent('input', newItems);

    expect(wrapper).toEmitInput(newItems);
  });

  it('Filter position changed after trigger draggable with removed event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const newItems = [items[0], items[2]];

    selectDraggable(wrapper).triggerCustomEvent('input', newItems);

    expect(wrapper).toEmitInput(newItems);
  });
});
