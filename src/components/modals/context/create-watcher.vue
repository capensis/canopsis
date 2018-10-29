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
          h3.text-xs-center {{ $t('filterEditor.title') }}
          v-divider
          filter-editor(v-model="form.mfilter")
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

const { mapActions: watcherMapActions } = createNamespacedHelpers('watcher');

export default {
  name: MODALS.createWatcher,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
  },
  mixins: [modalInnerMixin],
  data() {
    const { item } = this.modal.config;

    let form = {
      name: '',
      mfilter: '{}',
    };

    if (item) {
      form = { ...item };
    }

    return {
      form,
    };
  },
  methods: {
    ...watcherMapActions(['create', 'edit']),

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          ...this.form,
          _id: this.config.item ? this.config.item._id : this.form.name,
          display_name: this.form.name,
          type: this.$constants.ENTITIES_TYPES.watcher,
          mfilter: this.form.mfilter,
        };

        try {
          if (this.config.item) {
            await this.edit({ data });
          } else {
            await this.create({ data });
          }

          this.hideModal();
        } catch (err) {
          console.error(err);
        }
      }
    },
  },
};
</script>
