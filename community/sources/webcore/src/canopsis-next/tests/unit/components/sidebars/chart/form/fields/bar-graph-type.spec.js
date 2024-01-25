import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';

import BarGraphType from '@/components/sidebars/chart/form/fields/bar-graph-type.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createCheckboxInputStub('v-radio-group'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('bar-graph-type', () => {
  const factory = generateShallowRenderer(BarGraphType, { stubs });
  const snapshotFactory = generateRenderer(BarGraphType, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

  test('Value changed after trigger radio group field', () => {
    const value = Faker.datatype.boolean();
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectRadioGroup(wrapper).setChecked(!value);

    expect(wrapper).toEmit('input', !value);
  });

  test('Renders `bar-graph-type` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `bar-graph-type` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: true,
        name: 'custom_name',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
