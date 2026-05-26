import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { IDLE_RULE_TYPES, TIME_UNITS } from '@/constants';

import { idleRuleToForm } from '@/helpers/entities/idle-rule/form';

import IdleRuleGeneralForm from '@/components/other/idle-rule/form/idle-rule-general-form.vue';

const stubs = {
  'c-form-block': true,
  'c-form-block-row': true,
  'c-name-field': true,
  'c-priority-field': true,
  'idle-rule-alarm-type-field': true,
  'c-description-field': true,
  'c-duration-field': true,
  'c-disable-during-periods-field': true,
  'c-action-type-field': true,
  'action-parameters-form': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectActionTypeField = wrapper => wrapper.find('c-action-type-field-stub');

describe('idle-rule-general-form', () => {
  const alarmForm = idleRuleToForm({
    type: IDLE_RULE_TYPES.alarm,
  });

  const entityForm = idleRuleToForm({
    type: IDLE_RULE_TYPES.entity,
  });

  const factory = generateShallowRenderer(IdleRuleGeneralForm, { stubs });
  const snapshotFactory = generateRenderer(IdleRuleGeneralForm, { stubs });

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form: alarmForm,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({ ...alarmForm, name: newName });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form: alarmForm,
      },
    });

    const newDescription = Faker.datatype.string();

    selectDescriptionField(wrapper).triggerCustomEvent('input', newDescription);

    expect(wrapper).toEmitInput({ ...alarmForm, description: newDescription });
  });

  test('Action type field is rendered for alarm type', () => {
    const wrapper = factory({
      propsData: {
        form: alarmForm,
        isEntityType: false,
      },
    });

    expect(selectActionTypeField(wrapper).exists()).toBe(true);
  });

  test('Action type field is hidden for entity type', () => {
    const wrapper = factory({
      propsData: {
        form: entityForm,
        isEntityType: true,
      },
    });

    expect(selectActionTypeField(wrapper).exists()).toBe(false);
  });

  test('Renders `idle-rule-general-form` with alarm type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          ...alarmForm,
          duration: {
            value: 1,
            unit: TIME_UNITS.hour,
          },
        },
        isEntityType: false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `idle-rule-general-form` with entity type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          ...entityForm,
          duration: {
            value: 1,
            unit: TIME_UNITS.hour,
          },
        },
        isEntityType: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
