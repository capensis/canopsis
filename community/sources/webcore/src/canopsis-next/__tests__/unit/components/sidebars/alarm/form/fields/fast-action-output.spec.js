import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FastActionOutput from '@/components/sidebars/alarm/form/fields/fast-action-output.vue';

const stubs = {
  'widget-settings-item': true,
  'c-enabled-field': true,
  'c-text-field': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectTextField = wrapper => wrapper.find('c-text-field-stub');

describe('fast-action-output', () => {
  const factory = generateShallowRenderer(FastActionOutput, { stubs });

  it('Enabled changed after trigger the enabled field', () => {
    const value = {
      enabled: true,
      value: '',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', false);

    expect(wrapper).toEmitInput({
      ...value,
      enabled: false,
    });
  });

  it('Value changed after trigger the text field', () => {
    const value = {
      enabled: true,
      value: '',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const newValue = Faker.datatype.string();

    selectTextField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({
      ...value,
      value: newValue,
    });
  });
});
