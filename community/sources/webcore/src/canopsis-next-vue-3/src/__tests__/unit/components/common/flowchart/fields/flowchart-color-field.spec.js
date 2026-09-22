import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FlowchartColorField from '@/components/common/flowchart/fields/flowchart-color-field.vue';

const selectCheckboxField = wrapper => wrapper.find('v-checkbox-stub');
describe('flowchart-color-field', () => {
  const factory = generateShallowRenderer(FlowchartColorField);

  const pickerStubs = {
    'c-color-picker-field': {
      template: '<button class="color-picker" @click="$emit(\'update:modelValue\', value)" />',
      props: ['value'],
    },
  };

  test('Value changed to first palette color after trigger select field with true', () => {
    const firstColor = Faker.internet.color();
    const wrapper = factory({
      stubs: pickerStubs,
      propsData: {
        palette: [firstColor],
      },
    });

    const checkboxField = selectCheckboxField(wrapper);

    checkboxField.triggerCustomEvent('change', true);

    expect(wrapper.emitted('update:value')).toEqual([[firstColor]]);
  });

  test('Value changed to transparent after trigger select field with false', () => {
    const wrapper = factory({
      stubs: pickerStubs,
      propsData: {
        color: Faker.internet.color(),
      },
    });

    const checkboxField = selectCheckboxField(wrapper);

    checkboxField.triggerCustomEvent('change', false);

    expect(wrapper.emitted('update:value')).toEqual([['transparent']]);
  });

  test('Color picker value is emitted after picker change', async () => {
    const color = Faker.internet.color();
    const palette = [Faker.internet.color()];
    const wrapper = factory({ propsData: { value: color, palette }, stubs: pickerStubs });
    const newColor = Faker.internet.color();
    wrapper.vm.updateModel(newColor);
    expect(wrapper.emitted('update:value')).toEqual([[newColor]]);
  });
});
