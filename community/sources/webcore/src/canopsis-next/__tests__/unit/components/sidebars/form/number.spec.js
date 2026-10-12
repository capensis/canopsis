import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { NUMBERS_CHART_DEFAULT_FONT_SIZE } from '@/constants';

import FontSize from '@/components/sidebars/chart/form/fields/font-size.vue';

const stubs = {
  'widget-settings-item': true,
  'c-radio-group': true,
  'c-number-field': true,
};

const selectNumberField = wrapper => wrapper.find('c-number-field-stub');

describe('field-number', () => {
  const factory = generateShallowRenderer(FontSize, { stubs });

  test('Value changed after trigger number field', () => {
    const wrapper = factory({
      propsData: {
        value: NUMBERS_CHART_DEFAULT_FONT_SIZE,
      },
    });

    const newValue = Faker.datatype.number();

    selectNumberField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
