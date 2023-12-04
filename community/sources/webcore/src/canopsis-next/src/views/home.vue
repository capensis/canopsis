<template>
  <div v-if="!pendingDefaultView">
    <div id="brand">
      Canopsis Next
    </div>
  </div>
</template>

<script>
import { ROUTES_NAMES } from '@/constants';

import { getFirstRoleWithDefaultView } from '@/helpers/entities/user/entity';

import { authMixin } from '@/mixins/auth';
import { entitiesRoleMixin } from '@/mixins/entities/role';

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
      const { defaultview: roleDefaultView } = getFirstRoleWithDefaultView(this.currentUser) ?? {};

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
