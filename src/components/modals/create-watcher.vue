<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t(config.title) }}
    v-form
      v-layout(wrap, justify-center)
        v-flex(xs11)
          v-text-field(
            :label="$t('modals.createWatcher.displayName')",
            v-model="form.name",
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
import modalInnerMixin from '@/mixins/modal/modal-inner';

const { mapActions: watcherMapActions } = createNamespacedHelpers('watcher');
const { mapGetters: filterEditorMapGetters } = createNamespacedHelpers('mFilterEditor');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        name: '',
        mfilter: '',
      },
    };
  },
  computed: {
    ...filterEditorMapGetters(['request']),
  },
  mounted() {
    if (this.config.item) {
      this.form = { ...this.config.item };
    }
  },
  methods: {
    ...watcherMapActions(['create', 'edit']),
    async submit() {
      const formIsValid = await this.$validator.validateAll();

      if (formIsValid) {
        // If there's an item, means we're editing. If there's not, we're creating an entity
        if (this.config.item) {
          const formData = {
            ...this.form,
            _id: this.config.item._id,
            display_name: this.form.name,
            type: 'watcher',
            mfilter: JSON.stringify(this.request),
          };
          await this.edit({ watcher_id: formData._id, data: formData });
        } else {
          const formData = {
            ...this.form,
            _id: this.form.name,
            display_name: this.form.name,
            type: 'watcher',
            mfilter: JSON.stringify(this.request),
          };
          await this.create(formData);
        }
        this.hideModal();
      }
    },
  },
};
</script>
