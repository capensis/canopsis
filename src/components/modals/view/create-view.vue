<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.createView.title') }}
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
          :label="$t('modals.createView.fields.groupTags')",
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
          :label="$t('modals.createView.fields.groupIds')",
          :search-input.sync="search"
          data-vv-name="group",
          v-validate="'required'",
          :error-messages="errors.collect('group')",
          )
            template(slot="no-data")
              v-list-tile
                v-list-tile-content
                  v-list-tile-title(v-html="$t('modals.createView.noData')")

          span {{ form.group_id }}
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
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
  },
  mounted() {
    this.fetchGroupsList();
  },
  methods: {
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
          await this.createView({ data });
          await this.fetchGroupsList();

          this.addSuccessPopup({ text: this.$t('modals.createView.success') });

          this.hideModal();
        }
      } catch (err) {
        this.addErrorPopup({ text: this.$t('modals.createView.fail') });
        console.error(err.description);
      }
    },
  },
};
</script>
