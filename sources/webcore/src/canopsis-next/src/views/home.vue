<template lang="pug">
  div(v-if="!pendingDefaultView")
    div#brand Canopsis Next
</template>

<script>
import authMixin from '@/mixins/auth';
import entitiesRoleMixin from '@/mixins/entities/role';

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
      const { defaultview: defaultViewId } = this.currentUser;

      if (!defaultViewId) {
        await this.redirectToRoleDefaultView();
      } else if (!this.checkReadAccess(defaultViewId)) {
        this.addRedirectInfoPopup(this.$t('home.popups.info.noAccessToDefaultView'));

        await this.redirectToRoleDefaultView();
      } else {
        this.$router.push({ name: 'view', params: { id: defaultViewId } });
      }
    },

    async redirectToRoleDefaultView() {
      const { defaultview: roleDefaultViewId } = await this.fetchRoleWithoutStore({ id: this.currentUser.role });

      if (!roleDefaultViewId) {
        this.addRedirectInfoPopup(this.$t('home.popups.info.notSelectedRoleDefaultView'));
      } else if (!this.checkReadAccess(roleDefaultViewId)) {
        this.addRedirectInfoPopup(this.$t('home.popups.info.noAccessToRoleDefaultView'));
      } else {
        this.$router.push({ name: 'view', params: { id: roleDefaultViewId } });
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
