<template lang="pug">
  v-menu(bottom, offset-y, offset-x)
    v-btn.white--text(slot="activator", flat) {{ userName }}
    v-list.py-0
      v-list-tile(@click="showEditUserModal")
        v-list-tile-avatar
          v-icon(color="black") person
        v-list-tile-title.text-uppercase.body-2 {{ $t('users.seeProfile') }}
      v-list-tile(:to="profilePatternsLink", active-class="")
        v-list-tile-avatar
          v-icon(color="black") filter_list
        v-list-tile-title.text-uppercase.body-2 {{ $t('pattern.patterns') }}
      v-list-tile.logout-btn(@click="logoutHandler")
        v-list-tile-avatar
          v-icon(color="error") exit_to_app
        v-list-tile-title.text-uppercase.error--text.body-2 {{ $t('common.logout') }}
</template>

<script>
import { MODALS, ROUTES_NAMES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesUserMixin } from '@/mixins/entities/user';

export default {
  mixins: [authMixin, entitiesUserMixin],
  computed: {
    profilePatternsLink() {
      return { name: ROUTES_NAMES.profilePatterns };
    },

    userName() {
      return this.currentUser.name || this.currentUser._id;
    },
  },
  methods: {
    showEditUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          title: this.$t('common.profile'),
          user: this.currentUser,
          onlyUserPrefs: true,
          action: async (data) => {
            await this.updateCurrentUser({ data });

            this.$popups.success({ text: this.$t('success.default', data.ui_language) });

            await this.fetchCurrentUser();
          },
        },
      });
    },

    logoutHandler() {
      return this.logout({ redirectTo: { name: ROUTES_NAMES.login } });
    },
  },
};
</script>

<style lang="scss" scoped>
$btnErrorColor: rgba(255, 82, 82, .1);

.logout-btn {
  &:hover, &:active {
    background: $btnErrorColor;
  }
}
</style>
