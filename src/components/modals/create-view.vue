<template lang='pug'>
  v-card
    v-card-title
      span.headline {{ $t('modals.createWatcher.title') }}
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
          v-combobox(v-model='form.tags',
          label='Tags', tags='', clearable='',  multiple='', append-icon='')
            template(slot='selection', slot-scope='data')
              v-chip(:selected='data.selected', close='', @input='remove(data.item)') {{ data.item }}
          v-combobox(
          v-model='group_name',
          @click="fetchGroupList",
          :items="groupNames",
          :label="$t('modals.createView.fields.groupIds')",
          )
          span {{ this.form.group_id }}
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import viewMixin from '@/mixins/entities/viewV3/viewV3';
import groupMixin from '@/mixins/entities/viewV3/group';
import find from 'lodash/find';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalMixin, viewMixin, groupMixin],
  data() {
    return {
      group_name: '',
      form: {
        name: '',
        title: '',
        description: '',
        enabled: '',
        tags: [],
      },
    };
  },
  computed: {
    groupNames() {
      return this.groupList.map(group => group.name);
    },
  },
  methods: {
    remove(item) {
      this.form.tags.splice(this.form.tags.indexOf(item), 1);
      this.form.tags = [...this.form.tags];
    },
    async submit() {
      let groupId;
      if (this.groupNames.includes(this.group_name)) {
        groupId = find(this.groupList, { name: this.group_name })._id;
      } else {
        groupId = await this.createGroup({ name: this.group_name });
      }
      const data = {
        ...this.form,
        group_id: groupId,
      };
      const isFormValid = await this.$validator.validateAll();
      if (isFormValid) {
        await this.createView(data);
      }
    },
  },
};
</script>
