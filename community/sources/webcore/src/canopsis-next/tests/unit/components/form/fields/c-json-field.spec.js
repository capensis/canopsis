import Faker from 'faker';
import { Validator } from 'vee-validate';

import { shallowMount, createVueInstance, mount } from '@unit/utils/vue';

import CJsonField from '@/components/forms/fields/c-json-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-textarea': {
    props: ['value'],
    template: `
      <div class='v-textarea'>
        <textarea :value="value" @input="$listeners.input($event.target.value)" />
        <slot name="append" />
      </div>

    `,
  },
  'v-tooltip': {
    template: `
      <div class='v-tooltip'>
        <slot />
      </div>
    `,
  },
  'v-btn': {
    template: `
      <button class="v-btn" @click="$listeners.click">
        <slot />
      </button>
    `,
  },
};

const factory = (options = {}) => shallowMount(CJsonField, {
  localVue,
  stubs,
  ...options,
});

describe('c-json-field', () => {
  it('Object value set to the input', () => {
    const value = { key: 'value' };
    const wrapper = factory({ propsData: { value } });

    const textarea = wrapper.find('.v-textarea textarea');

    expect(textarea.element.value).toBe(JSON.stringify(value, undefined, 4));
  });

  it('Input event after set value', () => {
    const value = { key: 'value' };
    const wrapper = factory({ propsData: { value } });

    const textarea = wrapper.find('.v-textarea textarea');

    textarea.setValue('{ "newKey": "newValue" }');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
  });

  it('String value set to the input with variables', () => {
    const value = `{
    "extra_vars": {
        "version": 15,
        "state": {{ .Alarm.Value.StateValue }},
        "entity_id": {{ .Entity.ID }}
    }
}`;

    const wrapper = factory({
      propsData: {
        value,
        variables: true,
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    expect(textarea.element.value).toBe(value);
  });

  it('v-validate works correctly on valid json', async () => {
    const name = Faker.datatype.string();
    const value = { key: 'value' };
    const validator = new Validator();

    mount({
      inject: ['$validator'],
      components: {
        CJsonField,
      },
      props: ['name', 'value'],
      template: `
        <c-json-field :name="name" :value="value" />
      `,
    }, {
      localVue,
      stubs,
      provide: {
        $validator: validator,
      },
      propsData: {
        value,
        name,
      },
    });

    const isValid = await validator.validateAll();

    expect(isValid).toBeTruthy();
    expect(validator.fields.find({ name })).toBeTruthy();
  });

  it('v-validate works correctly on invalid json', async () => {
    const name = Faker.datatype.string();
    const value = { key: 'value' };
    const validator = new Validator();

    const wrapper = mount({
      inject: ['$validator'],
      components: {
        CJsonField,
      },
      props: ['name', 'value'],
      template: `
        <c-json-field :name="name" :value="value" />
      `,
    }, {
      localVue,
      stubs,
      provide: {
        $validator: validator,
      },
      propsData: {
        value,
        name,
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    textarea.setValue('asd');

    const isValid = await validator.validateAll();

    expect(isValid).toBeTruthy();
    expect(validator.fields.find({ name })).toBeTruthy();
  });

  it('Renders `c-json-field` with default props correctly', () => {
    const wrapper = shallowMount(CJsonField, {
      localVue,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly', () => {
    const wrapper = shallowMount(CJsonField, {
      localVue,
      propsData: {
        validateOn: 'button',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly and touched value', async () => {
    const wrapper = mount({
      components: {
        CJsonField,
      },
      props: ['name', 'value'],
      template: `
        <c-json-field ref="jsonField" :name="name" :value="value" validate-on="button" />
      `,
    }, {
      localVue,
      $_veeValidate: {
        validator: 'new',
      },
      propsData: {
        value: '{}',
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue('asd');

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly', () => {
    const wrapper = shallowMount(CJsonField, {
      localVue,
      propsData: {
        validateOn: 'button',
        readonly: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly', () => {
    const wrapper = shallowMount(CJsonField, {
      localVue,
      propsData: {
        value: { key: 'value' },
        validateOn: 'button',
        helpText: 'Test text',
        label: 'JSON field',
        name: 'jsonField',
        rows: 4,
        autoGrow: true,
        box: true,
        outline: true,
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
