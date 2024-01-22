import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CDraggableListField from '@/components/forms/fields/list/c-draggable-list-field.vue';

const stubs = {
  draggable: {
    template: `
      <div class="draggable" @change="$listeners.change($event)" />
    `,
  },
};

const snapshotStubs = {
  draggable: true,
};

const selectDraggable = wrapper => wrapper.find('.draggable');

describe('c-draggable-list-field', () => {
  const items = [
    { title: 'Item 1' },
    { title: 'Item 2' },
    { title: 'Item 3' },
  ];

  const factory = generateShallowRenderer(CDraggableListField, { stubs });
  const snapshotFactory = generateRenderer(CDraggableListField, { stubs: snapshotStubs });

  it('Item position changed after trigger draggable with moved event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const newOrder = [
      items[1],
      items[0],
      items[2],
    ];

    selectDraggable(wrapper).triggerCustomEvent('input', newOrder);

    expect(wrapper).toEmit('input', newOrder);
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

    const newItems = [
      items[0],
      newItem,
      items[1],
      items[2],
    ];

    selectDraggable(wrapper).triggerCustomEvent('input', newItems);

    expect(wrapper).toEmit('input', newItems);
  });

  it('Filter position changed after trigger draggable with removed event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const newItems = [
      items[0],
      items[2],
    ];

    selectDraggable(wrapper).triggerCustomEvent('input', newItems);

    expect(wrapper).toEmit('input', newItems);
  });

  it('Renders `c-draggable-list-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-draggable-list-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: items,
        disabled: true,
        animation: 2,
        component: 'span',
      },
      slots: {
        default: '<div class="custom-class">Custom slot</div>',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
