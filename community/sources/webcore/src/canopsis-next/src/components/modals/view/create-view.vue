<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-container(v-show="isDuplicating")
          v-alert(type="info") {{ $t('modals.view.duplicate.infoMessage') }}
        view-form(
          v-model="form",
          :groups="groups"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          v-if="hasUpdateViewAccess",
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
        v-btn.error(
          v-if="config.view && hasDeleteViewAccess && !isDuplicating",
          :disabled="submitting",
          @click="remove"
        ) {{ $t('common.delete') }}
</template>

<script>
import { find, isString } from 'lodash';

import { MODALS, ROUTES_NAMES } from '@/constants';

import { generateCopyOfViewTab, getViewsWidgetsIdsMappings } from '@/helpers/entities';
import { viewToForm, viewToRequest } from '@/helpers/forms/view';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import entitiesViewMixin from '@/mixins/entities/view';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

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
    modalInnerMixin,
    entitiesViewMixin,
    entitiesViewGroupMixin,
    entitiesUserPreferenceMixin,
    permissionsTechnicalViewMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: viewToForm(this.modal.config.view),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.view.create.title');
    },

    view() {
      return this.config.view;
    },

    isDuplicating() {
      return this.view && !this.view._id;
    },

    isEditing() {
      return this.view && this.view._id;
    },

    hasUpdateViewAccess() {
      if (this.view && !this.isDuplicating) {
        return this.checkUpdateAccess(this.view._id) && this.hasUpdateAnyViewAccess;
      }

      return this.hasUpdateAnyViewAccess;
    },

    hasDeleteViewAccess() {
      if (this.view && !this.isDuplicating) {
        return this.checkDeleteAccess(this.view._id) && this.hasDeleteAnyViewAccess;
      }

      return this.hasDeleteAnyViewAccess;
    },
  },
  mounted() {
    this.fetchAllGroupsListWithViewsWithCurrentUser();
  },
  methods: {
    /**
     * Redirect to home page if we surfing on this view at the moment
     */
    redirectToHomeIfCurrentRoute() {
      const { name, params = {} } = this.$route;

      if (name === ROUTES_NAMES.view && params.id === this.view._id) {
        this.$router.push({ name: ROUTES_NAMES.home });
      }
    },

    /**
     * Remove view
     */
    remove() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeViewWithPopup({ id: this.view._id });
              await this.fetchAllGroupsListWithViewsWithCurrentUser();

              this.redirectToHomeIfCurrentRoute();

              this.$modals.hide();
            } catch (err) {
              this.$popups.error({ text: this.$t('modals.view.fail.delete') });
            }
          },
        },
      });
    },

    /**
     * Create view group with special title
     *
     * @param {string} title
     * @return {Promise<ViewGroup>}
     */
    createGroupWithSpecialTitle(title) {
      return this.createGroup({
        data: { title },
      });
    },

    /**
     * Try to find view group by title or create a new one with special title
     *
     * @param {string} title
     * @return {ViewGroup | Promise<ViewGroup>}
     */
    prepareGroup(title) {
      return find(this.groups, { title })
        || this.createGroupWithSpecialTitle(this.form.group);
    },

    /**
     * Convert view form to request object with group and user preference creation if needed
     *
     * @return {Promise}
     */
    async formToRequest() {
      const group = isString(this.form.group)
        ? await this.prepareGroup(this.form.group)
        : this.form.group;

      const data = viewToRequest({
        ...this.view,
        ...this.form,

        group,
      });

      /**
       * If we're creating a new view, or duplicating an existing one.
       * Generate a new view. Then copy tabs and widgets if we're duplicating a view
       */
      if (this.isDuplicating) {
        data.tabs = data.tabs.map(generateCopyOfViewTab);

        const widgetsIdsMappings = getViewsWidgetsIdsMappings(this.view, data);

        await this.copyUserPreferencesByWidgetsIdsMappings(widgetsIdsMappings);
      }

      return data;
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          const data = await this.formToRequest();

          if (this.isEditing) {
            await this.updateViewWithPopup({ id: this.view._id, data });
          } else {
            await this.createViewWithPopup({ data });
          }

          await this.fetchAllGroupsListWithViewsWithCurrentUser();

          this.$modals.hide();
        } catch (err) {
          const text = this.isEditing ? this.$t('modals.view.fail.edit') : this.$t('modals.view.fail.create');

          this.$popups.error({ text });
        }
      }
    },
  },
};
</script>
