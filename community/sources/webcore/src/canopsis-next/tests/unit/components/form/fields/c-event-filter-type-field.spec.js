import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import CEventFilterTypeField from '@/components/forms/fields/c-event-filter-type-field.vue';
import { EVENT_FILTER_TYPES } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CEventFilterTypeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CEventFilterTypeField, {
  localVue,

  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-event-filter-type-field', () => {
  test('Value changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.setValue(EVENT_FILTER_TYPES.drop);

    expect(wrapper).toEmit('input', EVENT_FILTER_TYPES.drop);
  });

  test('Renders `c-event-filter-type-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Default value',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-event-filter-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-event-filter-type-field` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '',
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});
