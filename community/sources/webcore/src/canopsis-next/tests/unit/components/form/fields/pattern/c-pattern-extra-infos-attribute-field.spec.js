import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createInputStub } from '@unit/stubs/input';

import CPatternExtraInfosAttributeField from '@/components/forms/fields/pattern/c-pattern-extra-infos-attribute-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
};

const factory = (options = {}) => shallowMount(CPatternExtraInfosAttributeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternExtraInfosAttributeField, {
  localVue,

  ...options,
});

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-pattern-extra-infos-attribute-field', () => {
  it('Field changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        value: {
          field: 'test.name',
        },
      },
    });
    const fieldSelect = selectTextField(wrapper);

    const newField = 'newField';

    fieldSelect.setValue(newField);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      field: newField,
    });
  });

  it('Field changed after trigger the text field without initial field', () => {
    const wrapper = factory({
      propsData: {
        value: {
          field: '',
        },
      },
    });
    const fieldSelect = selectTextField(wrapper);

    const newField = 'newField';

    fieldSelect.setValue(newField);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      field: newField,
    });
  });

  it('Renders `c-pattern-extra-infos-attribute-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          field: '',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-pattern-extra-infos-attribute-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          field: 'test.test.test',
        },
        label: 'Custom label',
        name: 'custom_filter_extra_infos_attribute_name',
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
