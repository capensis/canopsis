import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

const stubs = {
  'variables-list': true,
};

const selectMenuNode = wrapper => wrapper.vm.$children[0];
const selectVariablesList = wrapper => wrapper.find('variables-list-stub');

describe('variables-menu', () => {
  const factory = generateShallowRenderer(VariablesMenu, { stubs });
  const snapshotFactory = generateRenderer(VariablesMenu, {
    attachTo: document.body,
    stubs,
  });

  test('Input event emitted after trigger variables', () => {
    const wrapper = factory();

    const variablesList = selectVariablesList(wrapper);

    const value = Faker.datatype.string();

    variablesList.triggerCustomEvent('input', value);

    expect(wrapper).toEmit('input', value);
  });

  test('Close event emitted after trigger menu', () => {
    const wrapper = factory();

    const menuNode = selectMenuNode(wrapper);

    menuNode.$emit('input');

    expect(wrapper).toEmit('close');
  });

  test('Renders `variables-menu` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `variables-menu` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        visible: true,
        value: 'entity._id',
        positionX: 2,
        positionY: 3,
        variables: [
          {
            value: 'entity._id',
            text: 'Variable',
          },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
