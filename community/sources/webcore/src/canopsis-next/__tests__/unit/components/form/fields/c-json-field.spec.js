import { createButtonStub } from '@unit/stubs/button';
import { createTextareaInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { stringifyJson } from '@/helpers/json';

import CJsonField from '@/components/forms/fields/compat/c-json-field.vue';

const stubs = {
  'c-textarea': createTextareaInputStub('v-textarea'),
  'c-payload-textarea-field': createTextareaInputStub('v-textarea'),
  'v-btn': createButtonStub('v-btn'),
};

const selectTextarea = wrapper => wrapper.find('.v-textarea textarea');
const selectParseButton = wrapper => wrapper.findAll('.v-btn').at(0);
const selectResetButton = wrapper => wrapper.findAll('.v-btn').at(1);

describe('c-json-field', () => {
  const name = 'jsonField';
  const defaultValue = '{}';
  const validJsonValue = { key: 'value' };
  const validJsonStringValue = stringifyJson({ newKey: 'newValue' });
  const validPayloadJsonValue = `{
    "extra_vars": {
        "version": 15,
        "state": {{ .Alarm.Value.StateValue }},
        "entity_id": {{ .Entity.ID }}
    }
}`;

  const invalidJsonStringValue = 'asd';
  const invalidPayloadJsonValue = `{
    "extra_vars": {
        "version": 15,
        "state": {{ .Alarm.Value.StateValue },
        "entity_id": {{ .Entity.ID }
    }
}`;

  const factory = generateShallowRenderer(CJsonField, { stubs });

  test('Object value as prop', () => {
    const wrapper = factory({
      propsData: { value: validJsonValue },
    });

    expect(selectTextarea(wrapper).element.value).toBe(stringifyJson(validJsonValue));
  });

  test('Input event after set value', async () => {
    const newValue = '"asd"';
    const wrapper = factory({
      propsData: { value: validJsonStringValue },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(newValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper).toEmitInput(newValue);
  });

  test('Payload json value as prop with variables', () => {
    const wrapper = factory({
      propsData: {
        value: validPayloadJsonValue,
        variables: true,
      },
    });

    expect(selectTextarea(wrapper).element.value).toBe(validPayloadJsonValue);
  });

  test('Payload json value as prop without variables', () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    const wrapper = factory({
      propsData: {
        value: validPayloadJsonValue,
      },
    });

    expect(selectTextarea(wrapper).element.value).toBe(defaultValue);
    expect(consoleErrorSpy).toHaveBeenCalled();

    consoleErrorSpy.mockRestore();
  });

  test('v-validate works correctly on valid json', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(validJsonStringValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBeFalsy();
  });

  test('v-validate works correctly on invalid json and blur event', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBe('Invalid JSON');
  });

  test('v-validate works correctly on invalid json and blur event (validate on `button`)', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
        validateOn: 'button',
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBeFalsy();
  });

  test('v-validate works correctly on invalid json and parse click', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
        validateOn: 'button',
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(invalidJsonStringValue);
    await selectParseButton(wrapper).trigger('click');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBe('Invalid JSON');
  });

  test('v-validate reset by previous value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonStringValue,
      },
    });

    const textarea = selectTextarea(wrapper);
    const originalValue = textarea.element.value;

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBe('Invalid JSON');

    await textarea.setValue(originalValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBeFalsy();
  });

  test('v-validate reset by another value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonStringValue,
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBe('Invalid JSON');

    await textarea.setValue(stringifyJson(validJsonValue));
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBeFalsy();
  });

  test('v-validate works correctly on valid payload json value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(validPayloadJsonValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBeFalsy();
  });

  test('v-validate works correctly on invalid payload json value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(invalidPayloadJsonValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(wrapper.vm.errorMessage).toBe('Invalid JSON');
  });

  test('Set new value by prop', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
      },
    });

    await wrapper.setProps({ value: validJsonStringValue });
    await flushPromises();

    expect(selectTextarea(wrapper).element.value).toBe(validJsonStringValue);
  });

  test('Set new value by prop (same value)', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
        value: validJsonValue,
      },
    });

    await wrapper.setProps({ value: { ...validJsonValue } });
    await flushPromises();

    expect(selectTextarea(wrapper).element.value).toBe(stringifyJson(validJsonValue));
  });

  test('Set the same values', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonStringValue,
        variables: true,
      },
    });

    const textarea = selectTextarea(wrapper);
    const originalValue = textarea.element.value;

    await textarea.setValue(originalValue);
    await textarea.trigger('blur');
    await flushPromises();

    expect(textarea.element.value).toBe(originalValue);
  });

  test('Check reset click', async () => {
    const wrapper = factory({
      propsData: {
        name,
        validateOn: 'button',
      },
    });

    const textarea = selectTextarea(wrapper);

    await textarea.setValue(validJsonStringValue);
    await selectResetButton(wrapper).trigger('click');
    await flushPromises();

    expect(textarea.element.value).toBe(defaultValue);
    expect(wrapper.vm.errorMessage).toBeFalsy();
  });
});
