import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

import { userToForm, formToUserRequest } from '@/helpers/entities/user/form';

import { useAuthStore } from './auth';

export const useUserStore = createCRUDStore({
  id: 'user',
  route: API_ROUTES.users,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    fetchTopByActiveConnectsWithoutStore() {
      return request.get(API_ROUTES.topUserByActiveConnects);
    },

    remove({ id } = {}) {
      return request.delete(`${API_ROUTES.users}/${encodeURIComponent(id)}`);
    },

    async updateCurrentUser({ data } = {}) {
      const authStore = useAuthStore();
      const newCurrentUser = await request.put(API_ROUTES.currentUser, data);

      authStore.setCurrentUser(newCurrentUser);
      await authStore.afterFetchCurrentUser(newCurrentUser);

      return newCurrentUser;
    },

    updateCurrentUserTours({ data } = {}) {
      const authStore = useAuthStore();
      const userForm = userToForm(authStore.currentUser);

      return request
        .put(
          API_ROUTES.currentUser,
          formToUserRequest({
            ...userForm,

            ui_tours: {
              ...userForm.ui_tours,
              ...data,
            },
          }),
        )
        .then(async (newCurrentUser) => {
          authStore.setCurrentUser(newCurrentUser);
          await authStore.afterFetchCurrentUser(newCurrentUser);

          return newCurrentUser;
        });
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkUsersEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkUsersDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkUsers, { data });
    },
  },
});
