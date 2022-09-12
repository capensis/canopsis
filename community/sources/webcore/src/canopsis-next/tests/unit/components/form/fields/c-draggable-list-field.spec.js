import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CDraggableListField from '@/components/forms/fields/c-draggable-list-field.vue';

const localVue = createVueInstance();

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

const factory = (options = {}) => shallowMount(CDraggableListField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CDraggableListField, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectDraggable = wrapper => wrapper.find('.draggable');

describe('c-draggable-list-field', () => {
  const items = [
    { title: 'Item 1' },
    { title: 'Item 2' },
    { title: 'Item 3' },
  ];

  it('Item position changed after trigger draggable with moved event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const draggable = selectDraggable(wrapper);

    draggable.trigger('change', {
      moved: {
        oldIndex: 0,
        newIndex: 1,
      },
    });

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];

    expect(eventData).toEqual([
      items[1],
      items[0],
      items[2],
    ]);
  });

  it('Filter position changed after trigger draggable with added event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const draggable = selectDraggable(wrapper);

    const newItem = {
      title: Faker.datatype.string(),
    };

    draggable.trigger('change', {
      added: {
        element: newItem,
        newIndex: 1,
      },
    });

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];

    expect(eventData).toEqual([
      items[0],
      newItem,
      items[1],
      items[2],
    ]);
  });

  it('Filter position changed after trigger draggable with removed event', () => {
    const wrapper = factory({
      propsData: {
        value: items,
      },
    });

    const draggable = selectDraggable(wrapper);

    draggable.trigger('change', {
      removed: {
        oldIndex: 1,
      },
    });

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];

    expect(eventData).toEqual([
      items[0],
      items[2],
    ]);
  });

  it('Renders `c-draggable-list-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
