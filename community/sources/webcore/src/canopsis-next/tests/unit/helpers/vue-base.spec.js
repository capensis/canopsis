import { cloneDeep } from 'lodash';

import { mergeReactiveChangedProperties } from '@/helpers/vue-base';

describe('mergeReactiveChangedProperties', () => {
  it('shouldn\'t change link to empty object', () => {
    const original = {};
    const updated = {};

    const result = mergeReactiveChangedProperties(original, updated);

    expect(result === original).toBeTruthy();
    expect(result).toEqual(updated);
  });

  it('shouldn\'t change link to object', () => {
    const original = {
      stringProperty: 'property value',
      numberProperty: 2,
      arrayProperty: [1, 2, 3, 4],
      objectProperty: {
        subObjectProperty: 1,
      },
      deepProperty: [{
        object: {
          property: 1,
          testArray: [123],
        },
      }, {
        prop: 1,
      }],
    };
    const updated = {
      ...cloneDeep(original),
      // Removed one value
      arrayProperty: [1, 2, 3],
      objectProperty: {
        subObjectProperty: 1,
        // New property
        newObjectProperty: 1,
      },
      deepProperty: [{
        object: {
          // Changed properties
          property: 2,
          testArray: [],
        },
      }, {
        // Changed value
        prop: 3,
      }],
      newProperty: '123',
    };

    const result = mergeReactiveChangedProperties(original, updated);

    expect(result === original).toBeTruthy();
    expect(result.arrayProperty === original.arrayProperty).toBeTruthy();
    expect(result.deepProperty === original.deepProperty).toBeTruthy();
    expect(result.deepProperty[0] === original.deepProperty[0]).toBeTruthy();
    expect(result.deepProperty[0].object === original.deepProperty[0].object).toBeTruthy();
    expect(result.deepProperty[1] === original.deepProperty[1]).toBeTruthy();
    expect(result).toEqual(updated);
  });

  it('should return null when receive null as new object', () => {
    expect(mergeReactiveChangedProperties({}, null)).toEqual(null);
    expect(mergeReactiveChangedProperties([1], undefined)).toEqual(undefined);
  });
});
