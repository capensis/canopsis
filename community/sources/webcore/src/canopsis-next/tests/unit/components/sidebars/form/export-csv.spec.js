import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_FIELDS, ENTITIES_TYPES, EXPORT_CSV_DATETIME_FORMATS, CSV_SEPARATORS } from '@/constants';

import ExportCsv from '@/components/sidebars/form/export-csv.vue';

const stubs = {
  'c-columns-with-template-field': true,
  'widget-settings-item': true,
  'c-csv-separator-field': true,
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'widget-settings-item': true,
  'c-csv-separator-field': true,
  'c-columns-with-template-field': true,
};

const selectColumnsWithTemplateField = wrapper => wrapper.find('c-columns-with-template-field-stub');
const selectSeparatorSelectField = wrapper => wrapper.find('c-csv-separator-field-stub');
const selectDatetimeFormatSelectField = wrapper => wrapper.find('select.v-select');

describe('export-csv', () => {
  const columns = [{
    label: 'Column label',
    value: ALARM_FIELDS.displayName,
    isHtml: false,
  }];

  const factory = generateShallowRenderer(ExportCsv, { stubs });
  const snapshotFactory = generateRenderer(ExportCsv, {
    stubs: snapshotStubs,

    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

  it('Separator changed after trigger separator select field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        form: {
          exportCsvSeparator: CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: columns,
        },
      },
    });

    selectSeparatorSelectField(wrapper).triggerCustomEvent('input', CSV_SEPARATORS.semicolon);

    expect(wrapper).toEmitInput({
      exportCsvSeparator: CSV_SEPARATORS.semicolon,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
      widgetExportColumns: columns,
    });
  });

  it('Datetime format changed after trigger datetime format select field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        form: {
          exportCsvSeparator: CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: columns,
        },
        datetimeFormat: true,
      },
    });

    selectDatetimeFormatSelectField(wrapper).setValue(EXPORT_CSV_DATETIME_FORMATS.dayOfMonthMonthNameYearTime.value);

    expect(wrapper).toEmitInput({
      exportCsvSeparator: CSV_SEPARATORS.comma,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.dayOfMonthMonthNameYearTime.value,
      widgetExportColumns: columns,
    });
  });

  it('Columns changed after trigger columns field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        form: {
          exportCsvSeparator: CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: [],
        },
      },
    });

    selectColumnsWithTemplateField(wrapper).triggerCustomEvent('input', columns);

    expect(wrapper).toEmitInput({
      exportCsvSeparator: CSV_SEPARATORS.comma,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
      widgetExportColumns: columns,
    });
  });

  it('Renders `export-csv` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
      },
    });

    const menuContents = wrapper.findAllMenus();

    expect(wrapper).toMatchSnapshot();
    menuContents.wrappers.forEach((menuContent) => {
      expect(menuContent.element).toMatchSnapshot();
    });
  });

  it('Renders `export-csv` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        form: {
          exportCsvSeparator: CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: columns,
        },
        datetimeFormat: true,
      },
    });

    const menuContents = wrapper.findAllMenus();

    expect(wrapper).toMatchSnapshot();
    menuContents.wrappers.forEach((menuContent) => {
      expect(menuContent.element).toMatchSnapshot();
    });
  });
});
