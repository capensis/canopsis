import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';

import PbehaviorExceptionField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exception-field.vue';

const stubs = {
  'date-time-splitted-range-picker-field': true,
  'date-time-splitted-range-picker-text': true,
  'c-pbehavior-type-field': true,
  'c-pbehavior-type-text': true,
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
};

const snapshotStubs = {
  'date-time-splitted-range-picker-field': true,
  'date-time-splitted-range-picker-text': true,
  'c-pbehavior-type-field': true,
  'c-pbehavior-type-text': true,
  'v-checkbox': true,
};

const selectDateTimePickerField = wrapper => wrapper.find('date-time-splitted-range-picker-field-stub');
const selectTypeField = wrapper => wrapper.find('c-pbehavior-type-field-stub');
const selectToggleEditButton = wrapper => wrapper.findAll('v-btn-stub').at(0);
const selectRemoveButton = wrapper => wrapper.findAll('v-btn-stub').at(1);
const selectFullDayCheckbox = wrapper => wrapper.find('.v-checkbox');

describe('pbehavior-exception-field', () => {
  const factory = generateShallowRenderer(PbehaviorExceptionField, {

    stubs,
  });
  const snapshotFactory = generateRenderer(PbehaviorExceptionField, {

    stubs: snapshotStubs,
  });

  test('Begin date changed after trigger date time picker', () => {
    const value = {
      begin: new Date(123),
      end: new Date(321),
      type: '',
    };
    const wrapper = factory({
      propsData: { value },
    });

    const newTimeStart = new Date(1233266565);

    selectDateTimePickerField(wrapper).triggerCustomEvent('update:start', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...value,
      begin: newTimeStart,
    });
  });

  test('Begin date changed after trigger date time picker with defined type', async () => {
    const value = {
      begin: new Date(123),
      end: new Date(321),
      type: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { value },
    });

    const newTimeStart = new Date(1233266565);

    selectToggleEditButton(wrapper).triggerCustomEvent('click');

    await flushPromises();

    selectDateTimePickerField(wrapper).triggerCustomEvent('update:start', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...value,
      begin: newTimeStart,
    });
  });

  test('End date changed after trigger date time picker', async () => {
    const value = {
      begin: new Date(123),
      end: new Date(321),
      type: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { value },
    });

    const newTimeStart = new Date(1233266565);

    selectToggleEditButton(wrapper).triggerCustomEvent('click');

    await flushPromises();

    selectDateTimePickerField(wrapper).triggerCustomEvent('update:end', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...value,
      end: newTimeStart,
    });
  });

  test('Type changed after trigger type field', async () => {
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

    selectToggleEditButton(wrapper).triggerCustomEvent('click');

    await flushPromises();

    selectTypeField(wrapper).triggerCustomEvent('input', newType);

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

    selectRemoveButton(wrapper).triggerCustomEvent('click');

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

    selectFullDayCheckbox(wrapper).triggerCustomEvent('change', true);
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

    selectFullDayCheckbox(wrapper).triggerCustomEvent('change', false);
    await flushPromises();

    expect(wrapper).toEmit('input', value);
  });

  test('Renders `pbehavior-exception-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehavior-exception-field` with custom props with defined type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          begin: new Date(123),
          end: new Date(321),
          type: 'type',
        },
        disabled: true,
        withType: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
