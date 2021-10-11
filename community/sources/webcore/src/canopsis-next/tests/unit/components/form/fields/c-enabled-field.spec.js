import { mount, createVueInstance } from '@/unit/utils/vue';

import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';

const localVue = createVueInstance();

const directives = {
  field: {
    bind: (el, directive, vnode) => {
      const { event } = vnode.context.$options.model;
      el.setAttribute('value', `${JSON.stringify(directive.value)}`);

      const updateField = $event => vnode.context.$emit(event || 'input', $event.target.value);

      el.addEventListener('change', updateField);
    },
  },
};

const stubs = {
  'v-switch': {
    props: ['value'],
    template: `
      <input :checked="value" type="checkbox" class="v-switch" @input="$emit('input')" />
    `,
  },
};

const factory = (options = {}) => mount(CEnabledField, {
  localVue,
  stubs,
  directives,
  ...options,
});

describe('c-enabled-field', () => {
  it('Value set in the input', () => {
    const wrapper = factory({ propsData: { value: false } });
    const input = wrapper.find('input.v-switch');

    expect(input.element.checked).toBe(false);
  });

  it('Value changed in the input', () => {
    const wrapper = factory({ propsData: { value: false } });
    const input = wrapper.find('input.v-switch');

    input.setChecked(true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([true]);
  });

  it('Renders `c-enabled-field` with default props correctly', () => {
    const wrapper = mount(CEnabledField, {
      localVue,
      directives,
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
