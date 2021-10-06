import Vuex from 'vuex';
import { cloneDeep } from 'lodash';
import { normalize } from 'normalizr';

import entitiesPlugin, { types } from '@/store/plugins/entities';
import { alarmSchema } from '@/store/schemas';
import { ENTITIES_TYPES } from '@/constants';

import { createVueInstance } from '@/unit/utils/vue';
import { fakeAlarms } from '@/unit/data/alarm';

const mockData = {
  alarms: fakeAlarms({ count: 1 }),
};

const storeConfig = {
  plugins: [entitiesPlugin],
};

describe('Entities plugin', () => {
  beforeAll(() => {
    const localVue = createVueInstance();

    localVue.use(Vuex);
  });

  it('Test', () => {
    const store = new Vuex.Store(cloneDeep(storeConfig));

    const { entities } = normalize(mockData.alarms, [alarmSchema]);

    store.commit(types.ENTITIES_UPDATE, entities);

    const alarm = store.getters['entities/getItem'](ENTITIES_TYPES.alarm, mockData.alarms[0]._id);

    expect(alarm).toEqual(mockData.alarms[0]);
  });
});
