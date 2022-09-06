import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

const localVue = createVueInstance();

const stubs = {
  'variables-list': true,
};

const factory = (options = {}) => shallowMount(VariablesMenu, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(VariablesMenu, {
  localVue,
  attachTo: document.body,
  stubs,

  ...options,
});

const selectMenuNode = wrapper => wrapper.vm.$children[0];
const selectVariablesList = wrapper => wrapper.find('variables-list-stub');

describe('variables-menu', () => {
  test('Input event emitted after trigger variables', () => {
    const wrapper = factory();

    const variablesList = selectVariablesList(wrapper);

    const value = Faker.datatype.string();

    variablesList.vm.$emit('input', value);

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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
