import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

const stubs = {
  'variables-list': true,
};

const selectMenuNode = wrapper => wrapper.vm.$children[0];
const selectVariablesList = wrapper => wrapper.find('variables-list-stub');

describe('variables-menu', () => {
  const factory = generateShallowRenderer(VariablesMenu, { stubs });

  test('Input event emitted after trigger variables', () => {
    const wrapper = factory();

    const variablesList = selectVariablesList(wrapper);

    const value = Faker.datatype.string();

    variablesList.triggerCustomEvent('input', value);

    expect(wrapper).toEmitInput(value);
  });

  test('Close event emitted after trigger menu', () => {
    const wrapper = factory();

    const menuNode = selectMenuNode(wrapper);

    menuNode.$emit('input');

    expect(wrapper).toHaveBeenEmit('close');
  });
});
