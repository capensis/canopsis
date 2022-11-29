import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createInputStub } from '@unit/stubs/input';
import { COLOR_INDICATOR_TYPES } from '@/constants';

import ColumnField from '@/components/forms/fields/columns/column-field.vue';

const localVue = createVueInstance();

const defaultColumn = {
  label: 'Default column label',
  value: 'column.defaultColumn',
  isHtml: false,
};
const columnWithColorIndicatorState = {
  label: 'Column with color indicator state label',
  value: 'column.columnWithColorIndicatorState',
  isHtml: false,
  colorIndicator: COLOR_INDICATOR_TYPES.state,
};
const columnWithColorIndicatorAndHtmlState = {
  label: 'Column with color indicator and html label',
  value: 'column.columnWithColorIndicatorAndHtml',
  isHtml: true,
  colorIndicator: COLOR_INDICATOR_TYPES.impactState,
};

const stubs = {
  'c-color-indicator-field': true,
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'c-color-indicator-field': true,
};

const selectLabelField = wrapper => wrapper.findAll('.v-text-field').at(0);
const selectValueField = wrapper => wrapper.findAll('.v-text-field').at(1);
const selectIsHTMLSwitchField = wrapper => wrapper.findAll('v-switch-stub').at(0);
const selectColorIndicatorSwitchField = wrapper => wrapper.findAll('v-switch-stub').at(1);
const selectColorIndicatorField = wrapper => wrapper.find('c-color-indicator-field-stub');
const selectRemoveButton = wrapper => wrapper.findAll('v-btn-stub').at(2);
const selectDownButton = wrapper => wrapper.findAll('v-btn-stub').at(1);
const selectUpButton = wrapper => wrapper.findAll('v-btn-stub').at(0);

describe('column-field', () => {
  const factory = generateShallowRenderer(ColumnField, {
    localVue,
    stubs,
  });
  const snapshotFactory = generateRenderer(ColumnField, {
    localVue,
    stubs: snapshotStubs,

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Column value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        column: defaultColumn,
      },
    });
    const newColumnValue = Faker.datatype.string();

    selectValueField(wrapper).setValue(newColumnValue);

    expect(wrapper).toEmit('input', {
      ...defaultColumn,
      value: newColumnValue,
    });
  });

  it('Column label changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        column: defaultColumn,
      },
    });
    const newColumnLabel = Faker.datatype.string();

    selectLabelField(wrapper).setValue(newColumnLabel);

    expect(wrapper).toEmit('input', {
      ...defaultColumn,
      label: newColumnLabel,
    });
  });

  it('Column isHtml value changed after trigger the switch', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        column: defaultColumn,
      },
    });
    const newColumnIsHtml = !defaultColumn.isHtml;

    selectIsHTMLSwitchField(wrapper).vm.$emit('change', newColumnIsHtml);

    expect(wrapper).toEmit('input', {
      ...defaultColumn,
      isHtml: newColumnIsHtml,
    });
  });

  it('Column colorIndicator enabled after trigger the switch', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        withColorIndicator: true,
        column: defaultColumn,
      },
    });

    selectColorIndicatorSwitchField(wrapper).vm.$emit('change', true);

    expect(wrapper).toEmit('input', {
      ...defaultColumn,
      colorIndicator: COLOR_INDICATOR_TYPES.state,
    });
  });

  it('Column colorIndicator disabled after trigger the switch', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        withColorIndicator: true,
        column: columnWithColorIndicatorState,
      },
    });
    selectColorIndicatorSwitchField(wrapper).vm.$emit('change', false);

    expect(wrapper).toEmit('input', {
      ...columnWithColorIndicatorState,
      colorIndicator: null,
    });
  });

  it('Column colorIndicator value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        withColorIndicator: true,
        column: columnWithColorIndicatorState,
      },
    });
    selectColorIndicatorField(wrapper).vm.$emit('input', COLOR_INDICATOR_TYPES.impactState);

    expect(wrapper).toEmit('input', {
      ...columnWithColorIndicatorState,
      colorIndicator: COLOR_INDICATOR_TYPES.impactState,
    });
  });

  it('Remove event emitted after click on remove button', () => {
    const wrapper = factory({
      propsData: {
        column: columnWithColorIndicatorAndHtmlState,
      },
    });

    selectRemoveButton(wrapper).vm.$emit('click');

    expect(wrapper).toEmit('remove');
  });

  it('Up event emitted after click on up button', () => {
    const wrapper = factory({
      propsData: {
        column: columnWithColorIndicatorState,
      },
    });

    selectUpButton(wrapper).vm.$emit('click');

    expect(wrapper).toEmit('up');
  });

  it('Down event emitted after click on up button', () => {
    const wrapper = factory({
      propsData: {
        column: columnWithColorIndicatorState,
      },
    });

    selectDownButton(wrapper).vm.$emit('click');

    expect(wrapper).toEmit('down');
  });

  it('Renders `column-field` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        column: defaultColumn,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `column-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        column: columnWithColorIndicatorAndHtmlState,
        withTemplate: true,
        withColorIndicator: true,
        withHtml: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
