import Faker from 'faker';
import { Validator } from 'vee-validate';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

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
      <button class="c-action-btn" @click="$listeners.click" />
    `,
  },
  'v-btn': {
    template: `
      <button class="v-btn" @click="$listeners.click" >
        <slot></slot>
      </button>
    `,
  },
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'c-action-btn': true,
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-array-text-field', () => {
  const factory = generateShallowRenderer(CArrayTextField, { stubs });
  const snapshotFactory = generateRenderer(CArrayTextField, { stubs: snapshotStubs });

  it('Empty string added after click on add button', () => {
    const initialValues = [{ key: mockData.key1, value: mockData.string }];
    const wrapper = factory({
      propsData: {
        values: initialValues,
      },
    });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper).toEmit('change', [
      initialValues[0],
      expect.objectContaining({
        key: expect.string(),
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
    const firstFieldElement = wrapper.findAll('v-layout-stub').at(1);

    selectTextField(firstFieldElement).setValue(newFieldValue);

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
    const secondFieldElement = wrapper.findAll('v-layout-stub').at(2);

    secondFieldElement.find('button.c-action-btn').trigger('click');

    expect(wrapper).toEmit('change', [{ key: mockData.key1, value: mockData.string }]);
  });

  it('Renders `c-array-text-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-array-text-field` with all field types correctly', () => {
    const wrapper = snapshotFactory({
      provide: {
        $validator: new Validator(),
      },
      propsData: {
        values: [
          { key: 'test-key-1', value: 'string' },
          { key: 'test-key-2', value: 123 },
          { key: 'test-key-3', value: false },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders disabled `c-array-text-field` correctly', () => {
    const wrapper = snapshotFactory({
      provide: {
        $validator: new Validator(),
      },
      propsData: {
        values: [
          { key: 'test-key-1', value: 'string' },
          { key: 'test-key-2', value: 123 },
        ],
        disabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
