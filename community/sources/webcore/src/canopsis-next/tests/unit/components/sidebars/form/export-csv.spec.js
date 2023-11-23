import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_FIELDS, ENTITIES_TYPES, EXPORT_CSV_DATETIME_FORMATS, EXPORT_CSV_SEPARATORS } from '@/constants';

import ExportCsv from '@/components/sidebars/form/export-csv.vue';

const stubs = {
  'c-columns-with-template-field': true,
  'widget-settings-item': true,
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'widget-settings-item': true,
  'c-columns-with-template-field': true,
};

const selectColumnsWithTemplateField = wrapper => wrapper.find('c-columns-with-template-field-stub');
const selectSeparatorSelectField = wrapper => wrapper.findAll('select.v-select').at(0);
const selectDatetimeFormatSelectField = wrapper => wrapper.findAll('select.v-select').at(1);

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
          exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: columns,
        },
      },
    });

    const separatorField = selectSeparatorSelectField(wrapper);

    separatorField.setValue(EXPORT_CSV_SEPARATORS.semicolon);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      exportCsvSeparator: EXPORT_CSV_SEPARATORS.semicolon,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
      widgetExportColumns: columns,
    });
  });

  it('Datetime format changed after trigger datetime format select field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        form: {
          exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: columns,
        },
        datetimeFormat: true,
      },
    });

    const datetimeFormatSelectField = selectDatetimeFormatSelectField(wrapper);

    datetimeFormatSelectField.setValue(EXPORT_CSV_DATETIME_FORMATS.dayOfMonthMonthNameYearTime.value);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.dayOfMonthMonthNameYearTime.value,
      widgetExportColumns: columns,
    });
  });

  it('Columns changed after trigger columns field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        form: {
          exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
          exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
          widgetExportColumns: [],
        },
      },
    });

    const columnsField = selectColumnsWithTemplateField(wrapper);

    columnsField.vm.$emit('input', columns);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
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
          exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
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
