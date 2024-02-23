import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { EVENT_FILTER_TYPES } from '@/constants';

import CEventFilterTypeField from '@/components/forms/fields/c-event-filter-type-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-event-filter-type-field', () => {
  const factory = generateShallowRenderer(CEventFilterTypeField, { stubs });
  const snapshotFactory = generateRenderer(CEventFilterTypeField, {
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Value changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.setValue(EVENT_FILTER_TYPES.drop);

    expect(wrapper).toEmitInput(EVENT_FILTER_TYPES.drop);
  });

  test('Renders `c-event-filter-type-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Default value',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-event-filter-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-event-filter-type-field` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '',
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});
