import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { uid } from '@/helpers/uid';

import CArrayTextField from '@/components/forms/fields/c-array-text-field.vue';

const mockData = {
  string: Faker.datatype.string(),
  number: Faker.datatype.number(),
  key1: uid(),
  key2: uid(),
  key3: uid(),
};

const stubs = {
  'c-action-btn': {
    template: `
      <button class="c-action-btn" @click="$emit('click')" />
    `,
  },
  'c-btn-with-error': {
    template: `
      <button class="c-btn-with-error" @click="$emit('click')">
        <slot />
      </button>
    `,
  },
  'c-text-field': createInputStub('c-text-field'),
};

const selectTextFields = wrapper => wrapper.findAll('input.c-text-field');
const selectAddButton = wrapper => wrapper.find('button.c-btn-with-error');
const selectRemoveButtons = wrapper => wrapper.findAll('button.c-action-btn');

describe('c-array-text-field', () => {
  const factory = generateShallowRenderer(CArrayTextField, { stubs });

  it('Empty string added after click on add button', () => {
    const initialValues = [{ key: mockData.key1, value: mockData.string }];
    const wrapper = factory({
      propsData: {
        values: initialValues,
      },
    });

    selectAddButton(wrapper).trigger('click');

    expect(wrapper).toEmit('change', [
      initialValues[0],
      expect.objectContaining({
        key: expect.any(String),
        value: '',
      }),
    ]);
  });

  it('Value changed after trigger mixed field', () => {
    const newFieldValue = Faker.datatype.string();
    const initialValues = [
      { key: mockData.key1, value: mockData.string },
      { key: mockData.key2, value: mockData.number },
    ];
    const wrapper = factory({
      propsData: {
        values: initialValues,
      },
    });

    selectTextFields(wrapper).at(0).triggerCustomEvent('update:modelValue', newFieldValue);

    expect(wrapper).toEmit('change', [
      { key: mockData.key1, value: newFieldValue },
      { key: mockData.key2, value: mockData.number },
    ]);
  });

  it('Field removed after click on remove button', () => {
    const initialValues = [
      { key: mockData.key1, value: mockData.string },
      { key: mockData.key2, value: mockData.number },
    ];
    const wrapper = factory({
      propsData: {
        values: initialValues,
      },
    });

    selectRemoveButtons(wrapper).at(1).trigger('click');

    expect(wrapper).toEmit('change', [{ key: mockData.key1, value: mockData.string }]);
  });
});
