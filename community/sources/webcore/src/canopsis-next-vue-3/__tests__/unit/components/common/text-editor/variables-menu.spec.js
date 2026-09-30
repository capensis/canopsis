import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

const stubs = {
  'v-menu': {
    template:
      '<div><slot /><button class="menu-close" @click="$emit(\'update:modelValue\', false)" /></div>',
  },
  'variables-list': {
    template: '<button class="variables-list" @click="$emit(\'update:value\', value)" />',
    props: ['value'],
  },
};

const selectVariablesList = wrapper => wrapper.find('variables-list-stub');

describe('variables-menu', () => {
  const factory = generateShallowRenderer(VariablesMenu, { stubs });

  test('Input event emitted after trigger variables', async () => {
    const wrapper = factory();

    const variablesList = selectVariablesList(wrapper);

    const value = Faker.datatype.string();

    wrapper.vm.$emit('update:value', value);

    expect(wrapper.emitted('update:value')).toEqual([[value]]);
  });

  test('Close event emitted after trigger menu', async () => {
    const wrapper = factory();

    wrapper.vm.$emit('close');

    expect(wrapper.emitted('close')).toHaveLength(1);
  });
});
