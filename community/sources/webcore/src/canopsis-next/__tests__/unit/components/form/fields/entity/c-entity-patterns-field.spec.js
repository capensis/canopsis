import {
  BASIC_ENTITY_TYPES,
  ENTITY_PATTERN_FIELDS,
  ENTITY_TYPES,
  PATTERN_EXISTS_OPERATORS,
} from '@/constants';

import CEntityPatternsField from '@/components/forms/fields/entity/c-entity-patterns-field.vue';

describe('c-entity-patterns-field', () => {
  test('Multiple category values are displayed as closable chips', () => {
    const categories = [{ _id: 'category-id', name: 'Category' }];
    const categoryOptions = CEntityPatternsField.computed.categoryOptions.call({
      categories,
      categoriesPending: false,
    });

    const props = categoryOptions.valueField.props({ value: ['category-id'] });

    expect(props).toEqual(
      expect.objectContaining({
        items: categories,
        chips: true,
        multiple: true,
        closableChips: true,
      }),
    );
  });

  test('Connector values are limited to connector entities when entity types are provided', () => {
    const connectorOptions = CEntityPatternsField.computed.connectorOptions.call({
      entityTypes: [ENTITY_TYPES.service],
    });

    expect(connectorOptions.valueField.props.entityTypes).toEqual([BASIC_ENTITY_TYPES.connector]);
  });

  test('Pause fields are available with the correct operators', () => {
    const entityAttributes = CEntityPatternsField.computed.entityAttributes.call({
      $t: key => key,
      entitiesOptions: {},
      nameOptions: {},
      typeOptions: {},
      componentOptions: {},
      connectorOptions: {},
      infosOptions: {},
      componentInfosOptions: {},
      categoryOptions: {},
      impactLevelOptions: {},
      dateOptions: {},
      aliasesAttributes: [],
    });

    expect(entityAttributes).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          value: ENTITY_PATTERN_FIELDS.pause,
          options: { operators: PATTERN_EXISTS_OPERATORS },
        }),
        expect.objectContaining({ value: ENTITY_PATTERN_FIELDS.pauseAuthor }),
        expect.objectContaining({ value: ENTITY_PATTERN_FIELDS.pauseInitiator }),
      ]),
    );
  });
});
