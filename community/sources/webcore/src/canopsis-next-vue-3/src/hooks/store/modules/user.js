import { storeToRefs } from 'pinia'

import { useUserStore } from '@/stores/modules/user'

export const useUser = () => {
  const userStore = useUserStore()
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(userStore)

  return {
    users: items,
    usersMeta: meta,
    usersPending: pending,
    usersFetchingParams: fetchingParams,

    fetchUsersList: userStore.fetchList,
    fetchUsersListWithPreviousParams: userStore.fetchListWithPreviousParams,
    fetchUsersListWithoutStore: userStore.fetchListWithoutStore,
    createUser: userStore.create,
    updateUser: userStore.update,
    removeUser: userStore.remove,
    updateCurrentUser: userStore.updateCurrentUser,
    updateCurrentUserTours: userStore.updateCurrentUserTours,
    bulkEnableUsers: userStore.bulkEnable,
    bulkDisableUsers: userStore.bulkDisable,
    bulkRemoveUsers: userStore.bulkRemove,
  }
}
