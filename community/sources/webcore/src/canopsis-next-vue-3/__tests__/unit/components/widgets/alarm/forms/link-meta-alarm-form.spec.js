import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import LinkMetaAlarmForm from '@/components/widgets/alarm/forms/link-meta-alarm-form.vue';

const stubs = {
  'c-meta-alarm-field': {
    ...createInputStub('c-meta-alarm-field'),
    name: 'CMetaAlarmField',
    props: {
      value: [String, Object],
      modelValue: [String, Object],
      returnObject: Boolean,
    },
  },
  'c-text-field': createInputStub('c-text-field'),
  'meta-alarm-rule-tags-form': true,
  'meta-alarm-rule-infos-form': true,
  'c-enabled-field': true,
  'c-name-field': createInputStub('c-name-field'),
  'c-help-icon': true,
};

const selectTextField = wrapper => wrapper.find('.c-name-field');
const selectMetaAlarmField = wrapper => wrapper.find('.c-meta-alarm-field');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');

describe('link-meta-alarm-form', () => {
  const factory = generateShallowRenderer(LinkMetaAlarmForm, { stubs });

  test('Meta alarm field returns the selected meta alarm object', () => {
    const wrapper = factory({
      propsData: {
        form: {
          metaAlarm: null,
          comment: '',
          auto_resolve: false,
        },
      },
    });

    expect(wrapper.findComponent({ name: 'CMetaAlarmField' }).props('returnObject')).toBe(true);
  });

  test('Meta alarm changed after trigger text field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      comment: Faker.datatype.string(),
      auto_resolve: false,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const metaAlarm = Faker.datatype.string();

    const comboboxField = selectMetaAlarmField(wrapper);

    comboboxField.setValue(metaAlarm);

    expect(wrapper).toEmit('update:form', { ...form, metaAlarm });
  });

  test('Comment changed after trigger description field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      comment: Faker.datatype.string(),
      auto_resolve: false,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const comment = Faker.datatype.string();
    selectTextField(wrapper).setValue(comment);

    expect(wrapper).toEmit('update:form', { ...form, comment });
  });

  test('Auto resolve changed after trigger enabled field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      comment: Faker.datatype.string(),
      auto_resolve: false,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const autoResolve = true;

    const enabledField = selectEnabledField(wrapper);

    enabledField.triggerCustomEvent('update:modelValue', autoResolve);

    expect(wrapper).toEmit('update:form', { ...form, auto_resolve: autoResolve });
  });
});
