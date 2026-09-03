import { flushPromises, mount } from '@vue/test-utils';
import { createPinia } from 'pinia';
import { describe, expect, it } from 'vitest';
import { defineComponent, ref } from 'vue';

import veeValidatePlugin from '@/plugins/vee-validate';

import { i18n } from '@/i18n';

import { useForm } from '@/hooks/form';

import CNumberField from '@/components/forms/fields/compat/c-number-field.vue';

import {
  useVuetifyFieldValidation,
  validationEmits,
  validationProps,
} from './validation';

const TestField = defineComponent({
  props: validationProps,
  emits: validationEmits,
  template: `
    <label>
      <input
        :name="inputName"
        :value="inputValue"
        @input="inputValue = $event.target.value"
        @blur="onBlur"
      >
      <span v-for="message in errorMessages" :key="message">{{ message }}</span>
    </label>
  `,
  setup(props, { emit, attrs }) {
    return useVuetifyFieldValidation(props, emit, attrs);
  },
});

const TestForm = defineComponent({
  components: { TestField },
  template: `
    <form @submit.prevent="submit">
      <test-field v-model="name" label="Name" name="name" required :validation-delay="0" />
      <button type="submit">Submit</button>
    </form>
  `,
  setup() {
    const form = useForm();
    const name = ref('');
    const submit = form.handleSubmit(() => {});

    return {
      name,
      submit,
    };
  },
});

const TestFieldWithExternalErrors = defineComponent({
  components: { TestField },
  template: `
    <test-field
      v-model="name"
      label="Name"
      name="name"
      required
      :error-messages="[]"
      :validation-delay="0"
    />
  `,
  setup() {
    const name = ref('');

    return {
      name,
    };
  },
});

const TestMaxLengthField = defineComponent({
  components: { TestField },
  template: `
    <test-field
      v-model="description"
      label="Description"
      name="description"
      :max-length="5"
      :validation-delay="0"
    />
  `,
  setup() {
    const description = ref('');

    return {
      description,
    };
  },
});

const TestNumberField = defineComponent({
  components: { CNumberField },
  props: {
    optimalValue: {
      type: [Number, String],
      default: 2,
    },
  },
  template: `
    <c-number-field
      v-model:value="minimal"
      name="minimal"
      label="Minimal"
      :max="optimalValue"
      :min="0"
      :error-messages="[]"
      :validation-delay="0"
    />
  `,
  setup() {
    const minimal = ref(1);

    return {
      minimal,
    };
  },
});

const mountTestForm = () =>
  mount(TestForm, {
    global: {
      plugins: [createPinia(), i18n, [veeValidatePlugin, { i18n }]],
    },
  });

describe('useVuetifyFieldValidation', () => {
  it('does not show required errors after blurring an unchanged empty field', async () => {
    const wrapper = mountTestForm();

    await wrapper.find('input[name="name"]').trigger('blur');
    await flushPromises();

    expect(wrapper.text()).not.toContain('The field is required');
  });

  it('shows required errors after submitting an unchanged empty field', async () => {
    const wrapper = mountTestForm();

    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(wrapper.text()).toContain('The field is required');
  });

  it('shows internal errors for dirty fields when external errors are empty', async () => {
    const wrapper = mount(TestFieldWithExternalErrors, {
      global: {
        plugins: [createPinia(), i18n, [veeValidatePlugin, { i18n }]],
      },
    });

    const input = wrapper.find('input[name="name"]');

    await input.setValue('test');
    await input.setValue('');
    await flushPromises();

    expect(wrapper.text()).toContain('The field is required');
  });

  it('shows internal errors for dirty number fields when external errors are empty', async () => {
    const wrapper = mount(TestNumberField, {
      global: {
        components: { CTextField: TestField },
        plugins: [createPinia(), i18n, [veeValidatePlugin, { i18n }]],
      },
    });

    await wrapper.find('input[name="minimal"]').setValue('3');
    await flushPromises();

    expect(wrapper.text()).toContain('The value must be less than or equal to 2');
  });

  it('shows max length errors with the configured length', async () => {
    const wrapper = mount(TestMaxLengthField, {
      global: {
        plugins: [createPinia(), i18n, [veeValidatePlugin, { i18n }]],
      },
    });

    await wrapper.find('input[name="description"]').setValue('too long');
    await flushPromises();

    expect(wrapper.text()).toContain('The field may not be greater than 5 characters');
  });

  it('does not add number limit rules for empty number rule values', async () => {
    const wrapper = mount(TestNumberField, {
      props: {
        optimalValue: '',
      },
      global: {
        components: { CTextField: TestField },
        plugins: [createPinia(), i18n, [veeValidatePlugin, { i18n }]],
      },
    });

    await wrapper.find('input[name="minimal"]').setValue('3');
    await flushPromises();

    expect(wrapper.text()).not.toContain('less than or equal');
    expect(wrapper.text()).not.toContain('NaN');
  });

  it('clamps number fields to non-negative minimum values', async () => {
    const wrapper = mount(TestNumberField, {
      global: {
        components: { CTextField: TestField },
        plugins: [createPinia(), i18n, [veeValidatePlugin, { i18n }]],
      },
    });

    await wrapper.find('input[name="minimal"]').setValue('-3');
    await flushPromises();

    expect(wrapper.vm.minimal).toBe(0);
  });
});
