import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';

import CMovableCardIteratorField from '@/components/forms/fields/card-iterator/c-movable-card-iterator-field.vue';
import MovableCardIteratorItem from '@/components/forms/fields/card-iterator/movable-card-iterator-item.vue';

const stubs = {
  'movable-card-iterator-item': MovableCardIteratorItem,
  'v-btn': createButtonStub('v-btn'),
};

const selectCards = wrapper => wrapper.findAll('v-card-stub');
const selectCardByIndex = (wrapper, index) => selectCards(wrapper).at(index);
const selectCardButtons = wrapper => wrapper.findAll('.v-btn');
const selectCardUpButton = wrapper => selectCardButtons(wrapper).at(0);
const selectCardDownButton = wrapper => selectCardButtons(wrapper).at(1);
const selectCardRemoveButton = wrapper => selectCardButtons(wrapper).at(2);

describe('c-movable-card-iterator-field', () => {
  const items = [
    { key: 'key-1', title: 'title-1' },
    { key: 'key-2', title: 'title-2' },
    { key: 'key-3', title: 'title-3' },
    { key: 'key-4', title: 'title-4' },
  ];

  const factory = generateShallowRenderer(CMovableCardIteratorField, { stubs });
  const snapshotFactory = generateRenderer(CMovableCardIteratorField);

  test('Card moved below after click on down button', async () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const movingCard = selectCardByIndex(wrapper, 1);

    selectCardDownButton(movingCard).trigger('click');

    expect(wrapper).toEmit('input', [
      items[0],
      items[2],
      items[1],
      items[3],
    ]);
  });

  test('Card doesn\'t moved below after try to move down last card', async () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const movingCard = selectCardByIndex(wrapper, 3);

    selectCardDownButton(movingCard).trigger('click');

    expect(wrapper).not.toEmit('input');
  });

  test('Card moved above after click on up button', async () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const movingCard = selectCardByIndex(wrapper, 3);

    selectCardUpButton(movingCard).trigger('click');

    expect(wrapper).toEmit('input', [
      items[0],
      items[1],
      items[3],
      items[2],
    ]);
  });

  test('Card doesn\'t moved above after try to move up first card', async () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const movingCard = selectCardByIndex(wrapper, 0);

    selectCardUpButton(movingCard).trigger('click');

    expect(wrapper).not.toEmit('input');
  });

  test('Card removed above after click on remove button', async () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const removingCard = selectCardByIndex(wrapper, 2);

    selectCardRemoveButton(removingCard).trigger('click');

    expect(wrapper).toEmit('input', [
      items[0],
      items[1],
      items[3],
    ]);
  });

  test('Renders `c-movable-card-iterator-field` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-movable-card-iterator-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [
          { id: 'id-1', text: 'text-1' },
          { id: 'id-2', text: 'text-2' },
        ],
        itemKey: 'id',
        addable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-movable-card-iterator-field` with slots props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [
          { id: 'id-1', text: 'text-1' },
          { id: 'id-2', text: 'text-2' },
        ],
        itemKey: 'id',
        addable: true,
      },
      slots: {
        prepend: '<div class="prepend-slot" />',
        append: '<div class="append-slot" />',
      },
      scopedSlots: {
        item(props) {
          return this.$createElement(
            'div',
            { attrs: { class: 'item-slot' } },
            JSON.stringify(props),
          );
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
