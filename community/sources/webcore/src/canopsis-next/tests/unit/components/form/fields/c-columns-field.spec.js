import Faker from 'faker';
import { Validator } from 'vee-validate';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import CColumnsField from '@/components/forms/fields/c-columns-field.vue';

const localVue = createVueInstance();

const mockData = {
  defaultColumn: {
    label: 'Default column label',
    value: 'column.defaultColumn',
    isHtml: false,
  },
  columnWithHtml: {
    label: 'Column with html label',
    value: 'column.columnWithHtml',
    isHtml: true,
  },
  columnWithTemplate: {
    label: 'Column with template',
    value: 'column.columnWithTemplate',
    template: '{{ value }}',
  },
  columnWithColorIndicatorState: {
    label: 'Column with color indicator state label',
    value: 'column.columnWithColorIndicatorState',
    isHtml: false,
    colorIndicator: COLOR_INDICATOR_TYPES.state,
  },
  columnWithColorIndicatorImpactState: {
    label: 'Column with color indicator impact state label',
    value: 'column.columnWithColorIndicatorImpactState',
    isHtml: false,
    colorIndicator: COLOR_INDICATOR_TYPES.impactState,
  },
  columnWithColorIndicatorAndHtmlState: {
    label: 'Column with color indicator and html label',
    value: 'column.columnWithColorIndicatorAndHtml',
    isHtml: true,
    colorIndicator: COLOR_INDICATOR_TYPES.impactState,
  },
};

const stubs = {
  'c-color-indicator-field': {
    props: ['value'],
    template: `
      <input
        :value="value"
        class="c-color-indicator-field"
        @input="$listeners.input($event.target.value)"
      />
    `,
  },
  'v-switch': {
    props: ['inputValue'],
    template: `
      <input
        :checked="inputValue"
        type="checkbox"
        class="v-switch"
        @change="$listeners.change($event.target.checked)" />
    `,
  },
  'v-text-field': {
    props: ['value'],
    template: `
      <input
        :value="value"
        class="v-text-field"
        @input="$listeners.input($event.target.value)"
      />
    `,
  },
  'v-btn': {
    template: `
      <button class="v-btn" @click="$listeners.click">
        <slot></slot>
      </button>
    `,
  },
  'c-help-icon': true,
};

const snapshotStubs = {
  'c-color-indicator-field': true,
  'c-help-icon': true,
  'v-text-field': {
    template: '<input class="v-text-field" />',
  },
};

const factory = (options = {}) => shallowMount(CColumnsField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CColumnsField, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

describe('c-columns-field', () => {
  it('Column value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.defaultColumn,
        ],
      },
    });
    const newColumnValue = Faker.datatype.string();
    const firstColumnElement = wrapper.findAll('v-card-stub').at(0);
    const inputValueElement = firstColumnElement.findAll('input.v-text-field').at(1);

    expect(inputValueElement.element.value).toBe(mockData.defaultColumn.value);

    inputValueElement.setValue(newColumnValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [firstColumn] = newColumnsData;

    expect(firstColumn.value).toBe(newColumnValue);
    expect(firstColumn).toEqual({
      ...mockData.defaultColumn,
      value: newColumnValue,
    });
  });

  it('Column label changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.defaultColumn,
        ],
      },
    });
    const newColumnLabel = Faker.datatype.string();
    const firstColumnElement = wrapper.findAll('v-card-stub').at(0);
    const inputLabelElement = firstColumnElement.findAll('input.v-text-field').at(0);

    expect(inputLabelElement.element.value).toBe(mockData.defaultColumn.label);

    inputLabelElement.setValue(newColumnLabel);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [firstColumn] = newColumnsData;

    expect(firstColumn.label).toBe(newColumnLabel);
    expect(firstColumn).toEqual({
      ...mockData.defaultColumn,
      label: newColumnLabel,
    });
  });

  it('Column isHtml value changed after trigger the switch', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        columns: [
          mockData.defaultColumn,
        ],
      },
    });
    const newColumnIsHtml = !mockData.defaultColumn.isHtml;
    const firstColumnElement = wrapper.findAll('v-card-stub').at(0);
    const switchIsHtmlElement = firstColumnElement.findAll('input.v-switch').at(0);

    switchIsHtmlElement.setChecked(newColumnIsHtml);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [firstColumn] = newColumnsData;

    expect(firstColumn.isHtml).toBe(newColumnIsHtml);
    expect(firstColumn).toEqual({
      ...mockData.defaultColumn,
      isHtml: newColumnIsHtml,
    });
  });

  it('Column colorIndicator enabled after trigger the switch', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        withColorIndicator: true,
        columns: [
          mockData.defaultColumn,
        ],
      },
    });
    const firstColumnElement = wrapper.findAll('v-card-stub').at(0);
    const firstColumnSwitchColorIndicatorElement = firstColumnElement.findAll('input.v-switch').at(1);

    firstColumnSwitchColorIndicatorElement.setChecked(true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [firstColumn] = newColumnsData;

    expect(firstColumn.colorIndicator).toBe(COLOR_INDICATOR_TYPES.state);
    expect(firstColumn).toEqual({
      ...mockData.defaultColumn,
      colorIndicator: COLOR_INDICATOR_TYPES.state,
    });
  });

  it('Column colorIndicator disabled after trigger the switch', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        withColorIndicator: true,
        columns: [
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const firstColumnElement = wrapper.findAll('v-card-stub').at(0);
    const firstColumnSwitchColorIndicatorElement = firstColumnElement.findAll('input.v-switch').at(1);

    firstColumnSwitchColorIndicatorElement.setChecked(false);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [firstColumn] = newColumnsData;

    expect(firstColumn.colorIndicator).toBe(null);
    expect(firstColumn).toEqual({
      ...mockData.columnWithColorIndicatorState,
      colorIndicator: null,
    });
  });

  it('Column colorIndicator value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        withColorIndicator: true,
        columns: [
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const firstColumnElement = wrapper.findAll('v-card-stub').at(0);
    const inputColorIndicatorElement = firstColumnElement.find('input.c-color-indicator-field');

    inputColorIndicatorElement.setValue(COLOR_INDICATOR_TYPES.impactState);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [firstColumn] = newColumnsData;

    expect(firstColumn.colorIndicator).toBe(COLOR_INDICATOR_TYPES.impactState);
    expect(firstColumn).toEqual({
      ...mockData.columnWithColorIndicatorState,
      colorIndicator: COLOR_INDICATOR_TYPES.impactState,
    });
  });

  it('Empty column added after click on add button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
        ],
      },
    });
    const addButtonElement = wrapper.find('v-container-stub > button.v-btn');

    addButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [oldColumn, newColumn] = newColumnsData;

    expect(oldColumn).toEqual(mockData.columnWithColorIndicatorAndHtmlState);
    expect(newColumn).toEqual({ label: '', value: '' });
  });

  it('Empty column with isHtml field added after click on add button', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
        ],
      },
    });
    const addButtonElement = wrapper.find('v-container-stub > button.v-btn');

    addButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];
    const [oldColumn, newColumn] = newColumnsData;

    expect(oldColumn).toEqual(mockData.columnWithColorIndicatorAndHtmlState);
    expect(newColumn).toEqual({ label: '', value: '', isHtml: false });
  });

  it('Column removed after click on remove button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.defaultColumn,
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const secondColumnElement = wrapper.findAll('v-card-stub').at(1);
    const secondColumnRemoveButtonElement = secondColumnElement.findAll('button.v-btn').at(2);

    secondColumnRemoveButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];

    expect(newColumnsData).toEqual([
      mockData.columnWithColorIndicatorAndHtmlState,
      mockData.columnWithColorIndicatorState,
    ]);
  });

  it('Column moved above after click on top button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.defaultColumn,
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const secondColumnElement = wrapper.findAll('v-card-stub').at(1);
    const secondColumnTopButtonElement = secondColumnElement.findAll('button.v-btn').at(0);

    secondColumnTopButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];

    expect(newColumnsData).toEqual([
      mockData.defaultColumn,
      mockData.columnWithColorIndicatorAndHtmlState,
      mockData.columnWithColorIndicatorState,
    ]);
  });

  it('First column didn\'t move above after click on top button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.defaultColumn,
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const secondColumnElement = wrapper.findAll('v-card-stub').at(0);
    const secondColumnTopButtonElement = secondColumnElement.findAll('button.v-btn').at(0);

    secondColumnTopButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toBeFalsy();
  });

  it('Column moved below after click on down button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.defaultColumn,
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const secondColumnElement = wrapper.findAll('v-card-stub').at(1);
    const secondColumnDownButtonElement = secondColumnElement.findAll('button.v-btn').at(1);

    secondColumnDownButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [newColumnsData] = inputEvents[0];

    expect(newColumnsData).toEqual([
      mockData.columnWithColorIndicatorAndHtmlState,
      mockData.columnWithColorIndicatorState,
      mockData.defaultColumn,
    ]);
  });

  it('Last column didn\'t move below after click on down button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.defaultColumn,
          mockData.columnWithColorIndicatorState,
        ],
      },
    });
    const secondColumnElement = wrapper.findAll('v-card-stub').at(2);
    const secondColumnDownButtonElement = secondColumnElement.findAll('button.v-btn').at(1);

    secondColumnDownButtonElement.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toBeFalsy();
  });

  it('Renders `c-columns-field` with default props correctly', () => {
    const wrapper = mount(CColumnsField, {
      localVue,
      stubs: snapshotStubs,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-columns-field` with all columns type correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        columns: [
          mockData.defaultColumn,
          mockData.columnWithHtml,
          mockData.columnWithColorIndicatorState,
          mockData.columnWithColorIndicatorImpactState,
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.columnWithTemplate,
        ],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-columns-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        columns: [
          mockData.defaultColumn,
          mockData.columnWithHtml,
          mockData.columnWithColorIndicatorState,
          mockData.columnWithColorIndicatorImpactState,
          mockData.columnWithColorIndicatorAndHtmlState,
          mockData.columnWithTemplate,
        ],
        withTemplate: true,
        withColorIndicator: true,
        withHtml: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-columns-field` with errors correctly', () => {
    const columns = [
      mockData.defaultColumn,
      mockData.columnWithHtml,
      mockData.columnWithColorIndicatorState,
      mockData.columnWithColorIndicatorImpactState,
      mockData.columnWithColorIndicatorAndHtmlState,
    ];
    const validator = new Validator();
    validator.errors.add([
      {
        field: 'label[1]',
        msg: 'Label 1 error',
      },
      {
        field: 'value[1]',
        msg: 'Value 1 error',
      },
      {
        field: 'value[3]',
        msg: 'Value 3 error',
      },
      {
        field: 'label[4]',
        msg: 'Label 4 error',
      },
    ]);

    const wrapper = mount(CColumnsField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        columns,
        withColorIndicator: true,
        withHtml: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
