<template lang="pug">
  v-form(data-test="createViewModal", @submit.prevent="submit")
    modal-wrapper(data-test="createViewModal")
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-container
          v-alert(v-show="config.isDuplicating", type="info") {{ $t('modals.view.duplicate.infoMessage') }}
        view-form(
          v-if="hasUpdateViewAccess",
          v-model="form",
          :groupName.sync="groupName",
          :groups="groups"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          v-if="hasUpdateViewAccess",
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="viewSubmitButton"
        ) {{ $t('common.submit') }}
        v-btn.error(
          v-if="config.view && hasDeleteViewAccess && !config.isDuplicating",
          :disabled="submitting",
          data-test="viewDeleteButton",
          @click="remove"
        ) {{ $t('common.delete') }}
</template>

<script>
import { find } from 'lodash';

import { MODALS, DEFAULT_PERIODIC_REFRESH } from '@/constants';
import {
  generateView,
  generateCopyOfView,
  getViewsWidgetsIdsMappings,
} from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesViewRightsMixin from '@/mixins/entities/view/rights';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesRightMixin from '@/mixins/entities/right';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

import ViewForm from '@/components/other/view/view-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ViewForm, ModalWrapper },
  mixins: [
    authMixin,
    modalInnerMixin,
    submittableMixin(),
    entitiesViewMixin,
    entitiesRoleMixin,
    entitiesRightMixin,
    entitiesViewGroupMixin,
    entitiesViewRightsMixin,
    entitiesUserPreferenceMixin,
    rightsTechnicalViewMixin,
  ],
  data() {
    return {
      groupName: '',
      form: {
        name: '',
        title: '',
        description: '',
        enabled: false,
        tags: [],
        periodicRefresh: DEFAULT_PERIODIC_REFRESH,
      },
    };
  },
  computed: {
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
        periodicRefresh: view.periodicRefresh
          ? {
            interval: view.periodicRefresh.interval || view.periodicRefresh.value,
            unit: view.periodicRefresh.unit || DEFAULT_PERIODIC_REFRESH.unit,
            enabled: view.periodicRefresh.enabled,
          }
          : DEFAULT_PERIODIC_REFRESH,
      };
    }
  },
  methods: {
    remove() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeView({ id: this.config.view._id });
              await Promise.all([
                this.removeRightByViewId(this.config.view._id),
                this.fetchGroupsList(),
              ]);

              this.$popups.success({ text: this.$t('modals.view.success.delete') });
              this.$modals.hide();
            } catch (err) {
              this.$popups.error({ text: this.$t('modals.view.fail.delete') });
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
           * Generate a new view. Then copy tabs, rows and widgets if we're duplicating a view
           */
          if (!this.config.view || this.config.isDuplicating) {
            const data = {
              ...generateView(),
              ...this.form,
              group_id: group._id,
            };

            if (this.config.isDuplicating) {
              const { tabs } = generateCopyOfView(this.config.view);

              data.tabs = tabs;

              const widgetsIdsMappings = getViewsWidgetsIdsMappings(this.config.view, data);

              await this.copyUserPreferencesByWidgetsIdsMappings(widgetsIdsMappings);
            }

            const response = await this.createView({ data });
            await this.createRightByViewId(response._id);
            this.$popups.success({ text: this.$t('modals.view.success.create') });
          } else {
            const data = {
              ...this.config.view,
              ...this.form,
              group_id: group._id,
            };

            await this.updateView({ id: this.config.view._id, data });
            this.$popups.success({ text: this.$t('modals.view.success.edit') });
          }

          await this.fetchGroupsList();
          this.$modals.hide();
        }
      } catch (err) {
        /**
         * If we got a view in modal's config, and if we're not duplicating a view, that
         * means we're editing a view
         */
        if (!this.config.isDuplicating && this.config.view) {
          this.$popups.error({ text: this.$t('modals.view.fail.edit') });
        }
        this.$popups.error({ text: this.$t('modals.view.fail.create') });
        console.error(err.description);
      }
    },
  },
};
</script>
