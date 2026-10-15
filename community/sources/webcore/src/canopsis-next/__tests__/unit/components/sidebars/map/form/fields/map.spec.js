import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

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
          register: vi.fn(),
          unregister: vi.fn(),
        },
        listClick: vi.fn(),
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

    expect(wrapper).toEmitInput(newMap);
  });
});
