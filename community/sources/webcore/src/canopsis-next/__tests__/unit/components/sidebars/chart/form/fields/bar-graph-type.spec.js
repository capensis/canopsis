import { createCheckboxInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import BarGraphType from '@/components/sidebars/chart/form/fields/bar-graph-type.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createCheckboxInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('bar-graph-type', () => {
  const factory = generateShallowRenderer(BarGraphType, { stubs });

  test('Value changed after trigger radio group field', () => {
    const value = Faker.datatype.boolean();
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectRadioGroup(wrapper).setChecked(!value);

    expect(wrapper).toEmitInput(!value);
  });
});
