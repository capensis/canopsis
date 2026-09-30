import { generateShallowRenderer } from '@unit/utils/vue';

import CCardIteratorForm from '@/components/forms/fields/card-iterator/c-card-iterator-form.vue';

const stubs = {
  'c-layout': true,
  'c-alert': true,
  'c-card-iterator-field': true,
  'c-btn-with-error': true,
};

const selectIteratorField = wrapper => wrapper.find('c-card-iterator-field-stub');
const selectAddButton = wrapper => wrapper.find('c-btn-with-error-stub');

describe('c-card-iterator-form', () => {
  const factory = generateShallowRenderer(CCardIteratorForm, { stubs });

  test('Add button and iterator are disabled in read-only mode', () => {
    const wrapper = factory({
      props: {
        addButtonLabel: 'Add item',
        disabled: true,
      },
    });

    expect(selectIteratorField(wrapper).attributes('disabled')).toBe('true');
    expect(selectAddButton(wrapper).attributes('disabled')).toBe('true');
  });
});
