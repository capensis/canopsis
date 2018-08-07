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
          v-combobox(v-model='tags',
          label='Tags', tags='', clearable='',  multiple='', append-icon='')
            template(slot='selection', slot-scope='data')
              v-chip(:selected='data.selected', close='', @input='remove(data.item)') {{ data.item }}
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalMixin],
  data() {
    return {
      tags: [],
      form: {
        name: '',
        title: '',
        description: '',
        enabled: '',
      },
    };
  },
  methods: {
    remove(item) {
      this.tags.splice(this.tags.indexOf(item), 1);
      this.tags = [...this.tags];
    },
    async submit() {
      /* const isFormValid = await this.$validator.validateAll();
       if (isFormValid) {

      } */
    },
  },
};
</script>
