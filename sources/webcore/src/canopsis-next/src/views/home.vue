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
  async created() {
    let defaultViewId = this.currentUser.defaultview;

    if (!defaultViewId) {
      const role = await this.fetchRoleWithoutStore({ id: this.currentUser.role });

      defaultViewId = role.defaultview;
    }

    if (defaultViewId) {
      this.$router.push({ name: 'view', params: { id: defaultViewId } });
    }

    this.pendingDefaultView = false;
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
