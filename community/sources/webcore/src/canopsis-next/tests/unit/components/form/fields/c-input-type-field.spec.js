import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { PATTERN_INPUT_TYPES } from '@/constants';

import CInputTypeField from '@/components/forms/fields/c-input-type-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CInputTypeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CInputTypeField, {
  localVue,

  ...options,
});

const selectTextField = wrapper => wrapper.find('select.v-select');

describe('c-input-type-field', () => {
  it('Input type changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: PATTERN_INPUT_TYPES.string,
      },
    });

    const textField = selectTextField(wrapper);

    textField.vm.$emit('input', PATTERN_INPUT_TYPES.number);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(PATTERN_INPUT_TYPES.number);
  });

  it('Renders `c-input-type-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-input-type-field` with default custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: PATTERN_INPUT_TYPES.number,
        label: 'Custom label',
        name: 'name',
        disabled: true,
        flat: true,
        errorMessages: ['Message'],
        types: [
          { value: PATTERN_INPUT_TYPES.number },
          { value: PATTERN_INPUT_TYPES.string },
          { value: PATTERN_INPUT_TYPES.array },
          { value: PATTERN_INPUT_TYPES.null },
          { text: 'Custom boolean', value: PATTERN_INPUT_TYPES.boolean },
        ],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
