import { generateShallowRenderer } from '@unit/utils/vue';

import { EXTERNAL_DATA_ON_MISMATCH_TYPES } from '@/constants';

import ExternalDataTableForm from '@/components/forms/external-data/external-data-table-form.vue';

const stubs = {
  'external-data-on-mismatch-field': true,
  'external-data-table-table-field': true,
  'external-data-table-condition-form': true,
};

const selectOnMismatchField = wrapper => wrapper.find('external-data-on-mismatch-field-stub');

describe('external-data-table-form', () => {
  const factory = generateShallowRenderer(ExternalDataTableForm, { stubs });
  const form = {
    optional: false,
    table: 'table-id',
    conditions: [{ key: 'condition-key' }],
    on_mismatch: EXTERNAL_DATA_ON_MISMATCH_TYPES.throwError,
  };

  it('shows mismatch behavior for required event-filter external data', () => {
    const wrapper = factory({
      propsData: {
        form,
        name: 'external_data.item-key',
        withOnMismatch: true,
      },
    });

    expect(selectOnMismatchField(wrapper).exists()).toBe(true);
  });

  it('hides mismatch behavior when external data is optional', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          optional: true,
        },
        name: 'external_data.item-key',
        withOnMismatch: true,
      },
    });

    expect(selectOnMismatchField(wrapper).exists()).toBe(false);
  });

  it('hides mismatch behavior outside event filters', () => {
    const wrapper = factory({
      propsData: {
        form,
        name: 'external_data.item-key',
      },
    });

    expect(selectOnMismatchField(wrapper).exists()).toBe(false);
  });

  it('updates mismatch behavior', () => {
    const wrapper = factory({
      propsData: {
        form,
        name: 'external_data.item-key',
        withOnMismatch: true,
      },
    });

    selectOnMismatchField(wrapper).triggerCustomEvent(
      'update:modelValue',
      EXTERNAL_DATA_ON_MISMATCH_TYPES.skipActions,
    );

    expect(wrapper.emitted('update:form')[0][0]).toEqual({
      ...form,
      on_mismatch: EXTERNAL_DATA_ON_MISMATCH_TYPES.skipActions,
    });
  });
});
