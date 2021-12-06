import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import Title from '@/components/side-bars/settings/fields/common/title.vue';

const localVue = createVueInstance();

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
};

const factory = (options = {}) => shallowMount(Title, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Title, {
  localVue,
  parentComponent: {
    provide: {
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('title', () => {
  it('Value changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.string();

    textField.setValue(newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  it('Renders `title` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `title` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        title: 'Custom title',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
