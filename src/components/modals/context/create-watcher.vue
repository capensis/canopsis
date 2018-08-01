<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.createWatcher.title') }}
    v-form
      v-layout(wrap, justify-center)
        v-flex(xs11)
          v-text-field(
            :label="$t('modals.createWatcher.displayName')",
            v-model="form.display_name",
            data-vv-name="name",
            v-validate="'required'",
            :error-messages="errors.collect('name')",
          )
      v-layout(wrap, justify-center)
        v-flex(xs11)
          h3.text-xs-center {{ $t('mFilterEditor.title') }}
          v-divider
          filter-editor
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import modalMixin from '@/mixins/modal/modal';

const { mapActions: watcherMapActions } = createNamespacedHelpers('watcher');
const { mapGetters: filterEditorMapGetters } = createNamespacedHelpers('mFilterEditor');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
  },
  mixins: [modalMixin],
  data() {
    return {
      form: {
        display_name: '',
        mfilter: '',
      },
    };
  },
  computed: {
    ...filterEditorMapGetters(['request']),
  },
  mounted() {
    if (this.config && this.config.item) {
      this.form = { ...this.config.item.props };
    }
  },
  methods: {
    ...watcherMapActions(['create']),
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const formData = {
          ...this.form,
          _id: this.form.display_name,
          type: 'watcher',
          mfilter: JSON.stringify(this.request),
        };
        this.create(formData);
        this.hideModal();
      }
    },
  },
};
</script>
