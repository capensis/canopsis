<template lang="pug">
  v-container.admin-rights
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.rights') }}
    div.position-relative
      v-fade-transition
        v-layout.white.progress(v-show="pending", column)
          v-progress-circular(indeterminate, color="primary")
      v-tabs(v-if="hasReadAnyRoleAccess", fixed-tabs, slider-color="primary")
        template(v-for="(rights, groupKey) in groupedRights")
          v-tab(:key="`tab-${groupKey}`") {{ groupKey }}
          v-tab-item.white(:key="`tab-item-${groupKey}`")
            rights-table-wrapper(
              :rights="rights",
              :roles="roles",
              :changedRoles="changedRoles",
              :disabled="!hasUpdateAnyActionAccess",
              @change="changeCheckboxValue"
            )
    v-layout.submit-button.mt-3(v-show="hasChanges")
      v-btn.primary.ml-3(@click="submit") {{ $t('common.submit') }}
      v-btn(@click="cancel") {{ $t('common.cancel') }}
    rights-fab-buttons(@refresh="fetchList")
</template>

<script>
import { get, isEmpty, isUndefined, transform } from 'lodash';

import { MODALS } from '@/constants';

import { getGroupedRights } from '@/helpers/right';
import { generateRoleRightByChecksum } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import entitiesRightMixin from '@/mixins/entities/right';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import entitiesPlaylistMixin from '@/mixins/entities/playlist';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';
import rightsTechnicalActionMixin from '@/mixins/rights/technical/action';

import RightsTableWrapper from '@/components/other/right/admin/rights-table-wrapper.vue';
import RightsFabButtons from '@/components/other/right/admin/rights-fab-buttons.vue';

export default {
  components: {
    RightsTableWrapper,
    RightsFabButtons,
  },
  mixins: [
    authMixin,
    entitiesRightMixin,
    entitiesRoleMixin,
    entitiesViewGroupMixin,
    entitiesPlaylistMixin,
    rightsTechnicalRoleMixin,
    rightsTechnicalActionMixin,
  ],
  data() {
    return {
      pending: false,
      groupedRights: { business: [], view: [], technical: [] },
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return !isEmpty(this.changedRoles);
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    /**
     * Clear changed roles
     *
     * @returns void
     */
    clearChangedRoles() {
      this.changedRoles = {};
    },

    /**
     * Show the confirmation modal window for clearing a changedRoles
     *
     * @returns void
     */
    cancel() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.clearChangedRoles,
        },
      });
    },

    /**
     * Show the confirmation modal window for submitting a changedRoles
     *
     * @returns void
     */
    submit() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.updateRoles,
        },
      });
    },

    /**
     * Send request for update changedRoles and fetchCurrentUser if it needed
     *
     * @returns void
     */
    async updateRoles() {
      try {
        this.pending = true;

        await Promise.all(Object.keys(this.changedRoles).map((roleId) => {
          const changedRoleRights = this.changedRoles[roleId];
          const role = this.getRoleById(roleId);

          const newRights = transform(
            changedRoleRights,
            (acc, value, key) => acc[key] = generateRoleRightByChecksum(value),
          );

          return this.createRole({ data: { ...role, rights: { ...role.rights, ...newRights } } });
        }));

        /**
         * If current user role changed
         */
        if (this.changedRoles[this.currentUser.role]) {
          await this.fetchCurrentUser();
        }

        this.$popups.success({ text: this.$t('success.default') });
        this.clearChangedRoles();

        this.pending = false;
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {Object} role
     * @param {Object} right
     * @param {number} mask
     * @returns void
     */
    changeCheckboxValue(value, role, right, mask) {
      const currentCheckSum = get(role, ['rights', right._id, 'checksum'], 0);
      const factor = value ? 1 : -1;

      /**
       * If we don't have changes for role
       */
      if (!this.changedRoles[role._id]) {
        const nextCheckSum = !mask ?
          Number(value) : currentCheckSum + (factor * mask);

        this.$set(this.changedRoles, role._id, { [right._id]: nextCheckSum });

        /**
         * If we have changes for role but we don't have changes for right
         */
      } else if (isUndefined(this.changedRoles[role._id][right._id])) {
        const nextCheckSum = !mask ?
          Number(value) : currentCheckSum + (factor * mask);

        this.$set(this.changedRoles[role._id], right._id, nextCheckSum);

        /**
         * If we have changes for role and for right
         */
      } else {
        const nextCheckSum = !mask ?
          Number(value) : this.changedRoles[role._id][right._id] + (factor * mask);

        if (currentCheckSum === nextCheckSum) {
          if (Object.keys(this.changedRoles[role._id]).length === 1) {
            this.$delete(this.changedRoles, role._id);
          } else {
            this.$delete(this.changedRoles[role._id], right._id);
          }
        } else {
          this.$set(this.changedRoles[role._id], right._id, nextCheckSum);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      }
    },

    /**
     * Fetch rights and roles lists
     *
     * @returns void
     */
    async fetchList() {
      this.pending = true;

      const [{ data: rights }] = await Promise.all([
        this.fetchRightsListWithoutStore({ params: { limit: 10000 } }),
        this.fetchRolesList({ params: { limit: 10000 } }),
      ]);

      const allViews = this.groups.reduce((acc, { views }) => acc.concat(views), []);

      this.groupedRights = getGroupedRights(rights, allViews, this.playlists);
      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
  .submit-button {
    position: sticky;
    bottom: 10px;
  }

  .admin-rights /deep/ {
    .v-window__container--is-active th {
      position: relative;
      top: 0;
    }
  }

  .progress {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    opacity: .4;
    z-index: 1;

    & /deep/ .v-progress-circular {
      top: 50%;
      left: 50%;
      margin-top: -16px;
      margin-left: -16px;
    }
  }
</style>
