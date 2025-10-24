import Faker from 'faker';

import { RESIZING_CELLS_CONTENTS_BEHAVIORS, COLOR_INDICATOR_TYPES } from '@/constants';

import { widgetColumnsAlarmMixin } from '@/mixins/widget/columns/alarm';
import { widgetColumnsEntityInfoPropertyMixin } from '@/mixins/widget/columns/entity-info-property';

describe('widget-columns-alarm', () => {
  const mockColumns = [
    {
      value: 'v.connector',
      label: '',
      text: 'Connector',
    },
    {
      value: 'entity.infos.cpu_usage.value',
      label: '',
      text: 'CPU',
    },
    {
      value: 'entity.infos.memory_usage.value',
      label: 'Custom Memory Label',
      text: 'Memory',
    },
    {
      value: 'v.state.val',
      label: '',
      text: 'State',
      colorIndicator: COLOR_INDICATOR_TYPES.state,
    },
  ];

  const mockEntityInfoProperties = [
    {
      _id: Faker.datatype.string(),
      name: 'cpu_usage',
      alias: 'CPU Usage',
      description: Faker.datatype.string(),
    },
    {
      _id: Faker.datatype.string(),
      name: 'memory_usage',
      alias: 'Memory Usage',
      description: Faker.datatype.string(),
    },
  ];

  const mockColumnsFilters = [
    {
      column: 'v.connector',
      filter: 'testFilter',
      attributes: ['attr1', 'attr2'],
    },
  ];

  let mockContext;

  beforeEach(() => {
    mockContext = {
      widget: {
        parameters: {
          showRootCauseByStateClick: true,
          columns: {
            cells_content_behavior: RESIZING_CELLS_CONTENTS_BEHAVIORS.truncate,
          },
          infoPopups: [
            {
              column: 'v.connector',
              template: 'Connector popup template',
            },
            {
              column: 'alarm.v.state.val',
              template: 'State popup template',
            },
          ],
        },
      },
      columns: mockColumns,
      entityInfoProperties: mockEntityInfoProperties,
      columnsFilters: mockColumnsFilters,
      columnsFiltersPending: false,
      fetchAlarmColumnsFiltersList: jest.fn().mockResolvedValue(mockColumnsFilters),
      fetchEntityInfoPropertiesList: jest.fn(),
      $i18n: {
        locale: 'en',
      },
      $options: {
        filters: {
          testFilter: jest.fn((value, ...attrs) => `filtered_${value}_${attrs.join('_')}`),
        },
      },
    };
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('data', () => {
    test('Should return initial data with empty columnsFilters and false columnsFiltersPending', () => {
      const data = widgetColumnsAlarmMixin.data.call(mockContext);

      expect(data).toEqual({
        columnsFilters: [],
        columnsFiltersPending: false,
      });
    });
  });

  describe('computed', () => {
    describe('infoPopupsMap', () => {
      test('Should create a map of info popups by column', () => {
        const { infoPopupsMap } = widgetColumnsAlarmMixin.computed;

        const result = infoPopupsMap.call(mockContext);

        expect(result).toEqual({
          'v.connector': 'Connector popup template',
          'alarm.v.state.val': 'State popup template',
        });
      });

      test('Should return empty object when infoPopups is undefined', () => {
        const { infoPopupsMap } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          widget: {
            parameters: {},
          },
        };

        const result = infoPopupsMap.call(context);

        expect(result).toEqual({});
      });

      test('Should return empty object when widget parameters is undefined', () => {
        const { infoPopupsMap } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          widget: {},
        };

        const result = infoPopupsMap.call(context);

        expect(result).toEqual({});
      });
    });

    describe('columnsFiltersMap', () => {
      test('Should create a map of columns filters with filter functions', () => {
        const { columnsFiltersMap } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          getFilter: widgetColumnsAlarmMixin.methods.getFilter,
        };

        const result = columnsFiltersMap.call(context);

        expect(result['v.connector']).toBeDefined();
        expect(typeof result['v.connector']).toBe('function');
        expect(result['v.connector']('test')).toBe('filtered_test_attr1_attr2');
      });

      test('Should handle columns filters without attributes', () => {
        const { columnsFiltersMap } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          getFilter: widgetColumnsAlarmMixin.methods.getFilter,
          columnsFilters: [
            {
              column: 'v.state.val',
              filter: 'testFilter',
            },
          ],
        };

        const result = columnsFiltersMap.call(context);

        expect(result['v.state.val']).toBeDefined();
        expect(result['v.state.val']('value')).toBe('filtered_value_');
      });

      test('Should return empty object when columnsFilters is undefined', () => {
        const { columnsFiltersMap } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          columnsFilters: undefined,
        };

        const result = columnsFiltersMap.call(context);

        expect(result).toEqual({});
      });

      test('Should return identity function when filter does not exist', () => {
        const { columnsFiltersMap } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          getFilter: widgetColumnsAlarmMixin.methods.getFilter,
          columnsFilters: [
            {
              column: 'v.connector',
              filter: 'nonExistentFilter',
            },
          ],
          $options: {
            filters: {},
          },
        };

        const result = columnsFiltersMap.call(context);

        expect(result['v.connector']('test')).toBe('test');
      });
    });

    describe('preparedColumns', () => {
      test('Should prepare columns with all properties', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const filterFunc = jest.fn(value => value);
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: jest.fn(() => filterFunc),
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result).toHaveLength(4);
        expect(result[0]).toMatchObject({
          value: 'v.connector',
          text: 'Connector',
          popupTemplate: 'Connector popup template',
        });
        expect(result[0].getComponent).toBeDefined();
        expect(result[0].filter).toBeDefined();
      });

      test('Should use alias when column value matches entity info property', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[1].text).toBe('CPU Usage');
      });

      test('Should use label when provided instead of alias', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[2].text).toBe('Custom Memory Label');
      });

      test('Should use original text when no label or alias', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[0].text).toBe('Connector');
      });

      test('Should enable color indicator when colorIndicator is valid', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[3].colorIndicatorEnabled).toBe(true);
      });

      test('Should disable color indicator when colorIndicator is invalid', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[0].colorIndicatorEnabled).toBe(false);
      });

      test('Should handle alarm.* prefix for popup template', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[3].popupTemplate).toBe('State popup template');
      });

      test('Should use truncate behavior when specified', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          infoPopupsMap: widgetColumnsAlarmMixin.computed.infoPopupsMap.call(mockContext),
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[0].getComponent).toBeDefined();
      });

      test('Should handle empty columns array', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          columns: [],
          infoPopupsMap: {},
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result).toEqual([]);
      });

      test('Should handle undefined columns', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          columns: undefined,
          infoPopupsMap: {},
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result).toEqual([]);
      });

      test('Should return undefined filter when locale is not set', () => {
        const { preparedColumns } = widgetColumnsAlarmMixin.computed;
        const context = {
          ...mockContext,
          $i18n: {
            locale: '',
          },
          infoPopupsMap: {},
          getColumnFilter: widgetColumnsAlarmMixin.methods.getColumnFilter,
          findAliasByColumnValue: widgetColumnsEntityInfoPropertyMixin.methods.findAliasByColumnValue,
          columnsFiltersMap: {},
        };

        const result = preparedColumns.call(context);

        expect(result[0].filter).toBeFalsy();
      });
    });
  });

  describe('mounted', () => {
    test('Should fetch column filters and entity info properties on mount', () => {
      const context = {
        ...mockContext,
        fetchColumnFilters: widgetColumnsAlarmMixin.methods.fetchColumnFilters,
      };

      widgetColumnsAlarmMixin.mounted.call(context);

      expect(context.fetchAlarmColumnsFiltersList).toHaveBeenCalled();
      expect(context.fetchEntityInfoPropertiesList).toHaveBeenCalledWith({
        params: { paginate: false },
      });
    });
  });

  describe('methods', () => {
    describe('getColumnFilter', () => {
      test('Should return filter from columnsFiltersMap when it exists', () => {
        const { getColumnFilter } = widgetColumnsAlarmMixin.methods;
        const filterFunc = jest.fn();
        const context = {
          ...mockContext,
          columnsFiltersMap: {
            'v.connector': filterFunc,
          },
        };

        const result = getColumnFilter.call(context, 'v.connector');

        expect(result).toBe(filterFunc);
      });

      test('Should return default filter when column not in columnsFiltersMap', () => {
        const { getColumnFilter } = widgetColumnsAlarmMixin.methods;
        const context = {
          ...mockContext,
          columnsFiltersMap: {},
        };

        const result = getColumnFilter.call(context, 'v.last_update_date');

        expect(result).toBeDefined();
        expect(typeof result).toBe('function');
      });
    });

    describe('getFilter', () => {
      test('Should return filter function with attributes', () => {
        const { getFilter } = widgetColumnsAlarmMixin.methods;

        const result = getFilter.call(mockContext, 'testFilter', ['arg1', 'arg2']);

        expect(typeof result).toBe('function');
        expect(result('value')).toBe('filtered_value_arg1_arg2');
      });

      test('Should return filter function without attributes', () => {
        const { getFilter } = widgetColumnsAlarmMixin.methods;

        const result = getFilter.call(mockContext, 'testFilter');

        expect(typeof result).toBe('function');
        expect(result('value')).toBe('filtered_value_');
      });

      test('Should return identity function when filter does not exist', () => {
        const { getFilter } = widgetColumnsAlarmMixin.methods;

        const result = getFilter.call(mockContext, 'nonExistentFilter');

        expect(typeof result).toBe('function');
        expect(result('value')).toBe('value');
      });
    });

    describe('fetchColumnFilters', () => {
      test('Should fetch column filters and update state', async () => {
        const { fetchColumnFilters } = widgetColumnsAlarmMixin.methods;
        const context = {
          ...mockContext,
          columnsFilters: [],
          columnsFiltersPending: false,
        };

        await fetchColumnFilters.call(context);

        expect(context.fetchAlarmColumnsFiltersList).toHaveBeenCalled();
        expect(context.columnsFilters).toEqual(mockColumnsFilters);
        expect(context.columnsFiltersPending).toBe(false);
      });

      test('Should set pending state during fetch', async () => {
        const { fetchColumnFilters } = widgetColumnsAlarmMixin.methods;
        let pendingDuringFetch = false;
        const context = {
          ...mockContext,
          columnsFilters: [],
          columnsFiltersPending: false,
          fetchAlarmColumnsFiltersList: jest.fn(async () => {
            pendingDuringFetch = context.columnsFiltersPending;
            return mockColumnsFilters;
          }),
        };

        await fetchColumnFilters.call(context);

        expect(pendingDuringFetch).toBe(true);
        expect(context.columnsFiltersPending).toBe(false);
      });

      test('Should handle fetch error gracefully', async () => {
        const { fetchColumnFilters } = widgetColumnsAlarmMixin.methods;
        const error = new Error('Fetch failed');
        const context = {
          ...mockContext,
          columnsFilters: [],
          columnsFiltersPending: false,
          fetchAlarmColumnsFiltersList: jest.fn().mockRejectedValue(error),
        };
        const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation();

        await fetchColumnFilters.call(context);

        expect(consoleWarnSpy).toHaveBeenCalledWith(error);
        expect(context.columnsFiltersPending).toBe(false);

        consoleWarnSpy.mockRestore();
      });
    });
  });
});
