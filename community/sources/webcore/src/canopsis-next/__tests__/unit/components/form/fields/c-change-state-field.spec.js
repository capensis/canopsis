import { createNumberInputStub, createTextareaInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { ALARM_STATES } from '@/constants';

import CChangeStateField from '@/components/forms/fields/c-change-state-field.vue';

vi.mock('@/hooks/store/modules/info', () => ({
  useInfo: () => ({
    allowChangeSeverityToInfo: { value: false },
  }),
}));

const stubs = {
  'c-form-block-row': true,
  'state-criticity-field': createNumberInputStub('state-criticity-field'),
  'c-payload-textarea-field': createTextareaInputStub('c-payload-textarea-field'),
};

const selectStateField = wrapper => wrapper.find('input.state-criticity-field');
const selectPayloadTextareaField = wrapper => wrapper.find('.c-payload-textarea-field');

describe('c-change-state-field', () => {
  const factory = generateShallowRenderer(CChangeStateField, { stubs });

  it('State changed after trigger the state field', () => {
    const initialValue = {
      state: ALARM_STATES.major,
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: {
        modelValue: initialValue,
      },
    });

    selectStateField(wrapper).triggerCustomEvent('update:modelValue', ALARM_STATES.critical);

    expect(wrapper).toEmitInput({
      ...initialValue,
      state: ALARM_STATES.critical,
    });
  });

  it('Output changed after trigger the textarea field', () => {
    const initialValue = {
      state: ALARM_STATES.major,
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: {
        modelValue: initialValue,
      },
    });

    const newOutput = Faker.datatype.string();

    selectPayloadTextareaField(wrapper).triggerCustomEvent('update:value', newOutput);

    expect(wrapper).toEmitInput({
      ...initialValue,
      output: newOutput,
    });
  });
});
