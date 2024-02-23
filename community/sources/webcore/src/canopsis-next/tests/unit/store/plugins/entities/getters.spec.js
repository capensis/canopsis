import Vuex from 'vuex';
import { cloneDeep } from 'lodash';
import { normalize } from 'normalizr';

import { createVueInstance } from '@unit/utils/vue';
import { fakeAlarms } from '@unit/data/alarm';

import { ENTITIES_TYPES } from '@/constants';

import entitiesPlugin, { types } from '@/store/plugins/entities';
import { alarmSchema } from '@/store/schemas';

const mockData = {
  alarms: fakeAlarms(1),
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
