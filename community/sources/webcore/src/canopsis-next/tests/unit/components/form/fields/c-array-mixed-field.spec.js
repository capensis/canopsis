import Faker from 'faker';
import { Validator } from 'vee-validate';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import CArrayTextField from '@/components/forms/fields/c-array-text-field.vue';

const mockData = {
  string: Faker.datatype.string(),
  number: Faker.datatype.number(),
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
    const wrapper = factory({
      propsData: {
        values: [mockData.string],
      },
    });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper).toEmit('change', [mockData.string, '']);
  });

  it('Value changed after trigger mixed field', () => {
    const newFieldValue = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        values: [mockData.string, mockData.number],
      },
    });
    const secondFieldElement = wrapper.findAll('v-layout-stub').at(0);

    selectTextField(secondFieldElement).setValue(newFieldValue);

    expect(wrapper).toEmit('change', [newFieldValue, mockData.number]);
  });

  it('Field removed after click on remove button', () => {
    const wrapper = factory({
      propsData: {
        values: [mockData.string, mockData.number],
      },
    });
    const secondFieldElement = wrapper.findAll('v-layout-stub').at(1);

    secondFieldElement.find('button.c-action-btn').trigger('click');

    expect(wrapper).toEmit('change', [mockData.string]);
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
          'string',
          123,
          false,
          null,
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
        values: ['string', 123],
        disabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
