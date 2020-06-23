<template lang="pug">
  v-container.admin-rights
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.rights') }}
    div.progress-wrapper
      v-fade-transition
        v-layout.white.progress(v-show="pending", column)
          v-progress-circular(indeterminate, color="primary")
      v-tabs(fixed-tabs, slider-color="primary")
        template(v-for="(rights, groupKey) in groupedRights")
          v-tab(:key="`tab-${groupKey}`") {{ groupKey }}
          v-tab-item.white(:key="`tab-item-${groupKey}`")
            div.pa-3(v-if="hasReadAnyRoleAccess")
              rights-groups-table(
                v-if="groupKey === 'business' || groupKey === 'technical'",
                :groups="rights",
                :roles="roles",
                :changedRoles="changedRoles",
                @change="changeCheckboxValue"
              )
              rights-table(
                v-else,
                :rights="rights",
                :roles="roles",
                :changedRoles="changedRoles",
                @change="changeCheckboxValue"
              )
    v-layout.submit-button.mt-3(v-show="hasUpdateAnyActionAccess && hasChanges")
      v-btn.primary.ml-3(@click="submit") {{ $t('common.submit') }}
      v-btn(@click="cancel") {{ $t('common.cancel') }}
    .fab(v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess || hasCreateAnyActionAccess")
      v-layout(column)
        refresh-btn(@click="fetchRightsList")
        v-speed-dial(
          v-model="fab",
          direction="left",
          transition="slide-y-reverse-transition"
        )
          v-btn(slot="activator", color="primary", fab, v-model="fab")
            v-icon add
            v-icon close
          v-tooltip(v-if="hasCreateAnyUserAccess", top)
            v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateUserModal")
              v-icon people
            span {{ $t('modals.createUser.title') }}
          v-tooltip(v-if="hasCreateAnyRoleAccess", top)
            v-btn(slot="activator", fab, dark, small, color="deep-purple ", @click.stop="showCreateRoleModal")
              v-icon supervised_user_circle
            span {{ $t('modals.createRole.title') }}
          v-tooltip(v-if="hasCreateAnyActionAccess", top)
            v-btn(slot="activator", fab, dark, small, color="teal", @click.stop="showCreateRightModal")
              v-icon verified_user
            span {{ $t('modals.createRight.title') }}
</template>

<script>
import { get, omit, isEmpty, isUndefined, transform } from 'lodash';
import flatten from 'flat';

import {
  MODALS,
  USERS_RIGHTS,
  NOT_COMPLETED_USER_RIGHTS_KEYS,
} from '@/constants';
import {
  prepareUserByData,
  generateRoleRightByChecksum,
} from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import entitiesRightMixin from '@/mixins/entities/right';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesUserMixin from '@/mixins/entities/user';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import entitiesPlaylistMixin from '@/mixins/entities/playlist';
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';
import rightsTechnicalActionMixin from '@/mixins/rights/technical/action';

import RightsGroupsTable from '@/components/other/right/admin/rights-groups-table.vue';
import RightsTable from '@/components/other/right/admin/rights-table.vue';
import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';

export default {
  components: {
    RightsGroupsTable,
    RightsTable,
    RefreshBtn,
  },
  mixins: [
    authMixin,
    entitiesRightMixin,
    entitiesRoleMixin,
    entitiesUserMixin,
    entitiesViewGroupMixin,
    entitiesPlaylistMixin,
    rightsTechnicalUserMixin,
    rightsTechnicalRoleMixin,
    rightsTechnicalActionMixin,
  ],
  data() {
    return {
      fab: false,
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
    this.fetchRightsList();
  },
  methods: {
    async fetchRightsList() {
      this.pending = true;

      await this.fetchList();

      this.pending = false;
    },

    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: data => this.createUser({ data: prepareUserByData(data) }),
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
          action: this.fetchRightsList,
        },
      });
    },

    clearChangedRoles() {
      this.changedRoles = {};
    },

    cancel() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.clearChangedRoles,
        },
      });
    },

    submit() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.updateRoles,
        },
      });
    },

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
     * @param {number} rightType
     */
    changeCheckboxValue(value, role, right, rightType) {
      const currentCheckSum = get(role, ['rights', right._id, 'checksum'], 0);
      const factor = value ? 1 : -1;

      /**
       * If we don't have changes for role
       */
      if (!this.changedRoles[role._id]) {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles, role._id, { [right._id]: nextCheckSum });

        /**
         * If we have changes for role but we don't have changes for right
         */
      } else if (isUndefined(this.changedRoles[role._id][right._id])) {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles[role._id], right._id, nextCheckSum);

        /**
         * If we have changes for role and for right
         */
      } else {
        const nextCheckSum = !rightType ?
          Number(value) : this.changedRoles[role._id][right._id] + (factor * rightType);

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
      const [{ data: rights }] = await Promise.all([
        this.fetchRightsListWithoutStore({ params: { limit: 10000 } }),
        this.fetchRolesList({ params: { limit: 10000 } }),
      ]);

      const allViews = this.groups.reduce((acc, { views }) => acc.concat(views), []);
      const allRightsIds = flatten(USERS_RIGHTS);
      const allBusinessRightsIds = flatten(USERS_RIGHTS.business);
      const { exploitation: exploitationTechnicalRights, ...adminTechnicalRights } = USERS_RIGHTS.technical;
      const adminTechnicalRightsValues = Object.values(adminTechnicalRights);
      const exploitationTechnicalRightsValues = Object.values(exploitationTechnicalRights);

      const groupedRights = rights.reduce((acc, right) => {
        const rightId = String(right._id, '\'');
        const view = allViews.find(({ _id }) => _id === rightId);
        const playlist = this.playlists.find(({ _id }) => _id === rightId);

        if (view) {
          acc.view.push({
            ...right,

            desc: right.desc.replace(view._id, view.name),
          });
        } else if (playlist) {
          acc.playlist.push({
            ...right,

            desc: right.desc.replace(playlist._id, playlist.name),
          });
        } else if (adminTechnicalRightsValues.indexOf(rightId) !== -1) {
          acc.technical.admin.push(right);
        } else if (exploitationTechnicalRightsValues.indexOf(rightId) !== -1) {
          acc.technical.exploitation.push(right);
        } else if (
          Object.values(allBusinessRightsIds).indexOf(rightId) !== -1 ||
          NOT_COMPLETED_USER_RIGHTS_KEYS.some(userRightKey => rightId.startsWith(allRightsIds[userRightKey]))
        ) {
          const [parentKey] = right._id.split('_');

          if (!acc.business[parentKey]) {
            acc.business[parentKey] = [right];
          } else {
            acc.business[parentKey].push(right);
          }
        }

        return acc;
      }, {
        business: {},
        view: [],
        playlist: [],
        technical: {
          admin: [],
          exploitation: [],
        },
      });

      groupedRights.business = Object.entries(groupedRights.business).map(([key, value]) => ({ key, rights: value }));
      groupedRights.technical = Object.entries(groupedRights.technical).map(([key, value]) => ({ key, rights: value }));
      groupedRights.view = [...groupedRights.view, ...groupedRights.playlist];

      this.groupedRights = omit(groupedRights, ['playlist']);
    },
  },
};
</script>

<style lang="scss" scoped>
  $firstTdWidth: 350px;

  .submit-button {
    position: sticky;
    bottom: 10px;
  }

  .admin-rights {
    & /deep/ {
      .v-table__overflow {
        overflow: visible;
        td {
          padding: 8px 10px;

          &:first-child {
            width: $firstTdWidth;
          }
        }

        th {
          transition: none;
          position: sticky;
          top: 48px;
          background: white;
          z-index: 1;
        }

        .v-datatable__expand-content .v-table {
          background: #f3f3f3;
        }
      }

      .v-expansion-panel__body {
        overflow: auto;
      }

      .v-window__container--is-active th {
        position: relative;
        top: 0;
      }

      .expand-rights-table {
        .v-table__overflow {
          tr td {
            &:first-child {
              padding-left: 36px;
            }
          }

          thead tr {
            height: 0;
            visibility: hidden;

            th {
              position: relative;
              height: 0;
              line-height: 0;
              padding: 0 24px;
            }
          }
        }
      }
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

  .progress-wrapper {
    position: relative;
  }
</style>
