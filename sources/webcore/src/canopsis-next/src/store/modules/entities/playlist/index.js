import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.playlist,
  entityType: ENTITIES_TYPES.playlist,
}, {
  actions: {
    fetchItemWithoutStore() {
      return {
        _id: 'asd',
        name: 'Playlist #1',
        fullscreen: true,
        interval: {
          value: 10,
          unit: 'm',
        },
        tabs: [
          '875df4c2-027b-4549-8add-e20ed7ff7d4f', // Alarm default
          'view-tab_5a339b3a-0611-4d4c-b307-dc1b92aeb27d', // Meteo technic
          'view-tab_c02ae48e-7f0a-4ba4-9215-ba5662e1550c', // Meteo correct
        ],
      };
    },
  },
});
