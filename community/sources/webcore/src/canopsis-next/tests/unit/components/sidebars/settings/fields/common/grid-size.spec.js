import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import FieldGridSize from '@/components/sidebars/settings/fields/common/grid-size.vue';

const localVue = createVueInstance();

const stubs = {
  'c-column-size-field': true,
};

const factory = (options = {}) => shallowMount(FieldGridSize, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      list: {
        register: jest.fn(),
        unregister: jest.fn(),
      },
      listClick: jest.fn(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(FieldGridSize, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      list: {
        register: jest.fn(),
        unregister: jest.fn(),
      },
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectColumnSizeField = wrapper => wrapper.find('c-column-size-field-stub');

describe('field-grid-size', () => {
  it('Column size changed after trigger field', () => {
    const wrapper = factory({
      propsData: {
        title: 'Title',
      },
    });

    const columnSizeField = selectColumnSizeField(wrapper);

    const newSize = Faker.datatype.number({ min: 1, max: 12 });

    columnSizeField.vm.$emit('input', newSize);

    expect(wrapper).toEmit('input', newSize);
  });

  it('Renders `field-grid-size` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `field-grid-size` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 12,
        title: 'Custom title',
        mobile: true,
        tablet: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
