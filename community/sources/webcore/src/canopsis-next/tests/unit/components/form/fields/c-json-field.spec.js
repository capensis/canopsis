import { shallowMount, createVueInstance } from '@unit/utils/vue';

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
    const input = wrapper.find('.v-textarea textarea');

    expect(input.element.value).toBe(JSON.stringify(value, undefined, 4));
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

    const input = wrapper.find('.v-textarea textarea');

    expect(input.element.value).toBe(value);
  });

  /*  it('Object value set to the input', () => {
    const value = { key: 'value' };
    const wrapper = factory({ propsData: { value } });
    const input = wrapper.find('.v-textarea textarea');

    expect(input.element.value).toBe(value);
  }); */

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
