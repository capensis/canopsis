import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import Title from '@/components/sidebars/form/fields/title.vue';

const stubs = {
  'widget-settings-item': true,
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('title', () => {
  const factory = generateShallowRenderer(Title, { stubs });
  const snapshotFactory = generateRenderer(Title, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

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

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `title` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        title: 'Custom title',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
