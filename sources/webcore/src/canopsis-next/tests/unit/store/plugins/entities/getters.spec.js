import { createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import { cloneDeep } from 'lodash';
import { normalize } from 'normalizr';

import entitiesPlugin, { types } from '@/store/plugins/entities';
import { alarmSchema } from '@/store/schemas';
import { ENTITIES_TYPES } from '@/constants';

import { alarms } from './data';

const storeConfig = {
  plugins: [entitiesPlugin],
};

describe('Entities plugin', () => {
  beforeAll(() => {
    const localVue = createLocalVue();
    localVue.use(Vuex);
  });

  it('Test', () => {
    const store = new Vuex.Store(cloneDeep(storeConfig));

    const { entities } = normalize(alarms, [alarmSchema]);

    store.commit(types.ENTITIES_UPDATE, entities);

    const alarm = store.getters['entities/getItem'](ENTITIES_TYPES.alarm, alarms[0]._id);

    expect(alarm).toEqual(alarms[0]);
  });
});
