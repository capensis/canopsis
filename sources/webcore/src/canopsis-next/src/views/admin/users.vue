<template lang="pug">
  div
    the-page-header {{ $t('common.users') }}
    users
    fab-buttons(
      :has-access="hasCreateAnyUserAccess",
      @refresh="fetchList",
      @create="showCreateUserModal"
    )
      span {{ $t('modals.createUser.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import { prepareUserByData } from '@/helpers/entities';

import entitiesUserMixin from '@/mixins/entities/user';
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';

import Users from '@/components/other/users/users.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export default {
  components: {
    FabButtons,
    Users,
  },
  mixins: [entitiesUserMixin, rightsTechnicalUserMixin],
  methods: {
    fetchList() {
      return this.fetchUsersListWithPreviousParams();
    },

    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await this.createUserWithPopup({ data: prepareUserByData(data) });

            await this.fetchList();
          },
        },
      });
    },
  },
};
</script>
