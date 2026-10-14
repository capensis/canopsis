import { generateShallowRenderer } from '@unit/utils/vue';

import CCardIteratorField from '@/components/forms/fields/card-iterator/c-card-iterator-field.vue';

const stubs = {
  'c-draggable-list-field': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: `
      <div class="c-draggable-list-field">
        <div v-for="(item, index) in modelValue" :key="item.key">
          <slot name="item" :element="item" :index="index" />
        </div>
      </div>
    `,
  },
};

const selectDraggableListField = wrapper => wrapper.find('.c-draggable-list-field');

describe('c-movable-card-iterator-field', () => {
  const items = [
    { key: 'key-1', title: 'title-1' },
    { key: 'key-2', title: 'title-2' },
    { key: 'key-3', title: 'title-3' },
    { key: 'key-4', title: 'title-4' },
  ];

  const factory = generateShallowRenderer(CCardIteratorField, { stubs });

  test('Cards reordered after trigger draggable list', async () => {
    const wrapper = factory({
      propsData: {
        data: items,
      },
    });

    const newOrder = [items[0], items[2], items[1], items[3]];

    await selectDraggableListField(wrapper).triggerCustomEvent('update:modelValue', newOrder);

    expect(wrapper).toEmitInput(newOrder);
  });

  test('Cards reordered to first position after trigger draggable list', async () => {
    const wrapper = factory({
      propsData: {
        data: items,
      },
    });

    const newOrder = [items[3], items[0], items[1], items[2]];

    await selectDraggableListField(wrapper).triggerCustomEvent('update:modelValue', newOrder);

    expect(wrapper).toEmitInput(newOrder);
  });
});
