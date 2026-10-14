import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FieldGridSize from '@/components/sidebars/form/fields/grid-size.vue';

const stubs = {
  'widget-settings-item': true,
  'c-column-size-field': true,
};

const selectColumnSizeField = wrapper => wrapper.find('c-column-size-field-stub');

describe('field-grid-size', () => {
  const factory = generateShallowRenderer(FieldGridSize, { stubs });

  it('Column size changed after trigger field', () => {
    const wrapper = factory({
      propsData: {
        title: 'Title',
      },
    });

    const columnSizeField = selectColumnSizeField(wrapper);

    const newSize = Faker.datatype.number({ min: 1, max: 12 });

    columnSizeField.triggerCustomEvent('input', newSize);

    expect(wrapper).toEmitInput(newSize);
  });
});
