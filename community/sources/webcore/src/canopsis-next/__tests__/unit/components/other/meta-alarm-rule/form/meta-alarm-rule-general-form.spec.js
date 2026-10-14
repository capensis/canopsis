import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleGeneralForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue';

const stubs = {
  'c-form-block': true,
  'c-form-block-row': true,
  'c-id-field': {
    name: 'CIdField',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'c-name-field': {
    name: 'CNameField',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'c-payload-textarea-field': {
    name: 'CPayloadTextareaField',
    emits: ['update:value'],
    template: '<div />',
  },
  'c-payload-text-field': true,
  'c-enabled-field': {
    name: 'CEnabledField',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'meta-alarm-rule-tags-form': {
    name: 'MetaAlarmRuleTagsForm',
    emits: ['update:form'],
    template: '<div />',
  },
  'meta-alarm-rule-infos-form': {
    name: 'MetaAlarmRuleInfosForm',
    emits: ['update:infos'],
    template: '<div />',
  },
  'meta-alarm-rule-type-field': true,
};

const selectIdField = wrapper => wrapper.findComponent({ name: 'CIdField' });
const selectPayloadTextareaField = wrapper =>
  wrapper.findComponent({ name: 'CPayloadTextareaField' });
const selectNameField = wrapper => wrapper.findComponent({ name: 'CNameField' });
const selectEnabledField = wrapper => wrapper.findComponent({ name: 'CEnabledField' });
const selectTagsForm = wrapper => wrapper.findComponent({ name: 'MetaAlarmRuleTagsForm' });
const selectInfosForm = wrapper => wrapper.findComponent({ name: 'MetaAlarmRuleInfosForm' });

describe('meta-alarm-rule-general-form', () => {
  const form = metaAlarmRuleToForm();

  const factory = generateShallowRenderer(MetaAlarmRuleGeneralForm, { stubs });

  test('ID changed after trigger id field', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const idField = selectIdField(wrapper);

    const newId = Faker.datatype.string();

    idField.vm.$emit('update:modelValue', newId);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('_id')).toBe(newId);
  });

  test('Output template changed after trigger description field', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const payloadTextareaField = selectPayloadTextareaField(wrapper);

    const outputTemplate = Faker.datatype.string();

    payloadTextareaField.vm.$emit('update:value', outputTemplate);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('output_template')).toBe(outputTemplate);
  });

  test('Name changed after trigger text field', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const nameField = selectNameField(wrapper);

    const name = Faker.datatype.string();

    nameField.vm.$emit('update:modelValue', name);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('name')).toBe(name);
  });

  test('Enabled changed after trigger enabled field', async () => {
    const autoResolve = Faker.datatype.boolean();
    const wrapper = factory({
      props: {
        form: {
          ...form,
          auto_resolve: autoResolve,
        },
      },
    });

    const enabledField = selectEnabledField(wrapper);

    const newAutoResolve = !autoResolve;

    enabledField.vm.$emit('update:modelValue', newAutoResolve);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('auto_resolve')).toBe(newAutoResolve);
  });

  test('Tags changed after trigger tags form', async () => {
    const copyFromChildren = Faker.datatype.boolean();
    const wrapper = factory({
      props: {
        form: {
          ...form,
          tags: { copy_from_children: copyFromChildren },
        },
      },
    });

    const newTags = {
      copy_from_children: !copyFromChildren,
    };

    selectTagsForm(wrapper).vm.$emit('update:form', newTags);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('tags')).toEqual(newTags);
  });

  test('Infos changed after trigger infos form', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const newInfos = [{ key: Faker.datatype.string() }];

    selectInfosForm(wrapper).vm.$emit('update:infos', newInfos);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('infos')).toEqual(newInfos);
  });
});
