import Faker from 'faker';

import { testsEntityModule } from '@unit/utils/store';

import { API_ROUTES } from '@/config';

import patternModule from '@/store/modules/entities/pattern';

describe('Entities pattern module', () => {
  const patternIds = [Faker.datatype.string(), Faker.datatype.string()];
  const patterns = patternIds.map(id => ({
    _id: id,
    name: `name-${id}`,
  }));

  const { axiosMockAdapter } = testsEntityModule({
    route: API_ROUTES.pattern.list,
    module: patternModule,
    entities: patterns,
    entityIds: patternIds,
  });

  test('Fetch list without store. Action: fetchListWithoutStore', async () => {
    const params = {
      param: Faker.datatype.string(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { params })
      .reply(200, patterns);

    const result = await patternModule.actions.fetchListWithoutStore({}, { params });

    expect(result).toEqual(patterns);
  });

  test('Bulk remove. Action: bulkRemove', async () => {
    const response = patterns.map(() => ({
      status: 200,
    }));
    axiosMockAdapter
      .onDelete(API_ROUTES.pattern.bulkList, patterns)
      .reply(200, response);

    const result = await patternModule.actions.bulkRemove({}, { data: patterns });

    expect(result).toEqual(response);
  });
});
