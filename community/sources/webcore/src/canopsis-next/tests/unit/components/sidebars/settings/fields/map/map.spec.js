import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import FieldMap from '@/components/sidebars/settings/fields/map/map.vue';

const localVue = createVueInstance();

const stubs = {
  'c-map-field': true,
};

const factory = (options = {}) => shallowMount(FieldMap, {
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

const snapshotFactory = (options = {}) => mount(FieldMap, {
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

const selectMapField = wrapper => wrapper.find('c-map-field-stub');

describe('field-map', () => {
  it('Info popup setting modal opened after trigger create button', () => {
    const wrapper = factory({
      propsData: {
        value: 'value',
      },
    });

    const mapField = selectMapField(wrapper);

    const newMap = Faker.datatype.string();

    mapField.vm.$emit('input', newMap);

    expect(wrapper).toEmit('input', newMap);
  });

  it('Renders `field-map` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `field-map` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Value',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
