<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-container
      v-alert(v-show="config.isDuplicating", type="info") {{ $t('modals.view.duplicate.infoMessage') }}
    v-card-text
      v-form(v-if="hasUpdateViewAccess")
        v-layout(wrap, justify-center)
          v-flex(xs11)
            v-text-field(
            :label="$t('common.name')",
            v-model="form.name",
            data-vv-name="name",
            v-validate="'required'",
            :error-messages="errors.collect('name')",
            )
            v-text-field(
            :label="$t('common.title')",
            v-model="form.title",
            data-vv-name="title",
            v-validate="'required'",
            :error-messages="errors.collect('title')",
            )
            v-text-field(
            :label="$t('common.description')",
            v-model="form.description",
            data-vv-name="description",
            )
            v-switch(v-model="form.enabled", :label="$t('common.enabled')")
        v-layout(wrap, justify-center)
          v-flex(xs11)
            v-combobox(
            v-model="form.tags",
            :label="$t('modals.view.fields.groupTags')",
            tags,
            clearable,
            multiple,
            append-icon,
            chips,
            deletable-chips,
            )
            v-combobox(
            ref="combobox",
            v-model="groupName",
            :items="groupNames",
            :label="$t('modals.view.fields.groupIds')",
            :search-input.sync="search"
            data-vv-name="group",
            v-validate="'required'",
            :error-messages="errors.collect('group')",
            @change="closeComboboxMenuOnChange()"
            )
              template(slot="no-data")
                v-list-tile
                  v-list-tile-content
                    v-list-tile-title(v-html="$t('modals.view.noData')")

            span {{ form.group_id }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(v-if="hasUpdateViewAccess", @click="submit") {{ $t('common.submit') }}
      v-btn.error(
      v-if="config.view && hasDeleteViewAccess && !config.isDuplicating",
      @click="remove"
      ) {{ $t('common.delete') }}
</template>

<script>
import find from 'lodash/find';
import omit from 'lodash/omit';

import { MODALS, USERS_RIGHTS_TYPES, USERS_RIGHTS_MASKS } from '@/constants';
import { generateView, generateViewRow, generateRight, generateRoleRightByChecksum } from '@/helpers/entities';
import uuid from '@/helpers/uid';
import authMixin from '@/mixins/auth';
import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesRightMixin from '@/mixins/entities/right';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';
import vuetifyComboboxMixin from '@/mixins/vuetify/combobox';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    authMixin,
    popupMixin,
    modalInnerMixin,
    entitiesViewMixin,
    entitiesRoleMixin,
    entitiesRightMixin,
    entitiesViewGroupMixin,
    rightsTechnicalViewMixin,
    vuetifyComboboxMixin,
  ],
  data() {
    return {
      search: '',
      groupName: '',
      form: {
        name: '',
        title: '',
        description: '',
        enabled: false,
        tags: [],
      },
    };
  },
  computed: {
    groupNames() {
      return this.groups.map(group => group.name);
    },

    title() {
      if (this.config.isDuplicating) {
        return `${this.$t('modals.view.duplicate.title')} - ${this.config.view.name}`;
      }

      if (this.config.view) {
        return this.$t('modals.view.edit.title');
      }

      return this.$t('modals.view.create.title');
    },

    hasUpdateViewAccess() {
      if (this.config.view) {
        return this.checkUpdateAccess(this.config.view._id) && this.hasUpdateAnyViewAccess;
      }

      return this.hasUpdateAnyViewAccess;
    },

    hasDeleteViewAccess() {
      if (this.config.view) {
        return this.checkDeleteAccess(this.config.view._id) && this.hasDeleteAnyViewAccess;
      }

      return this.hasDeleteAnyViewAccess;
    },
  },
  mounted() {
    const { view, isDuplicating } = this.config;

    if (view) {
      const group = find(this.groups, { _id: view.group_id });

      if (group) {
        this.groupName = group.name;
      }

      this.form = {
        name: isDuplicating ? '' : view.name,
        title: isDuplicating ? '' : view.title,
        description: view.description,
        enabled: view.enabled,
        tags: [...view.tags || []],
      };
    }
  },
  methods: {
    async createRightByViewId(viewId) {
      try {
        const checksum = USERS_RIGHTS_MASKS.read + USERS_RIGHTS_MASKS.update + USERS_RIGHTS_MASKS.delete;
        const role = await this.fetchRoleWithoutStore({ id: this.currentUser.role });

        const right = {
          ...generateRight(),

          _id: viewId,
          type: USERS_RIGHTS_TYPES.rw,
          desc: `Rights on view: ${viewId}`,
        };

        await this.createRight({ data: right });
        await this.createRole({
          data: {
            ...role,
            rights: {
              ...role.rights,

              [right._id]: generateRoleRightByChecksum(checksum),
            },
          },
        });

        return this.fetchCurrentUser();
      } catch (err) {
        this.addErrorPopup({ text: this.$t('modals.view.errors.rightCreating') });

        return Promise.resolve();
      }
    },

    async removeRightByViewId(viewId) {
      try {
        const { data: roles } = await this.fetchRolesListWithoutStore({ params: { limit: 10000 } });

        return Promise.all([
          this.removeRight({ id: viewId }),

          ...roles.map(role => this.createRole({
            data: {
              ...role,
              rights: omit(role.rights, [viewId]),
            },
          })),
        ]);
      } catch (err) {
        this.addErrorPopup({ text: this.$t('modals.view.errors.rightRemoving') });

        return Promise.resolve();
      }
    },

    remove() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeView({ id: this.config.view._id });
              await Promise.all([
                this.removeRightByViewId(this.config.view._id),
                this.fetchGroupsList(),
              ]);

              this.addSuccessPopup({ text: this.$t('modals.view.success') });
              this.hideModal();
            } catch (err) {
              this.addErrorPopup({ text: this.$t('modals.view.fail') });
            }
          },
        },
      });
    },

    async submit() {
      try {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          let group = find(this.groups, { name: this.groupName });

          if (!group) {
            group = await this.createGroup({ data: { name: this.groupName } });
          }

          /**
           * If we're creating a new view, or duplicating an existing one.
           * Generate a new view. Then copy rows and widgets if we're duplicating a view
           */
          if (!this.config.view || this.config.isDuplicating) {
            const data = {
              ...generateView(),
              ...this.form,
              group_id: group._id,
            };

            if (this.config.isDuplicating) {
              data.rows = this.config.view.rows.map(row => ({
                ...generateViewRow(),

                title: row.title,
                widgets: row.widgets.map(widget => ({ ...widget, _id: uuid(`widget_${widget.type}`) })),
              }));
            }

            const response = await this.createView({ data });
            await this.createRightByViewId(response._id);
            this.addSuccessPopup({ text: this.$t('modals.view.success.create') });
          } else {
            const data = {
              ...this.config.view,
              ...this.form,
              group_id: group._id,
            };

            await this.updateView({ id: this.config.view._id, data });
            this.addSuccessPopup({ text: this.$t('modals.view.success.edit') });
          }

          await this.fetchGroupsList();
          this.hideModal();
        }
      } catch (err) {
        /**
         * If we got a view in modal's config, and if we're not duplicating a view, that
         * means we're editing a view
        */
        if (!this.config.isDuplicating && this.config.view) {
          this.addErrorPopup({ text: this.$t('modals.view.fail.edit') });
        }
        this.addErrorPopup({ text: this.$t('modals.view.fail.create') });
        console.error(err.description);
      }
    },
  },
};
</script>
