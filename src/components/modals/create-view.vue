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
          v-combobox(v-model='form.group_id',
          @click="fetchGroupList",
          :items="groupList",
          :label="$t('modals.createView.fields.groupIds')",
          )
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import viewMixin from '@/mixins/entities/viewV3/viewV3';
import groupMixin from '@/mixins/entities/viewV3/group';

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
      form: {
        name: '',
        title: '',
        group_id: '',
        description: '',
        enabled: '',
        tags: [],
      },
    };
  },
  methods: {
    remove(item) {
      this.form.tags.splice(this.form.tags.indexOf(item), 1);
      this.form.tags = [...this.form.tags];
    },
    async submit() {
      /* if (!this.groupList.include(this.form.group_id)) {

      } */
      const isFormValid = await this.$validator.validateAll();
      if (isFormValid) {
        this.createView(this.form);
      }
    },
  },
};
</script>
