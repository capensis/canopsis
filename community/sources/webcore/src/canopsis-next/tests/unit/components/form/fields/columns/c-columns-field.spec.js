import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import CColumnsField from '@/components/forms/fields/columns/c-columns-field.vue';

const localVue = createVueInstance();

const defaultColumn = {
  label: 'Default column label',
  value: 'column.defaultColumn',
  isHtml: false,
};
const columnWithHtml = {
  label: 'Column with html label',
  value: 'column.columnWithHtml',
  isHtml: true,
};
const columnWithTemplate = {
  label: 'Column with template',
  value: 'column.columnWithTemplate',
  template: '{{ value }}',
};
const columnWithColorIndicatorState = {
  label: 'Column with color indicator state label',
  value: 'column.columnWithColorIndicatorState',
  isHtml: false,
  colorIndicator: COLOR_INDICATOR_TYPES.state,
};
const columnWithColorIndicatorImpactState = {
  label: 'Column with color indicator impact state label',
  value: 'column.columnWithColorIndicatorImpactState',
  isHtml: false,
  colorIndicator: COLOR_INDICATOR_TYPES.impactState,
};
const columnWithColorIndicatorAndHtmlState = {
  label: 'Column with color indicator and html label',
  value: 'column.columnWithColorIndicatorAndHtml',
  isHtml: true,
  colorIndicator: COLOR_INDICATOR_TYPES.impactState,
};

const stubs = {
  'column-field': true,
};

const selectColumnFieldByIndex = (wrapper, index) => wrapper.findAll('column-field-stub').at(index);
const selectAddButton = wrapper => wrapper.find('v-btn-stub');

describe('c-columns-field', () => {
  const factory = generateShallowRenderer(CColumnsField, {
    localVue,
    stubs,
  });
  const snapshotFactory = generateRenderer(CColumnsField, {
    localVue,
    stubs,

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Column value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          defaultColumn,
        ],
      },
    });

    selectColumnFieldByIndex(wrapper, 0).vm.$emit('input', columnWithColorIndicatorImpactState);

    expect(wrapper).toEmit('input', [columnWithColorIndicatorImpactState]);
  });

  it('Empty column added after click on add button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          columnWithColorIndicatorAndHtmlState,
        ],
      },
    });

    selectAddButton(wrapper).vm.$emit('click', new Event('click'));

    expect(wrapper).toEmit(
      'input',
      [columnWithColorIndicatorAndHtmlState, { label: '', value: '' }],
    );
  });

  it('Empty column with isHtml field added after click on add button', () => {
    const wrapper = factory({
      propsData: {
        withHtml: true,
        columns: [
          columnWithColorIndicatorAndHtmlState,
        ],
      },
    });

    selectAddButton(wrapper).vm.$emit('click', new Event('click'));

    expect(wrapper).toEmit(
      'input',
      [columnWithColorIndicatorAndHtmlState, { label: '', value: '', isHtml: false }],
    );
  });

  it('Column removed after emit click event', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          columnWithColorIndicatorAndHtmlState,
          defaultColumn,
          columnWithColorIndicatorState,
        ],
      },
    });

    selectColumnFieldByIndex(wrapper, 2).vm.$emit('remove');

    expect(wrapper).toEmit('input', [
      columnWithColorIndicatorAndHtmlState,
      defaultColumn,
    ]);
  });

  it('Column moved above after click on top button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          columnWithColorIndicatorAndHtmlState,
          defaultColumn,
          columnWithColorIndicatorState,
        ],
      },
    });

    selectColumnFieldByIndex(wrapper, 1).vm.$emit('up');

    expect(wrapper).toEmit('input', [
      defaultColumn,
      columnWithColorIndicatorAndHtmlState,
      columnWithColorIndicatorState,
    ]);
  });

  it('First column didn\'t move above after click on top button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          columnWithColorIndicatorAndHtmlState,
          defaultColumn,
          columnWithColorIndicatorState,
        ],
      },
    });

    selectColumnFieldByIndex(wrapper, 0).vm.$emit('up');

    expect(wrapper).not.toEmit('input');
  });

  it('Column moved below after click on down button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          columnWithColorIndicatorAndHtmlState,
          defaultColumn,
          columnWithColorIndicatorState,
        ],
      },
    });

    selectColumnFieldByIndex(wrapper, 1).vm.$emit('down');

    expect(wrapper).toEmit('input', [
      columnWithColorIndicatorAndHtmlState,
      columnWithColorIndicatorState,
      defaultColumn,
    ]);
  });

  it('Last column didn\'t move below after click on down button', () => {
    const wrapper = factory({
      propsData: {
        columns: [
          columnWithColorIndicatorAndHtmlState,
          defaultColumn,
          columnWithColorIndicatorState,
        ],
      },
    });

    selectColumnFieldByIndex(wrapper, 2).vm.$emit('down');

    expect(wrapper).not.toEmit('input');
  });

  it('Renders `c-columns-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-columns-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        columns: [
          defaultColumn,
          columnWithHtml,
          columnWithColorIndicatorState,
          columnWithColorIndicatorImpactState,
          columnWithColorIndicatorAndHtmlState,
          columnWithTemplate,
        ],
        withTemplate: true,
        withColorIndicator: true,
        withHtml: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
