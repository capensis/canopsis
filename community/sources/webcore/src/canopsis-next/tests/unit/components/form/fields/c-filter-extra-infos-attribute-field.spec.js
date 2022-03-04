import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createInputStub } from '@unit/stubs/input';

import CFilterExtraInfosAttributeField from '@/components/forms/fields/c-filter-extra-infos-attribute-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
};

const factory = (options = {}) => shallowMount(CFilterExtraInfosAttributeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CFilterExtraInfosAttributeField, {
  localVue,

  ...options,
});

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-filter-extra-infos-attribute-field', () => {
  it('Field changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        value: 'extra_infos.test.name',
      },
    });
    const fieldSelect = selectTextField(wrapper);

    expect(fieldSelect.vm.value).toBe('test.name');

    const newField = 'newField';

    fieldSelect.vm.$emit('input', newField);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(`extra_infos.${newField}`);
  });

  it('Field changed after trigger the text field without initial field', () => {
    const wrapper = factory({
      propsData: {
        value: 'extra_infos',
      },
    });
    const fieldSelect = selectTextField(wrapper);

    const newField = 'newField';

    fieldSelect.vm.$emit('input', newField);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(`extra_infos.${newField}`);
  });

  it('Renders `c-filter-extra-infos-attribute-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'extra_infos',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-filter-extra-infos-attribute-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'extra_infos.test.test.test',
        label: 'Custom label',
        name: 'custom_filter_extra_infos_attribute_name',
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
