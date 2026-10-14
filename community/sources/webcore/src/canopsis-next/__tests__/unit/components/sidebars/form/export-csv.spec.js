import { generateShallowRenderer } from '@unit/utils/vue';

import {
  ALARM_FIELDS,
  ENTITIES_TYPES,
  EXPORT_CSV_DATETIME_FORMATS,
  CSV_SEPARATORS,
} from '@/constants';

import ExportCsv from '@/components/sidebars/form/export-csv.vue';

const stubs = {
  'c-columns-with-template-field': true,
  'widget-settings-item': true,
  'c-csv-separator-field': true,
};

const selectColumnsWithTemplateField = wrapper =>
  wrapper.find('c-columns-with-template-field-stub');
const selectSeparatorSelectField = wrapper => wrapper.find('c-csv-separator-field-stub');
const selectDatetimeFormatSelectField = wrapper => wrapper.find('.c-select');

describe('export-csv', () => {
  const columns = [
    {
      label: 'Column label',
      value: ALARM_FIELDS.displayName,
      isHtml: false,
    },
  ];

  const factory = generateShallowRenderer(ExportCsv, { stubs });

  it('Separator changed after trigger separator select field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        value: {
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
        value: {
          exportCsvSeparator: CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: columns,
        },
        datetimeFormat: true,
      },
    });

    selectDatetimeFormatSelectField(wrapper).triggerCustomEvent(
      'input',
      EXPORT_CSV_DATETIME_FORMATS.dayOfMonthMonthNameYearTime.value,
    );

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
        value: {
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
});
