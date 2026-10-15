import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FieldSlider from '@/components/sidebars/form/fields/slider.vue';

const stubs = {
  'widget-settings-item': true,
  'v-slider': createNumberInputStub('v-slider'),
};

const selectSliderField = wrapper => wrapper.find('.v-slider');

describe('field-slider', () => {
  const factory = generateShallowRenderer(FieldSlider, { stubs });

  test('Value changed after trigger number field', () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        title: '',
        value,
      },
    });

    const newValue = Faker.datatype.number({ min: value });

    selectSliderField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
