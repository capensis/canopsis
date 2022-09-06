import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import VariablesList from '@/components/common/text-editor/variables-list.vue';

const localVue = createVueInstance();

const stubs = {
  'variables-list': true,
};

const factory = (options = {}) => shallowMount(VariablesList, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(VariablesList, {
  localVue,
  attachTo: document.body,
  stubs,

  ...options,
});

const selectVariableTiles = wrapper => wrapper.findAll('v-list-tile-stub');
const selectVariableTileByIndex = (wrapper, index) => selectVariableTiles(wrapper).at(index);
const selectMenu = wrapper => wrapper.find('v-menu-stub');
const selectVariablesList = wrapper => wrapper.find('variables-list-stub');

describe('variables-list', () => {
  test('Variable selected after click on tile', () => {
    const value = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        variables: [
          { value },
        ],
      },
    });

    const variableTile = selectVariableTileByIndex(wrapper, 0);

    variableTile.vm.$emit('click');

    expect(wrapper).toEmit('input', value);
  });

  test('Submenu opened when mouseover tile', async () => {
    const zIndex = Faker.datatype.number();
    const subVariables = [
      {
        value: 'child',
      },
    ];
    const wrapper = factory({
      propsData: {
        zIndex,
        value: 'parent.child',
        variables: [
          {
            value: 'parent',
            variables: subVariables,
          },
        ],
      },
    });

    const variableTile = selectVariableTileByIndex(wrapper, 0);

    jest.spyOn(variableTile.element, 'getBoundingClientRect').mockReturnValue({
      top: 101,
      left: 112,
      width: 88,
    });
    await variableTile.vm.$emit('mouseenter', { target: variableTile.element });

    const menu = selectMenu(wrapper);
    expect(menu.vm.positionX).toEqual(200);
    expect(menu.vm.positionY).toEqual(101);
    expect(menu.vm.zIndex).toEqual(zIndex);

    const variablesList = selectVariablesList(wrapper);
    expect(variablesList.vm.$attrs.variables).toEqual(subVariables);
    expect(variablesList.vm.$attrs['z-index']).toBe(zIndex + 1);
    expect(variablesList.vm.$attrs.value).toBe(subVariables[0].value);
  });

  test('Submenu closed after mouseover on other tile', async () => {
    const wrapper = factory({
      propsData: {
        variables: [
          {
            value: 'first',
            variables: [],
          },
          {
            value: 'second',
          },
        ],
      },
    });

    const firstVariableTile = selectVariableTileByIndex(wrapper, 0);

    await firstVariableTile.vm.$emit('mouseenter', {
      target: firstVariableTile.element,
    });

    expect(selectMenu(wrapper).element).toBeTruthy();

    const secondVariableTile = selectVariableTileByIndex(wrapper, 1);
    await secondVariableTile.vm.$emit('mouseenter', {
      target: secondVariableTile.element,
    });
  });

  test('Sub variable selected after trigger input on variables list', async () => {
    const parentValue = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        variables: [
          {
            value: parentValue,
            variables: [],
          },
        ],
      },
    });

    const firstVariableTile = selectVariableTileByIndex(wrapper, 0);

    await firstVariableTile.vm.$emit('mouseenter', {
      target: firstVariableTile.element,
    });

    const variablesList = selectVariablesList(wrapper);

    const value = Faker.datatype.string();

    variablesList.vm.$emit('input', value);

    expect(wrapper).toEmit('input', `${parentValue}.${value}`);
  });

  test('Renders `variables-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `variables-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'entity._id',
        zIndex: 2,
        variables: [
          {
            value: 'entity',
            text: 'Entity',
            variables: [
              {
                value: 'id',
                text: 'Id',
              },
            ],
          },
        ],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
