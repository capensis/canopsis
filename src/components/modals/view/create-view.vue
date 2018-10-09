<template lang="pug">
  v-card
    v-card-title
      span.headline {{ title }}
    v-form
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
          v-model="groupName",
          :items="groupNames",
          :label="$t('modals.view.fields.groupIds')",
          :search-input.sync="search"
          data-vv-name="group",
          v-validate="'required'",
          :error-messages="errors.collect('group')",
          )
            template(slot="no-data")
              v-list-tile
                v-list-tile-content
                  v-list-tile-title(v-html="$t('modals.view.noData')")

          span {{ form.group_id }}
      v-layout
        v-flex(xs6)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
        v-flex.text-xs-right(v-show="config.view", xs6)
          v-btn.red.darken-4.white--text(@click="remove") {{ $t('common.delete') }}
</template>

<script>
import find from 'lodash/find';

import { MODALS } from '@/constants';
import { generateView } from '@/helpers/entities';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import popupMixin from '@/mixins/popup';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    modalInnerMixin,
    entitiesViewMixin,
    entitiesViewGroupMixin,
    popupMixin,
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
      if (this.config.view) {
        return this.$t('modals.view.edit.title');
      }

      return this.$t('modals.view.create.title');
    },
  },
  mounted() {
    const { view } = this.config;

    if (view) {
      const group = find(this.groups, { _id: view.group_id });

      if (group) {
        this.groupName = group.name;
      }

      this.form = {
        name: view.name,
        title: view.title,
        description: view.description,
        enabled: view.enabled,
        tags: [...view.tags || []],
      };
    }
  },
  methods: {
    remove() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeView({ id: this.config.view._id });
            await this.fetchGroupsList();

            this.hideModal();
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

          const data = {
            ...generateView(),
            ...this.form,
            group_id: group._id,
          };

          if (this.config.view) {
            await this.updateView({ id: this.config.view._id, data });
          } else {
            await this.createView({ data });
          }

          await this.fetchGroupsList();

          this.addSuccessPopup({ text: this.$t('modals.view.success') });
          this.hideModal();
        }
      } catch (err) {
        this.addErrorPopup({ text: this.$t('modals.view.fail') });
        console.error(err.description);
      }
    },
  },
};
</script>
