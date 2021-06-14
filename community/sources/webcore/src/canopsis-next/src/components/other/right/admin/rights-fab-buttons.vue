<template lang="pug">
  c-fab-expand-btn(
    v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess || hasCreateAnyActionAccess",
    @refresh="refresh"
  )
    v-tooltip(v-if="hasCreateAnyUserAccess", top)
      v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateUserModal")
        v-icon people
      span {{ $t('modals.createUser.title') }}
    v-tooltip(v-if="hasCreateAnyRoleAccess", top)
      v-btn(slot="activator", fab, dark, small, color="deep-purple", @click.stop="showCreateRoleModal")
        v-icon supervised_user_circle
      span {{ $t('modals.createRole.title') }}
    v-tooltip(v-if="hasCreateAnyActionAccess", top)
      v-btn(slot="activator", fab, dark, small, color="teal", @click.stop="showCreateRightModal")
        v-icon verified_user
      span {{ $t('modals.createRight.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import { prepareUserByData } from '@/helpers/entities';

import entitiesUserMixin from '@/mixins/entities/user';
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';
import rightsTechnicalActionMixin from '@/mixins/rights/technical/action';

export default {
  mixins: [
    entitiesUserMixin,
    rightsTechnicalUserMixin,
    rightsTechnicalRoleMixin,
    rightsTechnicalActionMixin,
  ],
  methods: {
    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: data => this.createUserWithPopup({ data: prepareUserByData(data) }),
        },
      });
    },

    showCreateRoleModal() {
      this.$modals.show({
        name: MODALS.createRole,
      });
    },

    showCreateRightModal() {
      this.$modals.show({
        name: MODALS.createRight,
        config: {
          action: this.refresh,
        },
      });
    },

    refresh() {
      this.$emit('refresh');
    },
  },
};
</script>
