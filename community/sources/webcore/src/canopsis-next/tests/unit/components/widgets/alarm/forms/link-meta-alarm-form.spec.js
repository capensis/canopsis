import Faker from 'faker';
import { Validator } from 'vee-validate';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import LinkMetaAlarmForm from '@/components/widgets/alarm/forms/link-meta-alarm-form.vue';

const stubs = {
  'c-meta-alarm-field': createInputStub('c-meta-alarm-field'),
  'v-text-field': createInputStub('v-text-field'),
  'meta-alarm-rule-tags-form': true,
  'meta-alarm-rule-infos-form': true,
  'c-enabled-field': true,
  'c-name-field': createInputStub('c-name-field'),
  'c-help-icon': true,
};

const snapshotStubs = {
  'c-meta-alarm-field': true,
  'meta-alarm-rule-tags-form': true,
  'meta-alarm-rule-infos-form': true,
  'c-enabled-field': true,
  'c-name-field': true,
  'c-help-icon': true,
};

const selectTextField = wrapper => wrapper.find('.c-name-field');
const selectMetaAlarmField = wrapper => wrapper.find('.c-meta-alarm-field');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');

describe('link-meta-alarm-form', () => {
  const factory = generateShallowRenderer(LinkMetaAlarmForm, {
    stubs,
    provide: {
      validator: new Validator(),
    },
  });
  const snapshotFactory = generateRenderer(LinkMetaAlarmForm, {
    stubs: snapshotStubs,
    provide: {
      validator: new Validator(),
    },
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

    expect(wrapper).toEmitInput({ ...form, metaAlarm });
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

    expect(wrapper).toEmitInput({ ...form, comment });
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

    enabledField.triggerCustomEvent('input', autoResolve);

    expect(wrapper).toEmitInput({ ...form, auto_resolve: autoResolve });
  });

  test('Renders `link-meta-alarm-form` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          metaAlarm: 'metaAlarm',
          comment: 'comment',
          auto_resolve: true,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
