import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import ChartTitle from '@/components/sidebars/settings/fields/chart/chart-title.vue';

const stubs = {
  'widget-settings-item': true,
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('chart-title', () => {
  const factory = generateShallowRenderer(ChartTitle, {
    stubs,
  });

  const snapshotFactory = generateRenderer(ChartTitle, {
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

  it('Value changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = Faker.datatype.string();

    selectTextField(wrapper).setValue(newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  it('Renders `chart-title` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `chart-title` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        name: 'custom_name',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
