import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import Title from '@/components/sidebars/form/fields/title.vue';

const stubs = {
  'widget-settings-item': true,
  'v-text-field': createInputStub('v-text-field'),
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('title', () => {
  const factory = generateShallowRenderer(Title, { stubs });

  it('Value changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = Faker.datatype.string();

    selectTextField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
