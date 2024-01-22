import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import FieldMap from '@/components/sidebars/map/form/fields/map.vue';

const stubs = {
  'c-map-field': true,
};

const selectMapField = wrapper => wrapper.find('c-map-field-stub');

describe('field-map', () => {
  const factory = generateShallowRenderer(FieldMap, {
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
  });
  const snapshotFactory = generateRenderer(FieldMap, {
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
  });

  it('Info popup setting modal opened after trigger create button', () => {
    const wrapper = factory({
      propsData: {
        value: 'value',
      },
    });

    const mapField = selectMapField(wrapper);

    const newMap = Faker.datatype.string();

    mapField.triggerCustomEvent('input', newMap);

    expect(wrapper).toEmit('input', newMap);
  });

  it('Renders `field-map` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `field-map` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Value',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
