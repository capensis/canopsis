import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createTextareaInputStub } from '@unit/stubs/input';
import { createButtonStub } from '@unit/stubs/button';

import { stringifyJson } from '@/helpers/json';

import CJsonField from '@/components/forms/fields/c-json-field.vue';

const stubs = {
  'v-textarea': createTextareaInputStub('v-textarea'),
  'v-btn': createButtonStub('v-btn'),
  'c-help-icon': true,
};

const snapshotStubs = {
  'c-help-icon': true,
};

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

  const factory = generateShallowRenderer(CJsonField, {
    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });
  const snapshotFactory = generateRenderer(CJsonField, {
    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Object value as prop', () => {
    const wrapper = factory({
      propsData: { value: validJsonValue },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    expect(textarea.element.value).toBe(stringifyJson(validJsonValue));
  });

  it('Input event after set value', async () => {
    const newValue = '"asd"';
    const wrapper = factory({
      propsData: { value: validJsonStringValue },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(newValue);
    await textarea.trigger('blur');

    expect(wrapper).toEmit('input', newValue);
  });

  it('Payload json value as prop with variables', () => {
    const wrapper = factory({
      propsData: {
        value: validPayloadJsonValue,
        variables: true,
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    expect(textarea.element.value).toBe(validPayloadJsonValue);
  });

  it('Payload json value as prop without variables', () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();

    const popupErrorFn = jest.fn();
    const wrapper = factory({
      propsData: {
        value: validPayloadJsonValue,
      },
      mocks: {
        $popups: {
          error: popupErrorFn,
        },
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    expect(textarea.element.value).toBe(defaultValue);
    expect(popupErrorFn).toBeCalledWith({ text: 'Something went wrong...' });

    consoleErrorSpy.mockClear();
  });

  it('v-validate works correctly on valid json', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(validJsonStringValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeFalsy();
  });

  it('v-validate works correctly on invalid json and blur event', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeTruthy();
  });

  it('v-validate works correctly on invalid json and blur event (validate on `button`)', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
        validateOn: 'button',
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeFalsy();
  });

  it('v-validate works correctly on invalid json and parse click', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonValue,
        validateOn: 'button',
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');
    const button = wrapper.find('.v-btn');

    await textarea.setValue(invalidJsonStringValue);
    await button.trigger('click');

    expect(validator.errors.has(name)).toBeTruthy();
  });

  it('v-validate reset by previous value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonStringValue,
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');
    const originalValue = textarea.element.value;

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeTruthy();

    await textarea.setValue(originalValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeFalsy();
  });

  it('v-validate reset by another value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonStringValue,
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(invalidJsonStringValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeTruthy();

    await textarea.setValue(stringifyJson(validJsonValue));
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeFalsy();
  });

  it('v-validate works correctly on valid payload json value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(validPayloadJsonValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeFalsy();
  });

  it('v-validate works correctly on invalid payload json value', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(invalidPayloadJsonValue);
    await textarea.trigger('blur');

    expect(validator.errors.has(name)).toBeTruthy();
  });

  it('Set new value by prop', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    await wrapper.setProps({ value: validJsonStringValue });

    expect(textarea.element.value).toBe(validJsonStringValue);
  });

  it('Set new value by prop (same value)', async () => {
    const wrapper = factory({
      propsData: {
        name,
        variables: true,
        value: validJsonValue,
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');

    await wrapper.setProps({ value: { ...validJsonValue } });

    expect(textarea.element.value).toBe(stringifyJson(validJsonValue));
  });

  it('Set the same values', async () => {
    const wrapper = factory({
      propsData: {
        name,
        value: validJsonStringValue,
        variables: true,
      },
    });

    const textarea = wrapper.find('.v-textarea textarea');
    const originalValue = textarea.element.value;

    await textarea.setValue(originalValue);
    await textarea.trigger('blur');

    expect(textarea.element.value).toBe(originalValue);
  });

  it('Check reset click', async () => {
    const wrapper = factory({
      propsData: {
        name,
        validateOn: 'button',
      },
    });

    const validator = wrapper.getValidator();
    const textarea = wrapper.find('.v-textarea textarea');
    const button = wrapper.findAll('.v-btn').at(1);

    await textarea.setValue(validJsonStringValue);
    await button.trigger('click');

    expect(textarea.element.value).toBe(defaultValue);
    expect(validator.errors.has(name)).toBeFalsy();
  });

  it('Renders `c-json-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        validateOn: 'button',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        validateOn: 'button',
        readonly: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-json-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
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

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-json-field` with default props and changed value correctly', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        validateOn: 'button',
      },
    });
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(validJsonStringValue);

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-json-field` with default props and incorrect value correctly', async () => {
    const newValue = '{ key: value }';
    const wrapper = snapshotFactory();
    const textarea = wrapper.find('.v-textarea textarea');

    await textarea.setValue(newValue);
    await textarea.trigger('blur');

    expect(wrapper).toMatchSnapshot();
  });
});
