<template lang="pug">
  div(v-if="!pendingDefaultView")
    div#brand Canopsis Next
</template>

<script>
import { authMixin } from '@/mixins/auth';
import { entitiesRoleMixin } from '@/mixins/entities/role';

import { ROUTES_NAMES } from '@/constants';

export default {
  mixins: [authMixin, entitiesRoleMixin],
  data() {
    return {
      pendingDefaultView: true,
    };
  },
  async mounted() {
    await this.redirectToDefaultView();

    this.pendingDefaultView = false;
  },
  methods: {
    async redirectToDefaultView() {
      const { defaultview: userDefaultView } = this.currentUser;

      if (!userDefaultView) {
        await this.redirectToRoleDefaultView();
      } else if (!this.checkReadAccess(userDefaultView._id)) {
        this.addRedirectInfoPopup(this.$t('home.popups.info.noAccessToDefaultView'));

        await this.redirectToRoleDefaultView();
      } else {
        this.$router.push({ name: ROUTES_NAMES.view, params: { id: userDefaultView._id } });
      }
    },

    async redirectToRoleDefaultView() {
      const { defaultview: roleDefaultView } = this.currentUser.role;

      if (!roleDefaultView) {
        this.addRedirectInfoPopup(this.$t('home.popups.info.notSelectedRoleDefaultView'));
      } else if (!this.checkReadAccess(roleDefaultView._id)) {
        this.addRedirectInfoPopup(this.$t('home.popups.info.noAccessToRoleDefaultView'));
      } else {
        this.$router.push({ name: ROUTES_NAMES.view, params: { id: roleDefaultView._id } });
      }
    },

    addRedirectInfoPopup(text) {
      return this.$popups.info({ text, autoClose: 10000 });
    },
  },
};
</script>

<style lang="scss" scoped>
  #brand {
    text-align: center;
    position: relative;
    top: 25%;
    max-width: 50%;
    max-height: 5em;
    margin: auto;
    font-weight: bold;
    font-size: 2em;
  }
</style>
