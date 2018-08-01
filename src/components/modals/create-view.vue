<template lang='pug'>
  v-card
    v-card-title
      span.headline {{ $t('modals.createWatcher.title') }}
    v-form
      v-layout(wrap, justify-center)
        v-flex(xs11)
          v-text-field(
          :label="$t('modals.createWatcher.displayName')",
          v-model="form.crecord_name",
          data-vv-name="name",
          v-validate="'required'",
          :error-messages="errors.collect('name')",
          )
          v-text-field(
          :label="$t('common.description')",
          v-model="form.description",
          data-vv-name="description",
          )
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
      form: {
        crecord_name: '',
        description: '',
      },
    };
  },
  methods: {
  },
};
</script>
