import Faker from 'faker';
import { Validator } from 'vee-validate';

import { shallowMount, createVueInstance } from '@unit/utils/vue';

import CArrayMixedField from '@/components/forms/fields/c-array-mixed-field.vue';

const localVue = createVueInstance();

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
  'c-mixed-field': {
    props: ['value'],
    template: `
      <input
        :value="value"
        class="c-mixed-field"
        @input="$listeners.input($event.target.value)"
      />
    `,
  },
};

const snapshotStubs = {
  'c-mixed-field': true,
  'c-action-btn': true,
};

const factory = (options = {}) => shallowMount(CArrayMixedField, {
  localVue,
  stubs,
  ...options,
});

describe('c-array-mixed-field', () => {
  it('Empty string added after click on add button', () => {
    const wrapper = factory({
      propsData: {
        values: [mockData.string],
      },
    });
    const addButtonElement = wrapper.find('button.v-btn');

    addButtonElement.trigger('click');

    const changeEvents = wrapper.emitted('change');
    expect(changeEvents).toHaveLength(1);

    const [newFieldsData] = changeEvents[0];
    const [oldValue, newValue] = newFieldsData;

    expect(oldValue).toEqual(mockData.string);
    expect(newValue).toEqual('');
  });

  it('Value changed after trigger mixed field', () => {
    const newFieldValue = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        values: [mockData.string, mockData.number],
      },
    });
    const secondFieldElement = wrapper.findAll('v-layout-stub').at(0);
    const mixedFieldElement = secondFieldElement.find('input.c-mixed-field');

    mixedFieldElement.setValue(newFieldValue);

    const changeEvents = wrapper.emitted('change');
    expect(changeEvents).toHaveLength(1);

    const [newFieldsData] = changeEvents[0];

    expect(newFieldsData).toEqual([newFieldValue, mockData.number]);
  });

  it('Field removed after click on remove button', () => {
    const wrapper = factory({
      propsData: {
        values: [mockData.string, mockData.number],
      },
    });
    const secondFieldElement = wrapper.findAll('v-layout-stub').at(1);
    const removeButtonElement = secondFieldElement.find('button.c-action-btn');

    removeButtonElement.trigger('click');

    const changeEvents = wrapper.emitted('change');
    expect(changeEvents).toHaveLength(1);

    const [newFieldsData] = changeEvents[0];

    expect(newFieldsData).toEqual([mockData.string]);
  });

  it('Renders `c-array-mixed-field` with default props correctly', () => {
    const wrapper = shallowMount(CArrayMixedField, {
      localVue,
      stubs: snapshotStubs,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-array-mixed-field` with all field types correctly', () => {
    const wrapper = shallowMount(CArrayMixedField, {
      localVue,
      provide: {
        $validator: new Validator(),
      },
      stubs: snapshotStubs,
      propsData: {
        values: [
          'string',
          123,
          false,
          null,
          [null, 'string', 123],
        ],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders disabled `c-array-mixed-field` correctly', () => {
    const wrapper = shallowMount(CArrayMixedField, {
      localVue,
      provide: {
        $validator: new Validator(),
      },
      stubs: snapshotStubs,
      propsData: {
        values: ['string', 123],
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
