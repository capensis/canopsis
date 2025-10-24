import Faker from 'faker';

import { widgetColumnsEntityInfoPropertyMixin } from '@/mixins/widget/columns/entity-info-property';

describe('widget-columns-entity-info-property', () => {
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
    {
      _id: Faker.datatype.string(),
      name: 'disk_usage',
      alias: '',
      description: Faker.datatype.string(),
    },
    {
      _id: Faker.datatype.string(),
      name: 'network_usage',
      description: Faker.datatype.string(),
    },
  ];

  let mockContext;

  beforeEach(() => {
    mockContext = {
      entityInfoProperties: mockEntityInfoProperties,
      fetchEntityInfoPropertiesList: jest.fn(),
    };
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('mounted', () => {
    test('Should fetch entity info properties list without pagination on mount', () => {
      widgetColumnsEntityInfoPropertyMixin.mounted.call(mockContext);

      expect(mockContext.fetchEntityInfoPropertiesList).toHaveBeenCalledWith({
        params: { paginate: false },
      });
    });
  });

  describe('methods', () => {
    describe('findAliasByColumnValue', () => {
      test('Should return alias when column value matches entity info property with alias', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.infos.cpu_usage.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBe('CPU Usage');
      });

      test('Should return alias when column value matches entity info property without prefix', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'infos.memory_usage.value';

        const result = findAliasByColumnValue.call(mockContext, columnValue);

        expect(result).toBe('Memory Usage');
      });

      test('Should return undefined when column value matches entity info property without alias', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.infos.disk_usage.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBeUndefined();
      });

      test('Should return undefined when column value matches entity info property with empty alias', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.infos.network_usage.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBeUndefined();
      });

      test('Should return undefined when column value does not match any entity info property', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.infos.nonexistent.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBeUndefined();
      });

      test('Should return undefined when column value format is incorrect', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.cpu_usage';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBeUndefined();
      });

      test('Should return undefined when entityInfoProperties is empty', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const context = {
          ...mockContext,
          entityInfoProperties: [],
        };
        const columnValue = 'entity.infos.cpu_usage.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(context, columnValue, prefix);

        expect(result).toBeUndefined();
      });

      test('Should handle column value with different prefix', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'alarm.entity.infos.cpu_usage.value';
        const prefix = 'alarm.entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBe('CPU Usage');
      });

      test('Should handle column value with null prefix', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'infos.cpu_usage.value';
        const prefix = null;

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBe('CPU Usage');
      });

      test('Should handle column value with undefined prefix', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'infos.cpu_usage.value';
        const prefix = undefined;

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBe('CPU Usage');
      });

      test('Should handle column value with empty string prefix', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'infos.cpu_usage.value';
        const prefix = '';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBe('CPU Usage');
      });

      test('Should correctly filter Boolean values in join', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.infos.cpu_usage.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBe('CPU Usage');
      });

      test('Should not match partial property names', () => {
        const { findAliasByColumnValue } = widgetColumnsEntityInfoPropertyMixin.methods;
        const columnValue = 'entity.infos.cpu.value';
        const prefix = 'entity';

        const result = findAliasByColumnValue.call(mockContext, columnValue, prefix);

        expect(result).toBeUndefined();
      });
    });
  });
});
