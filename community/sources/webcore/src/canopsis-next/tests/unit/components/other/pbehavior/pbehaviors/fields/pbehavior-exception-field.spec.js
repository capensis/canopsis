import Faker from 'faker';
import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';

import PbehaviorExceptionField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-exception-field.vue';

const localVue = createVueInstance();

const stubs = {
  'date-time-splitted-range-picker-field': true,
  'c-pbehavior-type-field': true,
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
};

const snapshotStubs = {
  'date-time-splitted-range-picker-field': true,
  'c-pbehavior-type-field': true,
};

const selectDateTimePickerField = wrapper => wrapper.find('date-time-splitted-range-picker-field-stub');
const selectTypeField = wrapper => wrapper.find('c-pbehavior-type-field-stub');
const selectRemoveButton = wrapper => wrapper.find('v-btn-stub');
const selectFullDayCheckbox = wrapper => wrapper.find('.v-checkbox');

describe('pbehavior-exception-field', () => {
  const factory = generateShallowRenderer(PbehaviorExceptionField, {
    localVue,
    stubs,
  });
  const snapshotFactory = generateRenderer(PbehaviorExceptionField, {
    localVue,
    stubs: snapshotStubs,
  });

  test('Begin date changed after trigger date time picker', () => {
    const value = {
      begin: new Date(123),
      end: new Date(321),
      type: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { value },
    });

    const newTimeStart = new Date(1233266565);

    selectDateTimePickerField(wrapper).vm.$emit('update:start', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...value,
      begin: newTimeStart,
    });
  });

  test('End date changed after trigger date time picker', () => {
    const value = {
      begin: new Date(123),
      end: new Date(321),
      type: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { value },
    });

    const newTimeStart = new Date(1233266565);

    selectDateTimePickerField(wrapper).vm.$emit('update:end', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...value,
      end: newTimeStart,
    });
  });

  test('Type changed after trigger type field', () => {
    const value = {
      begin: new Date(123),
      end: new Date(321),
      type: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: {
        value,
        withType: true,
      },
    });

    const newType = {
      _id: Faker.datatype.string(),
    };

    selectTypeField(wrapper).vm.$emit('input', newType);

    expect(wrapper).toEmit('input', {
      ...value,
      type: newType,
    });
  });

  test('Delete event emitted after trigger button', () => {
    const wrapper = factory({
      propsData: {
        value: {},
        withType: true,
      },
    });

    selectRemoveButton(wrapper).vm.$emit('click');

    expect(wrapper).toEmit('delete');
  });

  test('Full day enabled after trigger checkbox', async () => {
    const value = {
      begin: new Date(1610393200000),
      end: new Date(1757480799000),
    };
    const wrapper = factory({
      propsData: {
        value,
        withType: true,
      },
    });

    selectFullDayCheckbox(wrapper).vm.$emit('change', true);
    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...value,

      begin: new Date(1610319600000),
      end: new Date(1757541599000),
    });
  });

  test('Full day disabled after trigger checkbox', async () => {
    const value = {
      begin: new Date(1610319600000),
      end: new Date(1757541599000),
    };
    const wrapper = factory({
      propsData: {
        value,
        withType: true,
      },
    });

    selectFullDayCheckbox(wrapper).vm.$emit('change', false);
    await flushPromises();

    expect(wrapper).toEmit('input', value);
  });

  test('Renders `pbehavior-exception-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-exception-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          begin: new Date(123),
          end: new Date(321),
          type: '',
        },
        disabled: true,
        withType: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
